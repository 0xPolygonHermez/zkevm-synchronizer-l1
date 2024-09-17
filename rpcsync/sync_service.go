package rpcsync

import (
	"github.com/0xPolygon/cdk-rpc/rpc"
	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
)

type SyncInterface interface {
	// IsSynced returns true if the synchronizer is synced or false if it's not
	IsSynced() bool
}

type SyncEndpoints struct {
	Sync SyncInterface
}

// curl -X POST http://localhost:8025/ -H "Con -application/json" -d '{"method":"sync_version", "params":[], "id":1}'
func (b *SyncEndpoints) Version() (interface{}, rpc.Error) {
	return zkevm_synchronizer_l1.Version, nil
}

// curl -X POST http://localhost:1025/ -H "Con -application/json" -d '{"method":"sync_isSynced", "params":[], "id":1}'
func (b *SyncEndpoints) IsSynced() (interface{}, rpc.Error) {
	return b.Sync.IsSynced(), nil
}
