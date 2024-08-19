package etherman

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/translator"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
	"github.com/ethereum/go-ethereum/common"
)

// Config represents the configuration of the etherman
type Config struct {
	L1URL                string         `mapstructure:"L1URL"`
	ForkIDChunkSize      uint64         `mapstructure:"ForkIDChunkSize"`
	L1ChainID            uint64         `mapstructure:"L1ChainID"`
	PararellBlockRequest bool           `mapstructure:"pararellBlockRequest"`
	Contracts            ContractConfig `mapstructure:"Contracts"`
	Validium             ValidiumConfig `mapstructure:"Validium"`
}

type ValidiumConfig struct {
	Enabled             bool   `mapstructure:"Enabled"`
	TrustedSequencerURL string `mapstructure:"TrustedSequencerURL"`
	// DataSourcePriority defines the order in which L2 batch should be retrieved: local, trusted, external
	DataSourcePriority      []dataavailability.DataSourcePriority `mapstructure:"DataSourcePriority"`
	Translator              translator.Config
	RetryOnDACErrorInterval types.Duration `mapstructure:"RetryOnDACErrorInterval"`
	RateLimit               utils.RateLimitConfig
}

type ContractConfig struct {
	GlobalExitRootManagerAddr common.Address
	RollupManagerAddr         common.Address
	ZkEVMAddr                 common.Address
}
