package syncconfig

import (
	"fmt"

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

	// L1BlockCheck is the configuration for the L1 Block Checker
	L1BlockCheck L1BlockCheckConfig `mapstructure:"L1BlockCheck"`
}

// L1BlockCheckConfig Configuration for L1 Block Checker
type L1BlockCheckConfig struct {
	// Enable if is true then the check l1 Block Hash is active
	Enable bool `mapstructure:"Enable"`
	// L1SafeBlockPoint is the point that a block is considered safe enough to be checked
	// it can be: finalized, safe,pending or latest
	L1SafeBlockPoint string `mapstructure:"L1SafeBlockPoint" jsonschema:"enum=finalized,enum=safe, enum=pending,enum=latest"`
	// L1SafeBlockOffset is the offset to add to L1SafeBlockPoint as a safe point
	// it can be positive or negative
	// Example: L1SafeBlockPoint= finalized, L1SafeBlockOffset= -10, then the safe block ten blocks before the finalized block
	L1SafeBlockOffset int `mapstructure:"L1SafeBlockOffset"`
	// ForceCheckBeforeStart if is true then the first time the system is started it will force to check all pending blocks
	ForceCheckBeforeStart bool `mapstructure:"ForceCheckBeforeStart"`

	// PreCheckEnable if is true then the pre-check is active, will check blocks between L1SafeBlock and L1PreSafeBlock
	PreCheckEnable bool `mapstructure:"PreCheckEnable"`
	// L1PreSafeBlockPoint is the point that a block is considered safe enough to be checked
	// it can be: finalized, safe,pending or latest
	L1PreSafeBlockPoint string `mapstructure:"L1PreSafeBlockPoint" jsonschema:"enum=finalized,enum=safe, enum=pending,enum=latest"`
	// L1PreSafeBlockOffset is the offset to add to L1PreSafeBlockPoint as a safe point
	// it can be positive or negative
	// Example: L1PreSafeBlockPoint= finalized, L1PreSafeBlockOffset= -10, then the safe block ten blocks before the finalized block
	L1PreSafeBlockOffset int `mapstructure:"L1PreSafeBlockOffset"`
}

func (c *L1BlockCheckConfig) String() string {
	return fmt.Sprintf("Enable: %v, L1SafeBlockPoint: %s, L1SafeBlockOffset: %d, ForceCheckBeforeStart: %v", c.Enable, c.L1SafeBlockPoint, c.L1SafeBlockOffset, c.ForceCheckBeforeStart)
}
