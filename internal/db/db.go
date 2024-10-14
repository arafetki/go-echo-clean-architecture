package database

import (
	"context"
	"errors"
	"time"

	"github.com/arafetki/go-echo-clean-architecture/assets"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	*sqlx.DB
}

type Options struct {
	ConnTimeout     time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	AutoMigrate     bool
}

func Postgres(dsn string, options Options) (*DB, error) {

	ctx, cancel := context.WithTimeout(context.Background(), options.ConnTimeout)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", "postgres://"+dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(options.MaxOpenConns)
	db.SetMaxIdleConns(options.MaxIdleConns)
	db.SetConnMaxIdleTime(options.ConnMaxIdleTime)
	db.SetConnMaxLifetime(options.ConnMaxLifetime)

	if options.AutoMigrate {

		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, err
		}
		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "postgres://"+dsn)
		if err != nil {
			return nil, err
		}
		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		default:
			return nil, err
		}

	}

	return &DB{db}, nil
}
