//go:build debug
// +build debug

package rpcsync

import (
	"log"

	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
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

func StartRPC(state interface{}) {
	cfg := jRPC.Config{
		Port:                      1025,
		MaxRequestsPerIPAndSecond: 1000,
	}
	server := jRPC.NewServer(cfg, []jRPC.Service{
		{
			Name: "debug",
			Service: &DebugEndpoints{
				State: state.(StateDebugInterface),
			},
		},
	})
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()
}
