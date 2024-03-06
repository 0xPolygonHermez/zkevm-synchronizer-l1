package etherman

import "github.com/ethereum/go-ethereum/common"

// Config represents the configuration of the etherman
type Config struct {
	L1URL           string         `mapstructure:"L1URL"`
	Contracts       ContractConfig `mapstructure:"contracts"`
	ForkIDChunkSize uint64         `mapstructure:"forkIDChunkSize"`
	L1ChainID       uint64         `mapstructure:"l1ChainID"`
}

type ContractConfig struct {
	GlobalExitRootManagerAddr common.Address
	RollupManagerAddr         common.Address
	ZkEVMAddr                 common.Address
	PolAddr                   common.Address
}
