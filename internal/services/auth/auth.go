package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/phenirain/sso/internal/domain"
	"github.com/phenirain/sso/internal/dto/auth"
	authErrors "github.com/phenirain/sso/internal/errors/auth"
	"github.com/phenirain/sso/internal/errors/jwt"
)

type Jwt interface {
	NewToken(userId int64) (accessToken string, refreshToken string, error error)
	ParseToken(tokenString string) (int64, error)
}

type Repository interface {
	GetUserByLogin(ctx context.Context, login string) (*domain.User, error)
	GetUserWithId(ctx context.Context, uid int64) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (int64, error)
}


type Auth struct {
	repo  Repository
	jwt Jwt
}

func New(repo Repository, jwt Jwt) *Auth {
	return &Auth{
		repo:  repo,
		jwt: jwt,
	}
}

func (a *Auth) Auth(ctx context.Context, request auth.AuthRequest, isNew bool) (*auth.AuthResponse, error) {
	const op string = "Auth.Login"

	user, err := a.repo.GetUserByLogin(ctx, request.Login)
	if err != nil {
		slog.Error("failed to get user", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var userId int64 
	// если создание
	if isNew {
		// если пользователь найден - уже существует
		if user != nil {
			return nil, authErrors.ErrUserAlreadyExists
		}

		user = domain.NewUser(request.Login, request.Password, nil, nil)
		userId, err = a.repo.CreateUser(ctx, user)
		if err != nil {
			errText := fmt.Errorf("ошибка в ходе создания пользователя: %w", err)
			slog.Error(errText.Error())
			return nil, errText
		}
	} else { // если авторизация
		// если пользователь не найден
		if user == nil {
			return nil, authErrors.ErrInvalidUserCredentials
		}
		valid := user.CheckPassword(request.Password)
		// если пароль не верен
		if !valid {
			return nil, authErrors.ErrInvalidUserCredentials
		}
		userId = user.Id
	}

	return a.getAuthResponse(userId)
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (*auth.AuthResponse, error) {

	// проверка токена
	userId, err := a.jwt.ParseToken(refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			return nil, err
		}
		slog.Error("ошибка парсинга токена", "err", err)
		return nil, err
	}

	// проверка пользователя
	user, err := a.repo.GetUserWithId(ctx, userId)
	if err != nil {
		errorText := fmt.Errorf("ошибка получения пользователя по идентфикатору: %w", err)
		slog.Error(errorText.Error())
		return nil, errorText
	}
	// если его нет или удален - нахуй
	if user == nil || user.IsArchived {
		return nil, authErrors.ErrUserNotFound
	}

	return a.getAuthResponse(userId)
}

func (a *Auth) getAuthResponse(userId int64) (*auth.AuthResponse, error) {
	accessToken, refreshToken, err := a.jwt.NewToken(userId)
	if err != nil {
		errorText := fmt.Errorf("ошибка генерации токенов доступа: %w", err)
		slog.Error(errorText.Error())
		return nil, errorText
	}

	return &auth.AuthResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}
