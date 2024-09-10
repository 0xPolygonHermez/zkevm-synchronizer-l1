package types

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types"
)

type BlockRetriever interface {
	RetrieveFullBlockForEvent(ctx context.Context, vLog types.Log) (*Block, error)
}
