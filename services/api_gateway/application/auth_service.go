package application

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secretKey   string
	issuer      string
	expireHours int
}

type CustomClaims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func NewJWTService(secretKey, issuer string, expireHours int) *JWTService {
	return &JWTService{
		secretKey:   secretKey,
		issuer:      issuer,
		expireHours: expireHours,
	}
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expireHours))),
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Если токен прошел проверку, извлекаем из него наши claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
