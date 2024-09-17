package rpcsync

import (
	jRPC "github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
)

func StartRPC(cfg Config, sync interface{}, state interface{}) {
	if !cfg.Enabled {
		log.Info("RPC server is disabled")
		return
	}
	cfgRPC := jRPC.Config{
		Port:                      cfg.Port,
		MaxRequestsPerIPAndSecond: cfg.MaxRequestsPerIPAndSecond,
	}

	serverCfg := buildRPCServiceEndPoints(sync, state)

	server := jRPC.NewServer(cfgRPC, serverCfg)
	log.Infof("RPC server started on port %d (MaxRequestsPerIPAndSecond=%f)", cfg.Port, cfg.MaxRequestsPerIPAndSecond)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal(err)
		}
	}()
}

func addSyncerEndpoint(sync interface{}, endPoints []jRPC.Service) []jRPC.Service {
	syncImpl, ok := sync.(SyncInterface)
	if !ok {
		log.Fatal("State must implement Sync")
	}
	endPoints = append(endPoints, jRPC.Service{
		Name: "sync",
		Service: &SyncEndpoints{
			Sync: syncImpl,
		},
	})
	return endPoints
}
