package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	node, err := chain.NewNode(ctx, c.String("node-endpoint"), c.String("node-token"))
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
				Cost:      msg.Value.String(),
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}

	var wg sync.WaitGroup
	wg.Add(1)

	indexer := indexer.NewIndexer(ctx, c.Int64("interval"), node, db, createMinerMsgHandler)
	go func() {
		defer wg.Done()
		indexer.Start()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	slog.Info("received termination signal, initiating shutdown...")

	cancel()

	wg.Wait()

	if err := indexer.Close(); err != nil {
		return err
	}

	slog.Info("graceful shutdown complete")
	return nil
}
