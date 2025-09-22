package chain

import (
	"context"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin"
	"github.com/filecoin-project/venus/venus-shared/types"
)

func TestNewNode(t *testing.T) {
	t.Skip("Skipping test that requires a running lotus node")
	ctx := context.Background()
	node, err := NewNode(ctx, "http://127.0.0.1:3463", "")
	if err != nil {
		t.Fatal(err)
	}

	height := abi.ChainEpoch(5282426)
	tipset, err := node.ChainGetTipSetByHeight(ctx, height, types.TipSetKey{})
	if err != nil {
		t.Fatal(err)
	}

	for _, blkHeader := range tipset.Blocks() {
		messages, err := node.ChainGetBlockMessages(ctx, blkHeader.Cid())
		if err != nil {
			t.Fatal(err)
		}

		for _, msg := range messages.BlsMessages {
			if msg.To == builtin.StoragePowerActorAddr && msg.Method == builtin.MethodsPower.CreateMiner {
				fmt.Printf("CreateMiner message found: From %s -> To %s, Method %d\n",
					msg.From, msg.To, msg.Method)
			}
		}
	}
}
