//go:build !debug
// +build !debug

package rpcsync

import "github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"

// On release (no debug) there are no RPC server
func StartRPC(state interface{}) {
	log.Debug("RPC server is disabled on release")
}
