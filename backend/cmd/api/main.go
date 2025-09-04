package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ipfs-force-community/janus/api"
	"github.com/ipfs-force-community/janus/database/mysql"
)

func main() {
	cmd := &cli.Command{
		Name:  "janus-backend",
		Usage: "Sync Filecoin chain data and serve API for visualization",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				Value:   "config.yaml",
			},
			&cli.Uint16Flag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "",
				Value:   10086,
			},
		},
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, c *cli.Command) error {
	configPath := c.String("config")

	config := mysql.Config{}
	if err := mysql.Load(configPath, &config); err != nil {
		return err
	}

	db, err := mysql.NewMySQLDB(config)
	if err != nil {
		return err
	}

	if err := api.NewServer(db).Run(c.Uint16("port")); err != nil {
		return err
	}

	return nil
}
