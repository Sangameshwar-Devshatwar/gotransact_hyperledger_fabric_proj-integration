package structutils

type PaymentRequest struct {
	CardNumber  string `json:"cardNumber" binding:"required" validate:"cardNumber"`
	ExpiryDate  string `json:"expiryDate" binding:"required" validate:"expiryDate"`
	CVV         string `json:"cvv" binding:"required" validate:"cvv"`
	Amount      string `json:"amount" binding:"required,gt=0"`
	Description string `json:"description"`
}
