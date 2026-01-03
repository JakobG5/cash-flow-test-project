package transactionservice

type ITransactionService interface {
	GetPaymentStatus(transactionID string) error
	ProcessPayment(transactionID string) error
}
