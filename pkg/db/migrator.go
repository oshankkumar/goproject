package db

import _ "github.com/golang-migrate/migrate/source/file"

import (
	"database/sql"
	"net/url"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/pkg/errors"
)

type Migrator interface {
	Up() error
	Rollback() error
}

type PostgresMigrator struct {
	db      *sql.DB
	migrate *migrate.Migrate
}

func NewPostgresMigrator(migrationPath string, db *sql.DB) (*PostgresMigrator, error) {
	f := &PostgresMigrator{db: db}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	sourceURL := (&url.URL{Scheme: "file", Path: migrationPath}).String()
	f.migrate, err = migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		return nil, errors.Wrap(err, "migrate.NewWithDatabaseInstance()")
	}

	return f, nil
}

func (f *PostgresMigrator) Up() error {
	return errors.Wrap(f.migrate.Up(), "Migrate.Up()")
}

func (f *PostgresMigrator) Rollback() error {
	return errors.Wrap(f.migrate.Steps(-1), "Migrate.Steps(-1)")
}
