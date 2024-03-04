package version

import (
	"os"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/urfave/cli/v2"
)

func VersionCmd(*cli.Context) error {
	zkevm.PrintVersion(os.Stdout)
	return nil
}
