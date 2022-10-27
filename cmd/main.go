package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "csv-app",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Action: runServer,
			},
			{
				Name:   "migrate",
				Action: runMigrate,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
