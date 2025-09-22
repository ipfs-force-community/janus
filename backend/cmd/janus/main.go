package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ipfs-force-community/janus/chain"
	"github.com/ipfs-force-community/janus/database/mysql"
	"github.com/ipfs-force-community/janus/database/orm"
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
				Name:    "node-endpoint",
				Aliases: []string{"n"},
				Usage:   "Filecoin node endpoint",
				Value:   "http://127.0.0.1:3463",
			},
			&cli.StringFlag{
				Name:     "node-token",
				Usage:    "Filecoin node endpoint token",
				Required: true,
			},
			&cli.Int64Flag{
				Name:  "start-epoch",
				Usage: "Start epoch to sync from",
				Value: 5260000,
			},
			&cli.Int64Flag{
				Name:  "end-epoch",
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

			if err := db.AutoMigrate(&orm.Miner{}); err != nil {
				return ctx, err
			}

			ctx = context.WithValue(ctx, contextKey("db"), db)

			node, err := chain.NewNode(ctx, c.String("node-endpoint"), c.String("node-token"))
			if err != nil {
				return ctx, err
			}

			ctx = context.WithValue(ctx, contextKey("node_endpoint"), node)
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
