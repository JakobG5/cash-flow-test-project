package checkout

import (
	"cash-flow-financial/internal/models"

	"github.com/go-playground/validator/v10"
)

func (h *CheckoutHandler) validateCreatePaymentIntentRequest(req models.CreatePaymentIntentRequest) []string {
	validate := validator.New()
	var errorMessages []string

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "Amount":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "amount is required")
				case "gt":
					errorMessages = append(errorMessages, "amount must be greater than 0")
				case "lte":
					errorMessages = append(errorMessages, "amount cannot exceed 100,000")
				}
			case "Currency":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "currency is required")
				case "len":
					errorMessages = append(errorMessages, "currency must be exactly 3 characters")
				case "oneof":
					errorMessages = append(errorMessages, "currency must be one of: ETB, USD, EUR, GBP")
				}
			case "CallbackURL":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "callback_url is required")
				case "url":
					errorMessages = append(errorMessages, "callback_url must be a valid URL")
				}
			case "Nonce":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "nonce is required")
				case "min":
					errorMessages = append(errorMessages, "nonce must be at least 16 characters")
				case "max":
					errorMessages = append(errorMessages, "nonce must be at most 64 characters")
				}
			case "Description":
				if fieldError.Tag() == "max" {
					errorMessages = append(errorMessages, "description must be at most 500 characters")
				}
			}
		}
	}

	return errorMessages
}

func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 10 {
		return apiKey
	}
	return apiKey[:10] + "..."
}
