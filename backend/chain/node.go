package chain

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc"
	v1 "github.com/filecoin-project/venus/venus-shared/api/chain/v1"
)

// Node wraps the Filecoin full node API client
type Node struct {
	ctx context.Context
	v1.FullNode
	closer jsonrpc.ClientCloser
}

// NewNode creates a new Node instance
func NewNode(ctx context.Context, url string, token string) (*Node, error) {
	node, closer, err := v1.DialFullNodeRPC(ctx, url, token, nil)
	if err != nil {
		return nil, err
	}

	return &Node{
		ctx:      ctx,
		FullNode: node,
		closer:   closer,
	}, nil
}

// ChainHeadHeight returns the current chain head height
func (n *Node) ChainHeadHeight() (int64, error) {
	head, err := n.FullNode.ChainHead(n.ctx)
	if err != nil {
		return 0, err
	}

	return int64(head.Height()), nil
}

// Close closes the client connection
func (n *Node) Close() {
	n.closer()
}
