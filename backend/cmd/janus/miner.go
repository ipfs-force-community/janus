package main

import (
	"context"

	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/venus/venus-shared/actors/types"
	"github.com/urfave/cli/v3"
	"gorm.io/gorm"

	"github.com/ipfs-force-community/janus/chain"
	"github.com/ipfs-force-community/janus/database/orm"
)

var miner = &cli.Command{
	Name:     "miner",
	Usage:    "Sync Filecoin chain data and serve API for visualization",
	Commands: []*cli.Command{},
	Action:   minerAction,
}

func minerAction(ctx context.Context, c *cli.Command) error {
	client := ctx.Value(contextKey("node_endpoint")).(*chain.Client)
	db := ctx.Value(contextKey("db")).(*gorm.DB)

	if err := client.SyncBlocks(c.Int64("start_epoch"), c.Int64("end_epoch"), func(blockMeta *chain.BlockMeta, msg *types.Message) error {
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
	}); err != nil {
		return err
	}

	return nil
}
