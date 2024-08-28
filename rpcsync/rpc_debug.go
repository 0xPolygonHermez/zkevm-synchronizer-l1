//go:build debug
// +build debug

package rpcsync

import (
	"log"

	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
)

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
