package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/urfave/cli/v2"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/db/migrations"
)

func runMigrate(*cli.Context) error {
	cfg, err := conf.New()
	if err != nil {
		return fmt.Errorf("conf new: %w", err)
	}

	src, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		return fmt.Errorf("iofs new: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, cfg.DBURI)
	if err != nil {
		return fmt.Errorf("migrate new with source instance: %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up: %w", err)
	}

	v, d, err := m.Version()
	if err != nil {
		return fmt.Errorf("migration get version: %w", err)
	}

	log.Printf("current migration %d, is %t\n", v, d)

	return nil
}
