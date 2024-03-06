package version

import (
	"os"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/urfave/cli/v2"
)

func VersionCmd(*cli.Context) error {
	zkevm_synchronizer_l1.PrintVersion(os.Stdout)
	return nil
}
