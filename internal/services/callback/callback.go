package callback

type CallbackRequest struct {
	PaymentIntentID    string                 `json:"payment_intent_id"`
	MerchantID         string                 `json:"merchant_id"`
	Amount             float64                `json:"amount"`
	Currency           string                 `json:"currency"`
	Status             string                 `json:"status"`
	AccountNumber      string                 `json:"account_number,omitempty"`
	PaymentMethod      string                 `json:"payment_method,omitempty"`
	ThirdPartyReference string                `json:"third_party_reference,omitempty"`
	FeeAmount          string                 `json:"fee_amount,omitempty"`
	ProcessedAt        string                 `json:"processed_at,omitempty"`
	Nonce              string                 `json:"nonce"`
	Metadata           map[string]interface{} `json:"metadata,omitempty"`
}

type CallbackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
}

type ICallbackService interface {
	SendCallback(callbackURL string, request CallbackRequest) error
}
