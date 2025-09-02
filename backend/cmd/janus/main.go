package main

import (
	"context"
	"log"
	"os"

	"github.com/ipfs-force-community/janus/chain"
	"github.com/ipfs-force-community/janus/database/mysql"
	"github.com/urfave/cli/v3"
)

type contextKey string

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
			&cli.StringFlag{
				Name:    "node_endpoint",
				Aliases: []string{"n"},
				Usage:   "Filecoin node endpoint",
				Value:   "http://127.0.0.1:3463",
			},
			&cli.Int64Flag{
				Name:  "start_epoch",
				Usage: "Start epoch to sync from",
				Value: 0,
			},
			&cli.Int64Flag{
				Name:  "end_epoch",
				Usage: "End epoch to sync to, 0 means sync to the latest",
				Value: 0,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			configPath := c.String("config")

			config := mysql.Config{}
			if err := mysql.Load(configPath, &config); err != nil {
				return ctx, err
			}

			db, err := mysql.NewMySQLDB(config)
			if err != nil {
				return ctx, err
			}
			ctx = context.WithValue(ctx, contextKey("db"), db)

			client, err := chain.NewClient(ctx, c.String("node_endpoint"), "")
			if err != nil {
				return ctx, err
			}

			ctx = context.WithValue(ctx, contextKey("node_endpoint"), client)
			return ctx, nil
		},
		Commands: []*cli.Command{
			miner,
		},

		Action: func(ctx context.Context, c *cli.Command) error {
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
