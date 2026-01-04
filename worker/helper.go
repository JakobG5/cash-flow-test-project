package worker

import (
	"cash-flow-financial/internal/db"
	"crypto/rand"
	"fmt"
	"math/big"
)

func generateThirdPartyReference() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}

func selectRandomPaymentMethod() db.PaymentMethodType {
	methods := []db.PaymentMethodType{
		db.PaymentMethodTypeCbe,
		db.PaymentMethodTypeMpesa,
		db.PaymentMethodTypeTelebirr,
		db.PaymentMethodTypeAwash,
	}
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(methods))))
	return methods[num.Int64()]
}

func generateAccountNumber(paymentMethod db.PaymentMethodType) string {
	digits := make([]byte, 8)
	for i := range digits {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		digits[i] = byte('0' + num.Int64())
	}

	var prefix string
	switch paymentMethod {
	case db.PaymentMethodTypeMpesa:
		prefix = "2517"
	default:
		prefix = "2519"
	}

	return fmt.Sprintf("%s%s", prefix, string(digits))
}
