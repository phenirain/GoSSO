package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/domain/models"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/lib/logger/sl"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct {
	log *slog.Logger
	db  *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUser(ctx context.Context, email string) (*models.User, error) {
	const op = "User.GetUser"
	log := u.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user")

	var user models.User

	err := u.db.QueryRowxContext(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(&user)
	if err != nil {
		log.Error("something went wrong", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (u *UserRepository) GetUserWithId(ctx context.Context, uid int64) (*models.User, error) {
	const op = "User.GetUserWithId"
	log := u.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user with id")

	var user models.User

	err := u.db.Get(&user, "SELECT * FROM users WHERE id = $1", uid)
	if err != nil {
		log.Error("something went wrong", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (u *UserRepository) AddUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	const op = "User.AddUser"
	log := u.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to create new user")
	tx, err := u.db.Begin()
	if err != nil {
		log.Error("transaction cat not init", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("create transaction")
	var uid int64
	row := tx.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", email, passwordHash)
	err = row.Scan(&uid)
	if err != nil {
		log.Error("can't insert data in users table", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully inserted data in users table")
	return uid, nil
}
