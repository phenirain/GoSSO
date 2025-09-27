package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidOldPassword error = errors.New("cтарый пароль не совпадает с текущим")
)

type User struct {
	Id           int64      `db:"id"`
	RoleId       int64      `db:"role_id"`
	Login        string     `db:"login"`
	PasswordHash []byte     `db:"password"`
	CreationTime time.Time  `db:"creation_datetime"`
	UpdateTime   *time.Time `db:"update_datetime"`
	IsArchived   bool       `db:"is_archived"`
}

func NewUser(login, password string, roleId *int64, isArchived *bool) *User {
	user := &User{
		Login: login,
	}
	user.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if roleId != nil {
		user.RoleId = *roleId
	} else {
		// 1 = покупатель - база
		user.RoleId = 1
	}
	user.CreationTime = time.Now()
	if isArchived != nil {
		user.IsArchived = *isArchived
	} else {
		user.IsArchived = false
	}

	return user
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	return err == nil
}

func (u *User) UpdateLogin(login string) {
	u.Login = login
	u.updateDateTime()
}

func (u *User) UpdatePassword(oldPass, newPass string) error {
	oldCorrect := u.CheckPassword(oldPass)
	if oldCorrect {
		u.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
		u.updateDateTime()
		return nil
	} else {
		return ErrInvalidOldPassword
	}
}

func (u *User) ChangeArchiveStatus(status bool) {
	u.IsArchived = status
	u.updateDateTime()
}

func (u *User) updateDateTime() {
	t := time.Now()
	u.UpdateTime = &t
}
