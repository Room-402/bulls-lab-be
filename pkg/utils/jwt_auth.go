package utils

import (
	"bulls-lab-be/internal/core/domain"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const JWT_SECRET = "bulls-lab-be-sample"

type JwtCustomClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func getSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = JWT_SECRET
	}
	return []byte(secret)
}

func GenerateToken(user domain.User) (string, error) {
	claims := JwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
}

func ValidateToken(tokenString string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
