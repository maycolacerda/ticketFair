// dto/validators.go
package dto

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	hasNumber        = regexp.MustCompile(`[0-9]`)
	hasUpperCase     = regexp.MustCompile(`[A-Z]`)
	hasLowerCase     = regexp.MustCompile(`[a-z]`)
	hasSymbol        = regexp.MustCompile(`[!@#$%&*]`)
	hasMinLength     = regexp.MustCompile(`.{8,}`)
	onlyLettersRegex = regexp.MustCompile(`^[a-zA-Z\s]+$`)
	onlyNumbersRegex = regexp.MustCompile(`^[0-9]+$`)
)

// RegisterCustomValidators registers all custom rules on the provided validator instance
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("onlyletters", validateOnlyLetters)
	v.RegisterValidation("onlynumbers", validateOnlyNumbers)
}

// validatePassword enforces complexity: upper, lower, number, symbol, min 8 chars
func validatePassword(fl validator.FieldLevel) bool {
	p := fl.Field().String()
	return hasNumber.MatchString(p) &&
		hasUpperCase.MatchString(p) &&
		hasLowerCase.MatchString(p) &&
		hasSymbol.MatchString(p) &&
		hasMinLength.MatchString(p)
}

func validateOnlyLetters(fl validator.FieldLevel) bool {
	return onlyLettersRegex.MatchString(fl.Field().String())
}

func validateOnlyNumbers(fl validator.FieldLevel) bool {
	return onlyNumbersRegex.MatchString(fl.Field().String())
}
