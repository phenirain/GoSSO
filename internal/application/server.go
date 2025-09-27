package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phenirain/sso/internal/application/auth"
	"github.com/phenirain/sso/internal/config"
	"github.com/phenirain/sso/pkg/echomiddleware"
)



func SetupHTTPServer(cfg *config.Config, authService auth.AuthService, jwt echomiddleware.Jwt) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echomiddleware.JwtValidation(jwt))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))


	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	registerAuthRoutes(e, authService)

	return e
}

func registerAuthRoutes(e *echo.Echo, authService auth.AuthService) {
	authHandler := auth.NewHandler(authService)
	auth := e.Group("/auth")
	auth.POST("/logIn", authHandler.LogIn)
	auth.POST("/signUp", authHandler.SignUp)
	auth.POST("/refresh", authHandler.Refresh)
}

