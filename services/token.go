package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

// The `createToken` function generates a JWT token with a username claim and expiration time, signed
// with a secret key.
func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// The `verifyToken` function in Go parses and verifies a JWT token using a secret key.
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
