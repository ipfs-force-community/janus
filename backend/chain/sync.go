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
func (n *Node) SyncBlocks(startEpoch, endEpoch int64, msgHandler MsgHandler) error {
	if startEpoch < 0 {
		return errors.New("startEpoch must be greater than 0")
	}

	headHeight, err := n.ChainHeadHeight()
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
		if err := n.syncBatch(startEpoch, startEpoch+batchBlockNum, msgHandler); err != nil {
			return err
		}

		startEpoch += batchBlockNum
	}

	return n.syncBatch(startEpoch, endEpoch, msgHandler)
}

func (n *Node) syncBatch(startEpoch, endEpoch int64, handler MsgHandler) error {
	g, ctx := errgroup.WithContext(n.ctx)
	for epoch := startEpoch; epoch <= endEpoch; epoch++ {
		g.Go(func() error {
			tipset, err := n.ChainGetTipSetByHeight(ctx, abi.ChainEpoch(epoch), types.TipSetKey{})
			if err != nil {
				return fmt.Errorf("failed to get tipset at epoch %d: %w", epoch, err)
			}

			seen := make(map[cid.Cid]struct{})
			for _, blk := range tipset.Blocks() {
				msgs, err := n.ChainGetBlockMessages(ctx, blk.Cid())
				if err != nil {
					return fmt.Errorf("get messages for block %s: %w", blk.Cid(), err)
				}

				process := func(cmsg cid.Cid, m *types.Message) error {
					if _, ok := seen[cmsg]; ok {
						return nil
					}

					seen[cmsg] = struct{}{}
					return handler(&BlockMeta{
						Height:    int64(blk.Height),
						Cid:       blk.Cid(),
						Timestamp: int64(blk.Timestamp),
					}, m)
				}

				for _, m := range msgs.BlsMessages {
					if err := process(m.Cid(), m); err != nil {
						return err
					}
				}
				for _, sm := range msgs.SecpkMessages {
					if err := process(sm.Cid(), &sm.Message); err != nil {
						return err
					}
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
