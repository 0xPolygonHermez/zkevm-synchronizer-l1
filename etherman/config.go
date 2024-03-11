package etherman

import "github.com/ethereum/go-ethereum/common"

// Config represents the configuration of the etherman
type Config struct {
	L1URL                string         `mapstructure:"L1URL"`
	Contracts            ContractConfig `mapstructure:"Contracts"`
	ForkIDChunkSize      uint64         `mapstructure:"ForkIDChunkSize"`
	L1ChainID            uint64         `mapstructure:"L1ChainID"`
	PararellBlockRequest bool           `mapstructure:"pararellBlockRequest"`
}

type ContractConfig struct {
	GlobalExitRootManagerAddr common.Address
	RollupManagerAddr         common.Address
	ZkEVMAddr                 common.Address
}
