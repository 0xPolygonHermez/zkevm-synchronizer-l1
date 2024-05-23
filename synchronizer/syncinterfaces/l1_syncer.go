package syncinterfaces

import "context"

type L1Syncer interface {
	SyncBlocks(ctx context.Context, lastEthBlockSynced *stateL1BlockType) (*stateL1BlockType, bool, error)
}
