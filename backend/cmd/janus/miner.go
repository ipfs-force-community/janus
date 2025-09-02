package main

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/venus/venus-shared/actors/types"
	"github.com/urfave/cli/v3"

	"github.com/ipfs-force-community/janus/chain"
)

var miner = &cli.Command{
	Name:     "miner",
	Usage:    "Sync Filecoin chain data and serve API for visualization",
	Commands: []*cli.Command{},
	Action: func(ctx context.Context, c *cli.Command) error {
		client := ctx.Value(contextKey("node_endpoint")).(*chain.Client)
		if err := client.SyncBlocks(c.Int64("start_epoch"), c.Int64("end_epoch"), func(epoch abi.ChainEpoch, msg *types.Message) error {
			if msg.To == builtin.StoragePowerActorAddr && msg.Method == builtin.MethodsPower.CreateMiner {
				fmt.Printf("CreateMiner message found: Msg ID: %s From %s -> To %s, Method %d\n", msg.Cid().String(),
					msg.From, msg.To, msg.Method)
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	},
}
