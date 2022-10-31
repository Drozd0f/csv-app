package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// @title       CSV-APP API
// @version     1.0
// @description Виконання тестового завдання в EVO 2022

// @query.collection.format  multi
// @host     localhost:4444
// @BasePath /api/v1
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
