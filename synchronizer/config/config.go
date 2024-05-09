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

	// SyncBlockProtection specify the state to sync (lastest, finalized, pending or safe)
	SyncBlockProtection string `jsonschema:"enum=lastest,enum=safe, enum=pending, enum=finalized" mapstructure:"SyncBlockProtection"`
	// Example: SyncBlockProtection= finalized, L1SafeBlockOffset= -10, then the safe block ten blocks before the finalized block
	SyncBlockProtectionOffset int `mapstructure:"SyncBlockProtectionOffset"`
}
