package zkevm_synchronizer_l1

import (
	"fmt"
	"io"
	"runtime"
)

// v0.2.0 - Add GetLeafsByL1InfoRoot
// v0.2.1 - Add L1InfoRoot to sequenced batches
var (
	Version = "v0.2.1"
)

// PrintVersion prints version info into the provided io.Writer.
func PrintVersion(w io.Writer) {
	fmt.Fprintf(w, "Version:      %s\n", Version)
	fmt.Fprintf(w, "Go version:   %s\n", runtime.Version())
	fmt.Fprintf(w, "OS/Arch:      %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
