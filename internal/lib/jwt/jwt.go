package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	jwtErrors "github.com/phenirain/sso/internal/errors/jwt"
	"time"
)

type JwtLib struct {
	duration time.Duration
	secret []byte
}

func NewJwtLib(duration time.Duration, secret []byte) *JwtLib {
	return &JwtLib{
		duration: duration,
		secret: secret,
	}
}

//TODO: класть роль
func (j *JwtLib) NewToken(userId int64) (accessToken string, refreshToken string, error error) {
	claims := jwt.MapClaims{
		"sub":  userId,
	}
	claims["exp"] = time.Now().Add(j.duration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(j.secret)
	if err != nil {
		return "", "", err
	}

	claims["exp"] = time.Now().Add(time.Hour*24*30).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err = token.SignedString(j.secret)
	if err != nil {
		return "", "", err
	}
	return
}

//TODO доставать роль
func (j *JwtLib) ParseToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return -1, fmt.Errorf("token parse error: %s", err.Error())
	}
	if !token.Valid {
		return -1, jwtErrors.ErrInvalidToken
	}
	uid, ok := token.Claims.(jwt.MapClaims)["sub"]
	if !ok {
		return -1, errors.New("can't get sub from claims")
	}
	return int64(uid.(float64)), nil
}
