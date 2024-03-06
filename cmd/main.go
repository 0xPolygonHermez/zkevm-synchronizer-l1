package main

import (
	"log"
	"os"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/cmd/run"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/cmd/version"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/config"
	"github.com/urfave/cli/v2"
)

const appName = "zkevm-sync-l1"

var (
	configFileFlag = cli.StringFlag{
		Name:     config.FlagCfg,
		Aliases:  []string{"c"},
		Usage:    "Configuration `FILE`",
		Required: true,
	}
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = zkevm_synchronizer_l1.Version

	app.Commands = []*cli.Command{
		{
			Name:    "version",
			Aliases: []string{},
			Usage:   "Application version and build",
			Action:  version.VersionCmd,
		},
		{
			Name:    "run",
			Aliases: []string{},
			Usage:   "Run synchronizer as standalone",
			Action:  run.RunCmd,
			Flags:   []cli.Flag{&configFileFlag},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
