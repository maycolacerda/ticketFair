// controllers/base.go
package controllers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/maycolacerda/ticketfair/dto"
)

var validate = validator.New()

func init() {
	dto.RegisterCustomValidators(validate) // ← register on startup
}

func formatValidationErrors(err error) map[string]string {
	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			errs[field] = field + " is required"
		case "email":
			errs[field] = "invalid email format"
		case "min":
			errs[field] = field + " must be at least " + e.Param() + " characters"
		case "max":
			errs[field] = field + " must be at most " + e.Param() + " characters"
		case "alphanum":
			errs[field] = field + " must contain only letters and numbers"
		case "oneof":
			errs[field] = field + " must be one of: " + e.Param()
		case "password":
			errs[field] = "password must contain at least 8 characters, one uppercase, one lowercase, one number, and one special character (!@#$%&*)"
		case "onlyletters":
			errs[field] = field + " must contain only letters"
		case "onlynumbers":
			errs[field] = field + " must contain only numbers"
		case "gtfield":
			errs[field] = field + " must be after " + e.Param()
		default:
			errs[field] = field + " is invalid"
		}
	}
	return errs
}
