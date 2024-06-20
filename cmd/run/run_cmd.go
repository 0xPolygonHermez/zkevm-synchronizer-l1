package run

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
	"github.com/urfave/cli/v2"
)

func RunCmd(cliCtx *cli.Context) error {
	config, err := config.Load(cliCtx)
	if err != nil {
		return err
	}
	log.Init(config.Log)
	log.Info("Running synchronizer")

	sync, err := synchronizer.NewSynchronizer(cliCtx.Context, *config)
	if err != nil {
		log.Error("Error creating synchronizer", err)
		return err
	}
	err = sync.Sync(true)
	if err != nil {
		log.Error("Error syncing", err)
	} else {
		log.Info("Syncing done")
	}
	return err
}
