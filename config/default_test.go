package config_test

import (
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	storage "github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage"
	syncconfig "github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/translator"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDefault(t *testing.T) {
	expectedCfg := &config.Config{
		Log: log.Config{
			Environment: "development",
			Level:       "info",
			Outputs:     []string{"stderr"},
		},
		SQLDB: storage.Config{
			DriverName: "set_driver. example: sqlite",
			DataSource: "example: file::memory:?cache=shared",
		},
		Synchronizer: syncconfig.Config{
			SyncInterval:         types.Duration{Duration: time.Second * 10},
			SyncChunkSize:        500,
			GenesisBlockNumber:   0,
			SyncUpToBlock:        "latest",
			BlockFinality:        "finalized",
			OverrideStorageCheck: false,
		},
		Etherman: etherman.Config{
			L1URL:           "http://localhost:8545",
			ForkIDChunkSize: 100,
			L1ChainID:       0,
			Contracts: etherman.ContractConfig{
				GlobalExitRootManagerAddr: common.HexToAddress("0x2968D6d736178f8FE7393CC33C87f29D9C287e78"),
				RollupManagerAddr:         common.HexToAddress("0xE2EF6215aDc132Df6913C8DD16487aBF118d1764"),
				ZkEVMAddr:                 common.HexToAddress("0x89BA0Ed947a88fe43c22Ae305C0713eC8a7Eb361"),
			},
			Validium: etherman.ValidiumConfig{
				Enabled:                 false,
				TrustedSequencerURL:     "",
				DataSourcePriority:      []dataavailability.DataSourcePriority{dataavailability.Trusted, dataavailability.External},
				Translator:              translator.Config{FullMatchRules: []translator.ConfigRuleFullMatch{}},
				RetryOnDACErrorInterval: types.Duration{Duration: time.Minute},
				RateLimit:               utils.NewRateLimitConfig(900, time.Second),
			},
		},
	}
	cfg, err := config.Default()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, expectedCfg, cfg)
}
