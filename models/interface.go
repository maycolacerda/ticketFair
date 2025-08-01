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
)

type ValidationsInterface interface {
	BeforeCreate(tx *gorm.DB) (err error)
	Validate() []string
}

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
	if User.Username == "" {
		errors = append(errors, "Username is required")
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
