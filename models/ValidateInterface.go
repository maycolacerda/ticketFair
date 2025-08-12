package models

import (
	"regexp"

	"github.com/maycolacerda/ticketfair/database"
)

var (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Pelo menos um dígito (0-9)
	hasNumber = regexp.MustCompile(`[0-9]`)
	// Pelo menos um caractere maiúsculo (A-Z)
	hasUpperCase = regexp.MustCompile(`[A-Z]`)
	// Pelo menos um caractere minúsculo (a-z)
	hasLowerCase = regexp.MustCompile(`[a-z]`)
	// Pelo menos um caractere especial (!@#$%&*)
	hasSymbol = regexp.MustCompile(`[!@#$%&*]`)
	// Pelo menos 8 caracteres no total
	hasMinLength = regexp.MustCompile(`.{8,}`)
	// Verifica se username possui algum caractere especial
	hasInvalidCharacter = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	// Verifica se a string possui apenas números
	hasonlynumbers = regexp.MustCompile(`^[0-9]*$`)
	// Verifica se a string possui apenas letras
	hasonlyletters = regexp.MustCompile(`^[a-zA-Z\s]*$`)
)

type ValidationsInterface interface {
	Validate() []string
}

// user validations

func (User *User) Validate() []string {
	var errors []string
	if User.Username == "" {
		errors = append(errors, "Username is required")
	}
	if !regexp.MustCompile(hasInvalidCharacter.String()).MatchString(User.Username) {
		errors = append(errors, "Username must not contain only letters and numbers")
	}
	if !regexp.MustCompile(emailRegex).MatchString(User.Email) {
		errors = append(errors, "Invalid email format")
	}
	if User.Email == "" {
		errors = append(errors, "Email is required")
	}
	if User.Password == "" {
		errors = append(errors, "Password is required")
	}
	if !validadePasswordComplexity(User.Password) {
		errors = append(errors, "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return errors
}

func validadePasswordComplexity(password string) bool {
	return hasNumber.MatchString(password) &&
		hasUpperCase.MatchString(password) &&
		hasLowerCase.MatchString(password) &&
		hasSymbol.MatchString(password) &&
		hasMinLength.MatchString(password)
}

// Login request validations
func (LoginRequest *LoginRequest) Validate() []string {
	var errors []string
	if LoginRequest.Email == "" {
		errors = append(errors, "Email is required")
	}
	if !regexp.MustCompile(emailRegex).MatchString(LoginRequest.Email) {
		errors = append(errors, "Invalid email format")
	}
	if LoginRequest.Password == "" {
		errors = append(errors, "Password is required")
	}
	if len(LoginRequest.Password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long")
	}

	return errors
}

func (Profile *Profile) Validate() []string {
	var errors []string
	if Profile.FirstName == "" {
		errors = append(errors, "First name is required")
	}
	if !regexp.MustCompile(hasonlyletters.String()).MatchString(Profile.FirstName) {
		errors = append(errors, "First name must contain only letters")
	}
	if Profile.LastName == "" {
		errors = append(errors, "Last name is required")
	}
	if !regexp.MustCompile(hasonlyletters.String()).MatchString(Profile.LastName) {
		errors = append(errors, "Last name must contain only letters")
	}
	if Profile.PhoneNumber == "" {
		errors = append(errors, "Phone number is required")
	}
	if !regexp.MustCompile(hasonlynumbers.String()).MatchString(Profile.PhoneNumber) {
		errors = append(errors, "Phone number must contain only numbers")
	}
	if Profile.Address == "" {
		errors = append(errors, "Address is required")
	}

	return errors
}

func (Merchant *Merchant) Validate() []string {
	var errors []string
	if Merchant.Name == "" {
		errors = append(errors, "Name is required")
	}
	if !regexp.MustCompile(hasInvalidCharacter.String()).MatchString(Merchant.Name) {
		errors = append(errors, "Name must contain only letters and numbers")
	}
	if Merchant.Description == "" {
		errors = append(errors, "Description is required")
	}
	if !regexp.MustCompile(hasInvalidCharacter.String()).MatchString(Merchant.Description) {
		errors = append(errors, "Description must contain only letters and numbers")
	}
	return errors
}

func (MerchantRep *MerchantRep) Validate() []string {
	var errors []string
	if MerchantRep.Name == "" {
		errors = append(errors, "Name is required")
	}
	if !regexp.MustCompile(hasInvalidCharacter.String()).MatchString(MerchantRep.Name) {
		errors = append(errors, "Name must contain only letters and numbers")
	}
	if MerchantRep.Role == "" {
		errors = append(errors, "Role is required")
	}
	if !regexp.MustCompile(hasInvalidCharacter.String()).MatchString(MerchantRep.Role) {
		errors = append(errors, "Role must contain only letters and numbers")
	}
	if MerchantRep.Email == "" {
		errors = append(errors, "Email is required")
	}
	if !regexp.MustCompile(emailRegex).MatchString(MerchantRep.Email) {
		errors = append(errors, "Invalid email format")
	}
	if MerchantRep.Password == "" {
		errors = append(errors, "Password is required")
	}
	if !validadePasswordComplexity(MerchantRep.Password) {
		errors = append(errors, "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}
	if MerchantRep.MerchantID == "" {
		errors = append(errors, "Merchant ID is required")
	}
	if err := database.DB.First(&MerchantRep, "merchant_id = ?", MerchantRep.MerchantID).Error; err != nil {
		errors = append(errors, "Merchant ID Not Found")
	}
	return errors
}

func (Event *Event) Validate() []string {
	var errors []string
	if Event.MerchantID == "" {
		errors = append(errors, "Merchant ID is required")
	}
	if Event.Name == "" {
		errors = append(errors, "Name is required")
	}
	if Event.Description == "" {
		errors = append(errors, "Description is required")
	}
	if Event.Date == "" {
		errors = append(errors, "Date is required")
	}
	if Event.Location == "" {
		errors = append(errors, "Location is required")
	}
	return errors
}
