package auth

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	authModels "github.com/phenirain/sso/internal/dto/auth"
	"github.com/phenirain/sso/internal/dto/response"
)

type AuthService interface {
	Auth(ctx context.Context, request authModels.AuthRequest, isNew bool) (*authModels.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*authModels.AuthResponse, error)
}

type Handler struct {
	s AuthService
}

func NewHandler(auth AuthService) *Handler {
	return &Handler{
		s: auth,
	}
}

// Login handles user login
func (h *Handler) LogIn(c echo.Context) error {
	return h.auth(c, false)
}

func (h *Handler) SignUp(c echo.Context) error {
	return h.auth(c, true)
}

// Refresh handles token refresh
func (h *Handler) Refresh(c echo.Context) error {
	ctx := c.Request().Context()

	// Получаем токен из заголовка Authorization
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Отсутствует токен", "Заголовок Authorization обязателен"))
	}

	// Проверяем формат "Bearer <token>"
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Неверный формат токена", "Используйте формат: Bearer <token>"))
	}

	refreshToken := authHeader[7:] // Убираем "Bearer "

	result, err := h.s.Refresh(ctx, refreshToken)
	if err != nil {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Ошибка обновления токена", err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}

func (h *Handler) auth(c echo.Context, isNew bool) error {
	ctx := c.Request().Context()

	var req authModels.AuthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Ошибка чтения json", err.Error()))
	}

	if req.Login == "" {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Отсутствует аргумент", "Логин обязателен"))
	}
	if req.Password == "" {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Отсутствует аргумент", "Пароль обязателен"))
	}

	result, err := h.s.Auth(ctx, req, isNew)
	if err != nil {
		return c.JSON(http.StatusOK, response.NewBadResponse[any]("Ошибка авторизации", err.Error()))
	}

	return c.JSON(http.StatusOK, response.NewSuccessResponse(result))
}
