package main

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/urfave/cli/v2"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/pkg/migrator"
)

func runMigrate(*cli.Context) error {
	cfg, err := conf.New()
	if err != nil {
		return fmt.Errorf("conf new: %w", err)
	}

	if err = migrator.MakeMigrate(cfg); err != nil {
		return fmt.Errorf("migrator make migrate: %w", err)
	}

	return nil
}
