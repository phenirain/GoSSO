package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/phenirain/sso/internal/domain"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type UserRepository struct {
	db  *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (user *domain.User, err error) {
	const op = "User.GetUser"
	log := slog.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user")

	err = u.db.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error("something went wrong", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return
}

func (u *UserRepository) GetUserWithId(ctx context.Context, uid int64) (*domain.User, error) {
	const op = "User.GetUserWithId"
	log := slog.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user with id")

	var user domain.User

	err := u.db.Get(&user, "SELECT * FROM users WHERE id = $1", uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error("something went wrong", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error){
	const query = `
		INSERT INTO users (role_id, login, password, creation_datetime, update_datetime, is_archived)
		VALUES (:role_id, :login, :password, :creation_datetime, :update_datetime, :is_archived)
		RETURNING id
	`

	rows, err := u.db.NamedQueryContext(ctx, query, u)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	defer rows.Close()

	var id int64
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("scan id: %w", err)
		}
	}

	return id, nil
}
