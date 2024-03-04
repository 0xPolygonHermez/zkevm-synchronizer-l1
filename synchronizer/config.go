package synchronizer

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config/types"
)

// Config represents the configuration of the synchronizer
type Config struct {
	// SyncInterval is the delay interval between reading new rollup information
	SyncInterval types.Duration `mapstructure:"SyncInterval"`

	// SyncChunkSize is the number of blocks to sync on each chunk
	SyncChunkSize uint64 `mapstructure:"SyncChunkSize"`

	// GenesisBlockNumber is the block number of the genesis block (first block to synchronize)
	GenesisBlockNumber uint64 `mapstructure:"GenesisBlockNumber"`
}
