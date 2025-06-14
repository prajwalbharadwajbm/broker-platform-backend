package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/prajwalbharadwajbm/broker/internal/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId string) (string, error) {
	// Access tokens should have a 10-minute validity.
	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "broker-platform",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := config.AppConfigInstance.JWTSecret
	if jwtSecret == "" {
		return "", errors.New("JWT secret not configured")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		jwtSecret := config.AppConfigInstance.JWTSecret
		if jwtSecret == "" {
			return nil, errors.New("JWT secret not configured")
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
