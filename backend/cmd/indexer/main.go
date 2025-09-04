package main

import (
	"context"
	"log"
	"os"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/venus/venus-shared/actors/types"
	"github.com/urfave/cli/v3"

	"github.com/ipfs-force-community/janus/chain"
	"github.com/ipfs-force-community/janus/database/mysql"
	"github.com/ipfs-force-community/janus/database/orm"
	"github.com/ipfs-force-community/janus/indexer"
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
				Name:  "interval",
				Usage: "Interval in seconds between indexing runs",
				Value: 10,
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

	if err := db.AutoMigrate(&orm.Miner{}, &orm.Chain{}); err != nil {
		return err
	}

	client, err := chain.NewClient(ctx, c.String("node-endpoint"), c.String("node-token"))
	if err != nil {
		return err
	}

	createMinerMsgHandler := func(blockMeta *chain.BlockMeta, msg *types.Message) error {
		if msg.To == builtin.StoragePowerActorAddr && msg.Method == builtin.MethodsPower.CreateMiner {
			if err := db.Create(&orm.Miner{
				Height:    blockMeta.Height,
				Cid:       blockMeta.Cid.String(),
				Timestamp: blockMeta.Timestamp,
				MsgCid:    msg.Cid().String(),
				From:      msg.From.String(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}

	indexer.NewIndexer(ctx, c.Int64("interval"), client, db, createMinerMsgHandler).Start()
	return nil
}
