package echomiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/phenirain/sso/pkg/contextkeys"
)

type Jwt interface {
	ParseToken(tokenString string) (int64, error)
}

func JwtValidation(jwt Jwt) echo.MiddlewareFunc {
	skip := map[string]struct{}{
		"/auth/logIn":   {},
		"/auth/signUp":  {},
		"/auth/refresh": {},
		"/health":       {},
		"/swagger/*":    {},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			p := c.Path()
			if _, ok := skip[p]; ok {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.ErrUnauthorized
			}
			tokenString := parts[1]

			userId, err := jwt.ParseToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": err.Error(),
				})
			}

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, contextkeys.UserIDCtxKey, userId)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
