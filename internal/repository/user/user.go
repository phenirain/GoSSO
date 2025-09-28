package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/phenirain/sso/internal/domain"
	"github.com/phenirain/sso/pkg/database"
)

type UserRepository struct {
	db  *sqlx.DB
}

func New(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (*domain.User, error) {
	const op = "User.GetUserByLogin"
	log := slog.With(
		slog.String("op", op),
	)
	log.Info("attempting to get user")

	var user domain.User
	err := u.db.Get(&user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		log.Error("something went wrong", "err", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
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

func (u *UserRepository) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	const query = `
		INSERT INTO users (role_id, login, password, creation_datetime, update_datetime, is_archived)
		VALUES (:role_id, :login, :password, :creation_datetime, :update_datetime, :is_archived)
		RETURNING id
	`

	result, err := database.WithUserTransaction(u.db, ctx, func(tx *sqlx.Tx) (int64, error) {
		rows, err := tx.NamedQuery(query, user)
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var id int64
		if rows.Next() {
			if err := rows.Scan(&id); err != nil {
				return 0, err
			}
		} else {
			return 0, fmt.Errorf("no id returned")
		}

		return id, nil
	})
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return result, nil
}

