package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"cash-flow-financial/internal/db"
	"cash-flow-financial/internal/managers/loggermanager"
	"cash-flow-financial/internal/managers/rabbitmqmanager"
	"cash-flow-financial/internal/services/callback"

	"go.uber.org/zap"
)

type Worker struct {
	queries      *db.Queries
	rabbitMQ     *rabbitmqmanager.RabbitMQManager
	logger       *loggermanager.Logger
	callbackSvc  callback.ICallbackService
	queueName    string
	exchangeName string
	routingKey   string
}

func NewWorker(queries *db.Queries, rabbitMQ *rabbitmqmanager.RabbitMQManager, logger *loggermanager.Logger, callbackSvc callback.ICallbackService) IWorker {
	return &Worker{
		queries:      queries,
		rabbitMQ:     rabbitMQ,
		logger:       logger,
		callbackSvc:  callbackSvc,
		queueName:    "payment_intents_queue",
		exchangeName: "payment_intents_exchange",
		routingKey:   "payment.intent.created",
	}
}

func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("Starting payment worker...")

	err := w.rabbitMQ.Channel.ExchangeDeclare(
		w.exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.logger.Error("Failed to declare exchange", zap.Error(err))
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	queue, err := w.rabbitMQ.Channel.QueueDeclare(
		w.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.logger.Error("Failed to declare queue", zap.Error(err))
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = w.rabbitMQ.Channel.QueueBind(
		queue.Name,
		w.routingKey,
		w.exchangeName,
		false,
		nil,
	)
	if err != nil {
		w.logger.Error("Failed to bind queue", zap.Error(err))
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := w.rabbitMQ.Channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		w.logger.Error("Failed to register consumer", zap.Error(err))
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	w.logger.Info("Payment worker started successfully", zap.String("queue", queue.Name))

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.logger.Info("Worker shutting down...")
				return
			case d, ok := <-msgs:
				if !ok {
					w.logger.Error("Message channel closed")
					return
				}

				w.logger.Info("Received payment intent message",
					zap.String("message_id", d.MessageId),
					zap.String("body", string(d.Body)))

				var msg PaymentIntentMessage
				if err := json.Unmarshal(d.Body, &msg); err != nil {
					w.logger.Error("Failed to unmarshal message", zap.Error(err), zap.String("body", string(d.Body)))
					d.Nack(false, false)
					continue
				}

				if err := w.ProcessPaymentIntent(ctx, msg); err != nil {
					w.logger.Error("Failed to process payment intent", zap.Error(err), zap.Any("message", msg))
					d.Nack(false, true)
					continue
				}

				w.logger.Info("Successfully processed payment intent", zap.String("payment_intent_id", msg.PaymentIntentID))
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (w *Worker) Stop() error {
	w.logger.Info("Stopping payment worker...")
	return nil
}

func (w *Worker) ProcessPaymentIntent(ctx context.Context, message PaymentIntentMessage) error {
	w.logger.Info("Processing payment intent", zap.String("payment_intent_id", message.PaymentIntentID))

	w.logger.Info("Getting payment intent by string ID", zap.String("payment_intent_id", message.PaymentIntentID))
	paymentIntentInfo, err := w.queries.GetPaymentIntent(ctx, message.PaymentIntentID)
	if err != nil {
		w.logger.Error("Failed to get payment intent", zap.String("payment_intent_id", message.PaymentIntentID), zap.Error(err))
		return fmt.Errorf("failed to get payment intent: %w", err)
	}

	w.logger.Info("Found payment intent",
		zap.String("payment_intent_id", message.PaymentIntentID),
		zap.String("uuid", paymentIntentInfo.ID.String()),
		zap.String("status", string(paymentIntentInfo.Status.PaymentStatus)),
		zap.Bool("status_valid", paymentIntentInfo.Status.Valid))

	if paymentIntentInfo.Status.Valid && paymentIntentInfo.Status.PaymentStatus == db.PaymentStatusProcessing {
		w.logger.Info("Payment intent already being processed, stopping execution",
			zap.String("payment_intent_id", message.PaymentIntentID))
		return nil
	}

	if !paymentIntentInfo.Status.Valid || paymentIntentInfo.Status.PaymentStatus != db.PaymentStatusPending {
		w.logger.Info("Payment intent not in pending status, stopping execution",
			zap.String("payment_intent_id", message.PaymentIntentID),
			zap.String("current_status", string(paymentIntentInfo.Status.PaymentStatus)))
		return nil
	}

	w.logger.Info("Payment intent is pending, changing to processing and continuing",
		zap.String("payment_intent_id", message.PaymentIntentID))

	_, err = w.queries.UpdatePaymentIntentStatus(ctx, &db.UpdatePaymentIntentStatusParams{
		ID:       paymentIntentInfo.ID,
		Status:   db.NullPaymentStatus{PaymentStatus: db.PaymentStatusProcessing, Valid: true},
		Status_2: db.NullPaymentStatus{PaymentStatus: db.PaymentStatusPending, Valid: true},
	})
	if err != nil {
		w.logger.Error("Failed to update payment intent to processing status", zap.String("payment_intent_id", message.PaymentIntentID), zap.Error(err))
		return fmt.Errorf("failed to update payment intent status: %w", err)
	}

	w.logger.Info("Successfully changed payment intent status to processing", zap.String("payment_intent_id", message.PaymentIntentID))

	w.logger.Info("Creating payment transaction",
		zap.String("payment_intent_id", message.PaymentIntentID))

	thirdPartyRef := generateThirdPartyReference()
	selectedPaymentMethod := selectRandomPaymentMethod()
	accountNumber := generateAccountNumber(selectedPaymentMethod)

	w.logger.Info("Selected payment method and generated account number",
		zap.String("payment_method", string(selectedPaymentMethod)),
		zap.String("account_number", accountNumber))

	// Calculate fee (1% of amount)
	amountFloat, _ := strconv.ParseFloat(paymentIntentInfo.Amount, 64)
	feeAmount := fmt.Sprintf("%.2f", amountFloat*0.01) // 1% fee

	transaction, err := w.queries.CreatePaymentTransaction(ctx, &db.CreatePaymentTransactionParams{
		PaymentIntentID: paymentIntentInfo.PaymentIntentID,
		MerchantID:      paymentIntentInfo.MerchantID,
		Amount:          paymentIntentInfo.Amount,
		Currency:        db.CurrencyType(paymentIntentInfo.Currency),
		PaymentMethod:   db.NullPaymentMethodType{PaymentMethodType: selectedPaymentMethod, Valid: true},
		FeeAmount:       sql.NullString{String: feeAmount, Valid: true},
		AccountNumber:   sql.NullString{String: accountNumber, Valid: true},
	})
	if err != nil {
		w.logger.Error("Failed to create payment transaction", zap.Error(err))
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}

	w.logger.Info("=== PAYMENT TRANSACTION CREATED ===",
		zap.String("transaction_id", transaction.ID.String()),
		zap.String("payment_intent_id", transaction.PaymentIntentID),
		zap.String("merchant_id", transaction.MerchantID),
		zap.String("amount", transaction.Amount),
		zap.String("currency", string(transaction.Currency)),
		zap.String("payment_method", string(transaction.PaymentMethod.PaymentMethodType)),
		zap.String("fee_amount", transaction.FeeAmount.String),
		zap.String("account_number", transaction.AccountNumber.String),
		zap.String("status", string(transaction.Status.TransactionStatus)))

	merchant, err := w.queries.GetMerchantByMerchantID(ctx, paymentIntentInfo.MerchantID)
	if err != nil {
		w.logger.Error("Failed to get merchant UUID", zap.String("custom_merchant_id", paymentIntentInfo.MerchantID), zap.Error(err))
		return fmt.Errorf("failed to get merchant UUID: %w", err)
	}

	merchantBalanceBefore, err := w.queries.GetMerchantBalance(ctx, &db.GetMerchantBalanceParams{
		MerchantID: merchant.ID,
		Currency:   db.CurrencyType(paymentIntentInfo.Currency),
	})
	if err != nil {
		w.logger.Warn("Could not get merchant balance before update", zap.Error(err))
	} else {
		w.logger.Info("=== MERCHANT BALANCE BEFORE UPDATE ===",
			zap.String("merchant_id", paymentIntentInfo.MerchantID),
			zap.String("currency", string(paymentIntentInfo.Currency)),
			zap.String("available_balance", merchantBalanceBefore.AvailableBalance.String),
			zap.String("total_deposit", merchantBalanceBefore.TotalDeposit.String),
			zap.Int32("transaction_count", merchantBalanceBefore.TotalTransactionCount.Int32))
	}

	w.logger.Info("Updating payment transaction status",
		zap.String("payment_transaction_id", transaction.ID.String()),
		zap.String("third_party_ref", thirdPartyRef))

	_, err = w.queries.UpdatePaymentTransactionStatus(ctx, &db.UpdatePaymentTransactionStatusParams{
		ID:                  transaction.ID,
		Status:              db.NullTransactionStatus{TransactionStatus: db.TransactionStatusSuccess, Valid: true},
		ThirdPartyReference: sql.NullString{String: thirdPartyRef, Valid: true},
		Status_2:            db.NullTransactionStatus{TransactionStatus: db.TransactionStatusPending, Valid: true},
	})
	if err != nil {
		w.logger.Error("Failed to update payment transaction status", zap.Error(err))
		return fmt.Errorf("failed to update payment transaction status: %w", err)
	}

	_, err = w.queries.UpdatePaymentIntentStatus(ctx, &db.UpdatePaymentIntentStatusParams{
		ID:       paymentIntentInfo.ID,
		Status:   db.NullPaymentStatus{PaymentStatus: db.PaymentStatusSuccess, Valid: true},
		Status_2: db.NullPaymentStatus{PaymentStatus: db.PaymentStatusProcessing, Valid: true},
	})
	if err != nil {
		w.logger.Error("Failed to update payment intent status", zap.Error(err))
		return fmt.Errorf("failed to update payment intent status: %w", err)
	}

	amountFloat, _ = strconv.ParseFloat(paymentIntentInfo.Amount, 64)
	feeFloat, _ := strconv.ParseFloat(feeAmount, 64)
	depositAmount := amountFloat
	netBalanceAmount := depositAmount - feeFloat
	depositAmountStr := fmt.Sprintf("%.2f", depositAmount)
	netBalanceStr := fmt.Sprintf("%.2f", netBalanceAmount)

	w.logger.Info("Updating merchant balance",
		zap.String("merchant_uuid", merchant.ID.String()),
		zap.String("custom_merchant_id", paymentIntentInfo.MerchantID),
		zap.String("deposit_amount", depositAmountStr),
		zap.String("fee_amount", feeAmount),
		zap.String("net_balance_credit", netBalanceStr),
		zap.String("currency", string(paymentIntentInfo.Currency)))

	w.logger.Info("=== MERCHANT BALANCE UPDATE START ===",
		zap.String("merchant_uuid", merchant.ID.String()),
		zap.String("custom_merchant_id", paymentIntentInfo.MerchantID),
		zap.String("currency", string(paymentIntentInfo.Currency)),
		zap.String("deposit_amount", depositAmountStr),
		zap.String("fee_deducted", feeAmount),
		zap.String("balance_after_fee", netBalanceStr))

	merchantBalance, err := w.queries.IncrementMerchantBalance(ctx, &db.IncrementMerchantBalanceParams{
		MerchantID: merchant.ID,
		Currency:   db.CurrencyType(paymentIntentInfo.Currency),
		Column3:    depositAmountStr,
		Column4:    feeAmount,
	})

	if err != nil {
		w.logger.Error("=== MERCHANT BALANCE UPDATE FAILED ===",
			zap.String("merchant_id", paymentIntentInfo.MerchantID),
			zap.String("currency", string(paymentIntentInfo.Currency)),
			zap.String("deposit_amount", depositAmountStr),
			zap.Error(err))
	} else {
		w.logger.Info("=== MERCHANT BALANCE UPDATE SUCCESSFUL ===",
			zap.String("merchant_id", paymentIntentInfo.MerchantID),
			zap.String("currency", string(paymentIntentInfo.Currency)),
			zap.String("new_available_balance", merchantBalance.AvailableBalance.String),
			zap.String("new_total_deposit", merchantBalance.TotalDeposit.String),
			zap.Int32("new_transaction_count", merchantBalance.TotalTransactionCount.Int32),
			zap.String("deposit_amount", depositAmountStr))

		if err := w.sendCallback(paymentIntentInfo, transaction, thirdPartyRef, feeAmount, depositAmount); err != nil {
			w.logger.Warn("Failed to send callback to merchant",
				zap.String("payment_intent_id", message.PaymentIntentID),
				zap.String("callback_url", paymentIntentInfo.CallbackUrl),
				zap.Error(err))
		} else {
			w.logger.Info("Callback sent successfully to merchant",
				zap.String("payment_intent_id", message.PaymentIntentID),
				zap.String("callback_url", paymentIntentInfo.CallbackUrl))
		}
	}

	w.logger.Info("Payment processing completed successfully",
		zap.String("payment_intent_id", message.PaymentIntentID),
		zap.String("payment_transaction_id", transaction.ID.String()),
		zap.String("third_party_ref", thirdPartyRef),
		zap.String("merchant_deposit", depositAmountStr))

	return nil
}

func (w *Worker) sendCallback(paymentIntentInfo *db.GetPaymentIntentRow, transaction *db.PaymentTransaction, thirdPartyRef, feeAmount string, depositAmount float64) error {
	if paymentIntentInfo.CallbackUrl == "" {
		w.logger.Info("No callback URL provided, skipping callback",
			zap.String("payment_intent_id", paymentIntentInfo.PaymentIntentID))
		return nil
	}

	var metadata map[string]interface{}
	if paymentIntentInfo.Metadata.Valid {
		if err := json.Unmarshal(paymentIntentInfo.Metadata.RawMessage, &metadata); err != nil {
			w.logger.Warn("Failed to parse payment intent metadata for callback",
				zap.String("payment_intent_id", paymentIntentInfo.PaymentIntentID),
				zap.Error(err))
			metadata = make(map[string]interface{})
		}
	}

	processedAt := ""
	if transaction.ProcessedAt.Valid {
		processedAt = transaction.ProcessedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	callbackReq := callback.CallbackRequest{
		PaymentIntentID:     paymentIntentInfo.PaymentIntentID,
		MerchantID:          paymentIntentInfo.MerchantID,
		Amount:              depositAmount,
		Currency:            string(paymentIntentInfo.Currency),
		Status:              "success",
		AccountNumber:       transaction.AccountNumber.String,
		PaymentMethod:       string(transaction.PaymentMethod.PaymentMethodType),
		ThirdPartyReference: thirdPartyRef,
		FeeAmount:           feeAmount,
		ProcessedAt:         processedAt,
		Nonce:               paymentIntentInfo.Nonce,
		Metadata:            metadata,
	}

	return w.callbackSvc.SendCallback(paymentIntentInfo.CallbackUrl, callbackReq)
}
