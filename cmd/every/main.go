package main

import (
	"errors"
	"io/fs"
	"log"
	"os"

	"github.com/r6m/every"
	"github.com/urfave/cli/v2"
)

var everyfileInit = `
  every "day" {
    user = "ubuntu"
		run = "hello world"
	}

	every "2 minutes at 12 am" {
		user = "ubuntu"
		run = "echo hello world"
	}`

type Config struct {
	Everies []Every `hcl:"every,block"`
}

// Every block every block data
type Every struct {
	Every string `hcl:"every,key"`
	User  string `hcl:"user"`
	Run   string `hcl:"run"`
}

func main() {
	app := &cli.App{
		Name:  "every",
		Usage: "Every command transles english to crontab expressions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Everyfile config path",
				Value:   "./Everyfile",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"u"},
				Usage:   "create a new Everyfile",
				Action: func(ctx *cli.Context) error {
					configPath := ctx.String("config")
					if exists(configPath) {

					}

					return nil
				},
			}, {
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update crontab",
				Action: func(ctx *cli.Context) error {
					configPath := ctx.String("config")

					config, err := every.Parse(configPath)
					if err != nil {
						return err
					}

					if err := every.WriteCrontab(config); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "clean",
				Usage: "clean crontab",
				Action: func(ctx *cli.Context) error {
					return every.CleanCrontab()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func exists(name string) bool {
	if _, err := os.Stat("/path/to/file"); errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}
