package chain

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
)

// Client wraps the Filecoin full node API client
type Client struct {
	ctx context.Context
	v1.FullNode
	closer jsonrpc.ClientCloser
}

// NewClient creates a new Filecoin full node API client
func NewClient(ctx context.Context, url string, token string) (*Client, error) {
	node, closer, err := v1.DialFullNodeRPC(ctx, url, token, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:      ctx,
		FullNode: node,
		closer:   closer,
	}, nil
}

// ChainHeadHeight returns the current chain head height
func (c *Client) ChainHeadHeight() (int64, error) {
	head, err := c.FullNode.ChainHead(c.ctx)
	if err != nil {
		return 0, err
	}

	return int64(head.Height()), nil
}
