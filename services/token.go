// services/token.go
package services

import (
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleClient          = "client"
	RoleMerchant        = "merchant"
	RoleMerchantStaff   = "staff"
	RoleMerchantAdmin   = "admin"
	RoleMerchantManager = "manager"
)

type Claims struct {
	UserID     string `json:"user_id"`
	Role       string `json:"role"`
	MerchantID string `json:"merchant_id,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, role, merchantID string) (string, int64, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", 0, errors.New("JWT_SECRET is not set")
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	claims := Claims{
		UserID:     userID,
		Role:       role,
		MerchantID: merchantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ticketfair",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, errors.New("failed to sign token")
	}

	return signed, expiresAt.Unix(), nil
}

// ParseToken is pure logic — no gin, fully testable
func ParseToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty") // ← was missing
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// extractBearerToken pulls the raw token string from the Authorization header
func ExtractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(prefix):], nil
}

// extractClaims is the single internal entry point for parsing claims from a request
func extractClaims(c *gin.Context) (*Claims, error) {
	tokenStr, err := ExtractBearerToken(c)
	if err != nil {
		return nil, err
	}
	return ParseToken(tokenStr)
}

// Public helpers — all share a single extractClaims call

func ExtractTokenID(c *gin.Context) (string, error) {
	claims, err := extractClaims(c)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func ExtractTokenRole(c *gin.Context) (string, error) {
	claims, err := extractClaims(c)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

func ExtractMerchantID(c *gin.Context) (string, error) {
	claims, err := extractClaims(c) // ← fixed: was ExtractClaims (undefined)
	if err != nil {
		return "", err
	}
	if claims.MerchantID == "" {
		return "", errors.New("merchant_id not present in token")
	}
	return claims.MerchantID, nil
}
