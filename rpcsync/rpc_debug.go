//go:build debug
// +build debug

package rpcsync

import (
	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
)

/*
To activate this you need the build tag `debug`:
- Example of launch.json entry
"name": "run CARDONA",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "buildFlags": "-tags 'debug'",
            "program": "cmd/main.go",


This RPC is deployed on port 1025:
- To execute rollbackBatches: lastBatch, accInputHash, l1BlockNumber
curl -X POST http://localhost:1025/ -H "Con -application/json" -d '{"method":"debug_rollbackBatches", "params":[53992, "0x1234", 5159957], "id":1}'
- To executge forceReorg: firstL1BlockNumberToKeep
curl -X POST http://localhost:1025/ -H "Con -application/json" -d '{"method":"debug_forceReorg", "params":[5159956], "id":1}'
*/

func buildRPCServiceEndPoints(sync interface{}, state interface{}) []jRPC.Service {
	stateImpl, ok := state.(StateDebugInterface)
	if !ok {
		log.Fatal("State must implement StateDebugInterface")
	}
	endPoints := []jRPC.Service{
		{
			Name: "debug",
			Service: &DebugEndpoints{
				State: stateImpl,
			},
		},
	}
	endPoints = addSyncerEndpoint(sync, endPoints)
	return endPoints
}
