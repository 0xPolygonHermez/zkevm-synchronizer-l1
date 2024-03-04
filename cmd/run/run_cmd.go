package run

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config/types"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
)

func getConfigForInternalLocalDockers() (db.Config, synchronizer.Config, etherman.Config) {
	configDB := db.Config{
		Name:     "state_db",
		User:     "test_user",
		Password: "test_password",
		Host:     "localhost",
		Port:     "5434",
		MaxConns: 10,
	}
	configSync := synchronizer.Config{
		SyncInterval:       types.NewDuration(time.Duration(10)),
		SyncChunkSize:      1000,
		GenesisBlockNumber: 4794475,
	}
	configEtherman := etherman.Config{
		L1URL:                             "https://sepolia.infura.io/v3/79f68fdb480a422886a39053af20cea7",
		PolygonZkEVMGlobalExitRootAddress: common.HexToAddress("0x2968D6d736178f8FE7393CC33C87f29D9C287e78"),
		PolygonRollupManagerAddress:       common.HexToAddress("0xE2EF6215aDc132Df6913C8DD16487aBF118d1764"),
		PolygonZkEvmAddress:               common.HexToAddress("0x89BA0Ed947a88fe43c22Ae305C0713eC8a7Eb361"),
	}
	return configDB, configSync, configEtherman
}

func getConfigForLocalDockers() (db.Config, synchronizer.Config, etherman.Config) {
	configDB := db.Config{
		Name:     "state_db",
		User:     "test_user",
		Password: "test_password",
		Host:     "localhost",
		Port:     "5434",
		MaxConns: 10,
	}
	configSync := synchronizer.Config{
		SyncInterval:       types.NewDuration(time.Duration(10)),
		SyncChunkSize:      1000,
		GenesisBlockNumber: 1,
	}
	configEtherman := etherman.Config{
		L1URL:                             "http://localhost:8545",
		PolygonZkEVMGlobalExitRootAddress: common.HexToAddress("0x8A791620dd6260079BF849Dc5567aDC3F2FdC318"),
		PolygonRollupManagerAddress:       common.HexToAddress("0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e"),
		PolygonZkEvmAddress:               common.HexToAddress("0x8dAF17A20c9DBA35f005b6324F493785D239719d"),
	}
	return configDB, configSync, configEtherman
}

func RunCmd(cliCtx *cli.Context) error {
	log.Info("Running synchronizer")
	configDB, configSync, configEtherman := getConfigForInternalLocalDockers()
	//setupLog(c.Log)
	err := db.RunMigrations(configDB)
	if err != nil {
		log.Error(err)
		return err
	}
	configStorage := pgstorage.Config{
		Name:     configDB.Name,
		User:     configDB.User,
		Password: configDB.Password,
		Host:     configDB.Host,
		Port:     configDB.Port,
		MaxConns: configDB.MaxConns,
	}
	storage, err := pgstorage.NewPostgresStorage(configStorage)
	if err != nil {
		log.Error(err)
		return err
	}
	etherman, err := etherman.NewClient(configEtherman)
	if err != nil {
		log.Error("Error creating etherman", err)
		return err
	}
	sync, err := synchronizer.NewSynchronizer(cliCtx.Context, storage, etherman, configSync)
	if err != nil {
		log.Error("Error creating synchronizer", err)
		return err
	}
	return sync.Sync()
}
