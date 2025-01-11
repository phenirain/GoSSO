package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/EtoNeAnanasbI95/auth-grpc-demo/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/ilyakaznacheev/cleanenv"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var configPath string
	var mp string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&mp, "migrations", "", "path to migrations dir")
	flag.Parse()

	var cfg config.Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Failed to parse config with: " + err.Error())
	}

	if cfg.ConnectionString == "" {
		panic("connection string is required")
	}

	m, err := migrate.New(
		"file://"+mp,
		cfg.ConnectionString,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied")
}

type Log struct {
	verbose bool
}

func (l *Log) Printf(format string, v ...any) {
	fmt.Printf(format, v...)
}

func (l *Log) Verbose() bool {
	return false
}
