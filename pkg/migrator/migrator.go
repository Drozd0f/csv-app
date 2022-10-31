package migrator

import (
	"errors"
	"fmt"
	"log"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func MakeMigrate(cfg *conf.Config) error {
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
