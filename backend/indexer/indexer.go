package indexer

import (
	"context"
	"log/slog"
	"time"

	"github.com/filecoin-project/venus/venus-shared/actors/types"
	"gorm.io/gorm"

	"github.com/ipfs-force-community/janus/chain"
	"github.com/ipfs-force-community/janus/database/orm"
)

const (
	minFetchHeight = 5200000
	safeConfirmNum = 20
	batchBlockNum  = 1000
	maxRetries     = 3
)

type Indexer struct {
	ctx         context.Context
	interval    int64
	node        *chain.Node
	db          *gorm.DB
	msgHandlers []chain.MsgHandler
}

func NewIndexer(ctx context.Context, interval int64, node *chain.Node, db *gorm.DB, msgHandlers ...chain.MsgHandler) *Indexer {
	return &Indexer{
		ctx:         ctx,
		interval:    interval,
		node:        node,
		db:          db,
		msgHandlers: msgHandlers,
	}
}

func (i *Indexer) Start() {
	ticker := time.NewTicker(time.Duration(i.interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := i.sync(); err != nil {
				slog.Error("indexer sync error", "error", err)
			}

		case <-i.ctx.Done():
			slog.Info("indexer context done, exiting...")
			return
		}
	}
}

func (i *Indexer) sync() error {
	latestHeight, err := i.localHeight()
	if err != nil {
		return err
	}

	headHeight, err := i.node.ChainHeadHeight()
	if err != nil {
		return err
	}

	headHeight -= safeConfirmNum
	if latestHeight >= headHeight {
		return nil
	}

	currentHeight := latestHeight
	for currentHeight < headHeight {
		endHeight := currentHeight + batchBlockNum
		if endHeight > headHeight {
			endHeight = headHeight
		}

		var lastErr error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			if err := i.node.SyncBlocks(currentHeight+1, endHeight, func(blockMeta *chain.BlockMeta, msg *types.Message) error {
				for _, handle := range i.msgHandlers {
					if err := handle(blockMeta, msg); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				lastErr = err
				if attempt < maxRetries {
					slog.Warn("batch sync failed, retrying", "from", currentHeight+1, "to", endHeight, "attempt", attempt, "error", err)
				} else {
					slog.Warn("batch sync failed, skipping", "from", currentHeight+1, "to", endHeight, "error", err)
				}
			} else {
				lastErr = nil
				break
			}
		}

		if lastErr != nil {
			slog.Warn("skipping failed batch, moving to next", "from", currentHeight+1, "to", endHeight)
		}

		if err := i.updateHeight(endHeight); err != nil {
			slog.Error("failed to update height", "error", err)
			return err
		}

		currentHeight = endHeight
		if lastErr == nil {
			slog.Info("batch synced", "from", currentHeight+1, "to", endHeight)
		}
	}

	slog.Info("indexer sync completed", "from", latestHeight+1, "to", headHeight)
	return nil
}

func (i *Indexer) updateHeight(height int64) error {
	return i.db.Model(&orm.Chain{}).Where("id = 1").Update("height", height).Error
}

func (i *Indexer) localHeight() (int64, error) {
	var latestChain orm.Chain
	if err := i.db.First(&latestChain).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		if err := i.db.Create(&orm.Chain{
			Height: minFetchHeight,
		}).Error; err != nil {
			return 0, err
		}

		return minFetchHeight, nil
	}

	return latestChain.Height, nil
}

func (i *Indexer) Close() error {
	i.node.Close()

	db, err := i.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
