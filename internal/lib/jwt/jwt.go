package jwt

import (
	"errors"
	"fmt"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func NewToken(user *models.User, duration time.Duration, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.Id,
		"name": user.Name,
	}
	claims["exp"] = time.Now().Add(duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string, secret []byte) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return -1, fmt.Errorf("token parse error: %s", err.Error())
	}
	if !token.Valid {
		return -1, ErrInvalidToken
	}
	uid, ok := token.Claims.(jwt.MapClaims)["sub"]
	if !ok {
		return -1, errors.New("can't get sub from claims")
	}
	return int64(uid.(float64)), nil
}
