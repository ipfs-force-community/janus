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
)

type Indexer struct {
	ctx         context.Context
	interval    int64
	node        *chain.Client
	db          *gorm.DB
	msgHandlers []chain.MsgHandler
}

func NewIndexer(ctx context.Context, interval int64, node *chain.Client, db *gorm.DB, msgHandlers ...chain.MsgHandler) *Indexer {
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
			return
		}
	}
}

func (i *Indexer) sync() error {
	latestHeight, err := i.localHeight()
	if err != nil {
		return err
	}

	// get the current chain head height
	headHeight, err := i.node.ChainHeadHeight()
	if err != nil {
		return err
	}

	// only sync up to headHeight - safeConfirmNum to avoid chain reorg issues
	headHeight -= safeConfirmNum
	if latestHeight >= headHeight {
		return nil
	}

	if err := i.node.SyncBlocks(latestHeight+1, headHeight, func(blockMeta *chain.BlockMeta, msg *types.Message) error {
		for _, handle := range i.msgHandlers {
			if err := handle(blockMeta, msg); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// update the latest synced height in the database
	if err := i.db.Model(&orm.Chain{}).Where("id = 1").Update("height", headHeight).Error; err != nil {
		return err
	}

	return nil
}

func (i *Indexer) localHeight() (int64, error) {
	var latestChain orm.Chain
	if err := i.db.First(&latestChain).Error; err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	} else if err == gorm.ErrRecordNotFound {
		// If no record found, initialize the height to minFetchHeight
		if err := i.db.Create(&orm.Chain{
			Height: minFetchHeight,
		}).Error; err != nil {
			return 0, err
		}

		return minFetchHeight, nil
	}

	return latestChain.Height, nil
}
