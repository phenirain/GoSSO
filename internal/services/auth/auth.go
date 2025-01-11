package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/domain/models"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/lib/jwt"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/lib/logger/sl"
	userRepo "github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log      *slog.Logger
	tokenTTL time.Duration
	usrRepo  UsrRepo
	secret   []byte
}

func New(tokenTTL time.Duration, log *slog.Logger, repo UsrRepo, secret []byte) *Auth {
	return &Auth{
		tokenTTL: tokenTTL,
		log:      log,
		usrRepo:  repo,
		secret:   secret,
	}
}

var (
	ErrInvalidUserCredentials = errors.New("invalid user credentials")
	ErrUserAlreadyExists      = errors.New("user already exists")
)

//go:generate mockery --name=UsrRepo
type UsrRepo interface {
	GetUser(ctx context.Context, login string) (*models.User, error)
	GetUserWithId(ctx context.Context, uid int64) (*models.User, error)
	AddUser(ctx context.Context, login string, passwordHash []byte) (int64, error)
}

func (a *Auth) Login(ctx context.Context, login string, passwordHash string) (string, string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("attempting to login user")

	user, err := a.usrRepo.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, userRepo.ErrUserNotFound) {
			log.Info("user not found")
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidUserCredentials)
		}
		log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passwordHash)); err != nil {
		log.Info("invalid credentials", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, ErrInvalidUserCredentials)
	}

	log.Info("user logged in successfully")
	accessToken, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err := jwt.NewToken(user, time.Hour*24*30, a.secret)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

func (a *Auth) Register(ctx context.Context, login string, password string) (int64, error) {
	const op = "Auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
		slog.String("password", password),
	)

	log.Info("attempting to register user")

	user, err := a.usrRepo.GetUser(ctx, login)
	if user != nil {
		log.Info("found user")
		return -1, fmt.Errorf("%s: user with this login alredy exists, %w", op, ErrUserAlreadyExists)
	}
	if err != nil && !errors.Is(err, userRepo.ErrUserNotFound) {
		log.Error("failed to get user", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	newUid, err := a.usrRepo.AddUser(ctx, login, passHash)
	if err != nil {
		log.Error("failed to add user", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return newUid, nil
}

func (a *Auth) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	const op = "Auth.Refresh"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to refresh token")

	uid, err := a.Validate(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			log.Info("incorrect token", sl.Err(err))
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidUserCredentials)
		}
		log.Error("failed to refresh token", sl.Err(err))
		return "", "", err
	}

	user, err := a.usrRepo.GetUserWithId(ctx, uid)
	if err != nil {
		if errors.Is(err, userRepo.ErrUserNotFound) {
			log.Info("user not found")
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidUserCredentials)
		}
		log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	refreshToken, err = jwt.NewToken(user, time.Hour*24*30, a.secret)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, refreshToken, nil
}

func (a *Auth) Validate(ctx context.Context, token string) (int64, error) {
	const op = "Auth.Validate"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to validate jwt token")

	uid, err := jwt.ParseToken(token, a.secret)
	if err != nil {
		if errors.Is(err, jwt.ErrInvalidToken) {
			log.Info("invalid jwt token", sl.Err(err))
			return -1, fmt.Errorf("%s: %w", op, ErrInvalidUserCredentials)
		}
		log.Error("failed to validate jwt", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
