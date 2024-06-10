package syncconfig

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
	// if it's zero it finds the etrog upgrade block
	GenesisBlockNumber uint64 `mapstructure:"GenesisBlockNumber"`

	// SyncUpToBlock specifies the synchronization highest block to sync (latest, finalized, pending or safe). also can have an offset
	// example: safe/-10, then the safe block ten blocks before the safe block
	SyncUpToBlock string `jsonschema:"enum=latest,enum=safe, enum=pending, enum=finalized" mapstructure:"SyncBlockProtection"`
	// BlockFinality indicates which block is considered finalized (latest, finalized, pending or safe) also can have an offset
	// example: safe/-10, then the safe block ten blocks before the finalized block
	// if not set assuming 'finalized'
	BlockFinality string `jsonschema:"enum=latest,enum=safe, enum=pending, enum=finalized" mapstructure:"SyncBlockFinalized"`

	// OverrideStorageCheck is a flag to override the storage check
	// take in account that without that check you can merge data from different rollups or differents L1 networks
	OverrideStorageCheck bool `mapstructure:"OverrideStorageCheck"`
}
