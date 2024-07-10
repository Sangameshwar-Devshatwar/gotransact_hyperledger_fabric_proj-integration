package validators

import (
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func cardNumberValidator(fl validator.FieldLevel) bool {
	cardNumber := fl.Field().String()
	// Check if card number is 16 or 18 digits
	match, _ := regexp.MatchString(`^\d{16}|\d{18}$`, cardNumber)
	return match
}

func expiryDateValidator(fl validator.FieldLevel) bool {
	expiryDate := fl.Field().String()
	// Check if expiry date is in the format MM/YY and within 10 years span
	t, err := time.Parse("01/06", expiryDate)
	if err != nil {
		return false
	}
	currentYear := time.Now().Year() % 100
	expiryYear := t.Year() % 100
	if expiryYear < currentYear || expiryYear > currentYear+10 {
		return false
	}
	return true
}

func cvvValidator(fl validator.FieldLevel) bool {
	cvv := fl.Field().String()
	// Check if CVV is exactly 3 digits
	match, _ := regexp.MatchString(`^\d{3}$`, cvv)
	return match
}

// CustomErrorMessages contains custom error messages for validation
var CustomErrorMessages = map[string]string{
	"cardNumber": "Card number must be 16 or 18 digits.",
	"expiryDate": "Expiry date must be in MM/YY format and within a 10 year span.",
	"cvv":        "CVV must be exactly 3 digits.",
	"amount":     "Amount must be greater than 0.",
}

// InitValidation initializes the custom validators

func InitValidation() {
	validate = validator.New()
	validate.RegisterValidation("cardNumber", cardNumberValidator)
	validate.RegisterValidation("expiryDate", expiryDateValidator)
	validate.RegisterValidation("cvv", cvvValidator)
}
func GetValidator() *validator.Validate {
	return validate
}
