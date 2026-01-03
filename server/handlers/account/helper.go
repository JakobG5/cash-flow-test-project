package account

import (
	"cash-flow-financial/internal/models"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (h *AccountHandler) validateCreateMerchantRequest(req models.CreateMerchantRequest) []string {
	validate := validator.New()
	var errorMessages []string

	if err := validate.Struct(req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "Name":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "name is required")
				case "min":
					errorMessages = append(errorMessages, "name must be at least 2 characters")
				case "max":
					errorMessages = append(errorMessages, "name must be at most 100 characters")
				}
			case "Email":
				switch fieldError.Tag() {
				case "required":
					errorMessages = append(errorMessages, "email is required")
				case "email":
					errorMessages = append(errorMessages, "invalid email format")
				case "max":
					errorMessages = append(errorMessages, "email must be at most 255 characters")
				}
			}
		}
	}

	// Custom validation: name can only contain letters and spaces
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if !nameRegex.MatchString(strings.TrimSpace(req.Name)) {
		errorMessages = append(errorMessages, "name can only contain letters and spaces")
	}

	return errorMessages
}
