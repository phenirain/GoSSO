package user

import (
	"context"
	"database/sql"
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

func New(db *sqlx.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (u *UserRepository) GetUser(ctx context.Context, login string) (*models.User, error) {
	const op = "User.GetUser"
	log := u.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user")

	var user models.User

	err := u.db.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("user not found")
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
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

func (u *UserRepository) AddUser(ctx context.Context, login string, passwordHash []byte) (int64, error) {
	const op = "User.AddUser"
	log := u.log.With(
		slog.String("op", op),
		slog.String("login", login),
		slog.String("passwordHash", string(passwordHash)),
	)
	log.Info("attempting to create new user")
	tx, err := u.db.Begin()
	if err != nil {
		log.Error("transaction cat not init", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("create transaction")
	var id int64
	row := tx.QueryRow("INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", login, passwordHash)
	err = row.Scan(&id)
	if err != nil {
		_ = tx.Rollback()
		log.Error("can't insert data in users table", sl.Err(err))
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully inserted data in users table")
	return id, tx.Commit()
}
