package synchronizer

import (
	"context"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthermanInterface contains the methods required to interact with ethereum.
type EthermanInterface interface {
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	GetRollupInfoByBlockRange(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]etherman.Block, map[common.Hash][]etherman.Order, error)
	EthBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error)
	//GetNetworkID(ctx context.Context) (uint, error)
	GetRollupID() uint
	GetL1BlockUpgradeLxLy(ctx context.Context, genesisBlock *uint64) (uint64, error)
	GetForks(ctx context.Context, genBlockNumber uint64, lastL1BlockSynced uint64) ([]etherman.ForkIDInterval, error)
	GetFinalizedBlockNumber(ctx context.Context) (uint64, error)
}
