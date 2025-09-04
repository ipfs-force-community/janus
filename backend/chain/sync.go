package chain

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/venus/venus-shared/types"
	"github.com/ipfs/go-cid"
	"golang.org/x/sync/errgroup"
)

const batchBlockNum = 1000

// BlockMeta contains base info about a block
type BlockMeta struct {
	Height    int64
	Cid       cid.Cid
	Timestamp int64
}

// MsgHandler defines the function type for handling messages during block synchronization
type MsgHandler func(blockMeta *BlockMeta, msg *types.Message) error

// SyncBlocks synchronizes blocks from startEpoch to endEpoch and processes messages using the provided MsgHandler
func (c *Client) SyncBlocks(startEpoch, endEpoch int64, msgHandler MsgHandler) error {
	if startEpoch < 0 {
		return errors.New("startEpoch must be greater than 0")
	}

	headHeight, err := c.ChainHeadHeight()
	if err != nil {
		return err
	}

	slog.Info("chain head height", slog.Int64("height", headHeight))
	if endEpoch == 0 || endEpoch > headHeight {
		endEpoch = headHeight
	}

	if startEpoch > endEpoch {
		return errors.New("startEpoch must be less than or equal to endEpoch")
	}

	slog.Info("start syncing", slog.Int64("startEpoch", startEpoch), slog.Int64("endEpoch", endEpoch))

	// batch download blocks
	for endEpoch-startEpoch > batchBlockNum {
		slog.Info("syncing batch", slog.Int64("startEpoch", startEpoch), slog.Int64("endEpoch", startEpoch+batchBlockNum))
		if err := c.syncBatch(startEpoch, startEpoch+batchBlockNum, msgHandler); err != nil {
			return err
		}

		startEpoch += batchBlockNum
	}

	return c.syncBatch(startEpoch, endEpoch, msgHandler)
}

func (c *Client) syncBatch(startEpoch, endEpoch int64, msgHandler MsgHandler) error {
	g, ctx := errgroup.WithContext(c.ctx)
	for epoch := startEpoch; epoch < endEpoch; epoch++ {
		g.Go(func() error {
			msgHandle := make(map[cid.Cid]struct{})
			tipset, err := c.ChainGetTipSetByHeight(ctx, abi.ChainEpoch(epoch), types.TipSetKey{})
			if err != nil {
				return fmt.Errorf("failed to get tipset at epoch %d: %w", epoch, err)
			}

			for _, blkHeader := range tipset.Blocks() {
				messages, err := c.ChainGetBlockMessages(ctx, blkHeader.Cid())
				if err != nil {
					return fmt.Errorf("failed to get block messages for block %s: %w", blkHeader.Cid(), err)
				}

				for _, msg := range messages.BlsMessages {
					if _, exists := msgHandle[msg.Cid()]; exists {
						continue
					}

					if err := msgHandler(&BlockMeta{
						Height:    int64(blkHeader.Height),
						Cid:       blkHeader.Cid(),
						Timestamp: int64(blkHeader.Timestamp),
					}, msg); err != nil {
						return err
					}

					msgHandle[msg.Cid()] = struct{}{}
				}
			}
			return nil
		})
	}

	// Wait for all goroutines to complete and return any error
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
