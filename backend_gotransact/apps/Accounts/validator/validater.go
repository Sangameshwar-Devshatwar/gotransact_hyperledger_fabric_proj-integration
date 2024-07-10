package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// ValidatePassword checks if the password meets the complexity requirements
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen    = len(password) >= 8
		hasUpperCase = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowerCase = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial   = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	)

	return hasMinLen && hasUpperCase && hasLowerCase && hasNumber && hasSpecial
}

var validate *validator.Validate

// Init initializes the custom validator
func Init() {
	validate = validator.New()
	validate.RegisterValidation("password_complexity", ValidatePassword)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
