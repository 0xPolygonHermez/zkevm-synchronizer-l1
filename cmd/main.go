package main

import (
	"log"
	"os"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/cmd/run"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/cmd/version"
	"github.com/urfave/cli/v2"
)

const appName = "zkevm-sync-l1"

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = zkevm.Version

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
			Usage:   "Run",
			Action:  run.RunCmd,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
