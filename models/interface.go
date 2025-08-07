package models

import (
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	BeforeCreate(tx *gorm.DB) (err error)
	Validate() []string
}

// user validations

func (User *User) BeforeCreate(tx *gorm.DB) (err error) {

	User.UserID = uuid.New().String()

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(User.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	User.Password = string(hashedPassword)
	hashedPassword = nil
	return nil
}

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
