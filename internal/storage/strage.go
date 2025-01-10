package storage

import "github.com/jmoiron/sqlx"

func MustInitDb(cs string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", cs)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}
