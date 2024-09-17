//go:build !debug
// +build !debug

package rpcsync

import (
	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
)

func buildRPCServiceEndPoints(sync interface{}, _ interface{}) []jRPC.Service {
	endPoints := []jRPC.Service{}
	endPoints = addSyncerEndpoint(sync, endPoints)
	return endPoints
}
