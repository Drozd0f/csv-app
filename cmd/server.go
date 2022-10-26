package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/repository"
	"github.com/Drozd0f/csv-app/server"
	"github.com/Drozd0f/csv-app/service"
)

func runServer(c *cli.Context) error {
	cfg, err := conf.New()
	if err != nil {
		return fmt.Errorf("conf new: %w", err)
	}

	r, err := repository.New(
		c.Context,
		cfg.DBURI,
	)
	if err != nil {
		return fmt.Errorf("repository new: %w", err)
	}

	s := service.New(r, cfg)
	serv := server.New(s, cfg)

	if err = serv.Run(cfg.Addr); err != nil {
		return fmt.Errorf("server run: %w", err)
	}

	return nil
}
