package etherman

import "github.com/ethereum/go-ethereum/common"

// Config represents the configuration of the etherman
type Config struct {
	L1URL  string   `mapstructure:"L1URL"`
	L2URLs []string `mapstructure:"L2URLs"`

	PolygonZkEVMGlobalExitRootAddress common.Address
	PolygonRollupManagerAddress       common.Address
	PolygonZkEvmAddress               common.Address
}
