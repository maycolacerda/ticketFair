package models

import (
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
	Username string `json:"username" gorm:"uniqueIndex"`
	UserID   string `json:"user_id" gorm:"uniqueIndex"` // Unique identifier for the user
}

var (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	// Pelo menos um dígito (0-9)
	hasNumber = regexp.MustCompile(`[0-9]`)
	// Pelo menos um caractere maiúsculo (A-Z)
	hasUpperCase = regexp.MustCompile(`[A-Z]`)

	hasSymbol = regexp.MustCompile(`[!@#$%^&*]`)
	// Pelo menos 8 caracteres no total
	hasMinLength = regexp.MustCompile(`.{8,}`)
)

func (User *User) Validate() []string {
	var errors []string
	if !regexp.MustCompile(emailRegex).MatchString(User.Email) {
		errors = append(errors, "Invalid email format;")
	}
	if User.Email == "" {
		errors = append(errors, "Email is required;")
	}
	if User.Password == "" {
		errors = append(errors, "Password is required;")
	}
	if !validadePassword(User.Password) {
		errors = append(errors, "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character;")
	}
	if User.Username == "" {
		errors = append(errors, "Username is required;")
	}
	if User.UserID == "" {
		errors = append(errors, "UserID is required;")
	}

	return errors
}

func validadePassword(password string) bool {
	return hasNumber.MatchString(password) &&
		hasUpperCase.MatchString(password) &&
		hasSymbol.MatchString(password) &&
		hasMinLength.MatchString(password)
}
