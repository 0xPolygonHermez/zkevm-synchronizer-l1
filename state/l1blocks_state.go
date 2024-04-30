package state

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool // The block is safe (have past the safe point, e.g. Finalized in L1)
	SyncVersion string
}

type StorageL1BlockInterface interface {
}

type L1BlocksState struct {
	storage StorageL1BlockInterface
}

func (b *L1BlocksState) AddBlock(ctx context.Context, block *L1Block, dbTx pgx.Tx) error {
	pgBlock
}

func (b *L1BlocksState) GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) (*state.L1Block, error) {
	return nil, nil
}
