package rpcsync

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygon/cdk-rpc/rpc"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/model"
	"github.com/ethereum/go-ethereum/common"
)

type StateDebugInterface interface {
	ExecuteRollbackBatches(ctx context.Context, rollbackBatchesRequest model.RollbackBatchesRequest, dbTx entities.Tx) (*model.RollbackBatchesExecutionResult, error)
	ExecuteReorg(ctx context.Context, reorgRequest model.ReorgRequest, dbTx entities.Tx) model.ReorgExecutionResult
	BeginTransaction(ctx context.Context) (entities.Tx, error)
}

type DebugEndpoints struct {
	State StateDebugInterface
}

func (b *DebugEndpoints) Hello() (interface{}, rpc.Error) {
	return "Hello", nil
}

func (b *DebugEndpoints) RollbackBatches(lastBatchNumber uint64, accInputHash string, l1BlockNumber uint64) (interface{}, rpc.Error) {
	req := model.RollbackBatchesRequest{
		LastBatchNumber:       lastBatchNumber,
		LastBatchAccInputHash: common.HexToHash(accInputHash),
		L1BlockNumber:         l1BlockNumber,
		L1BlockTimestamp:      time.Now(),
	}
	ctx := context.Background()
	dbTx, err := b.State.BeginTransaction(ctx)
	if err != nil {
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, err.Error())
	}
	log.Warnf("RPC: Execute fake RollbackBatches %v", req)
	res, err := b.State.ExecuteRollbackBatches(ctx, req, dbTx)
	if err != nil {
		errRollback := dbTx.Rollback(ctx)
		log.Warnf("RPC: RollbackBatches fails %v. Rollback DB result: %v", err, errRollback)
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, err.Error())
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Warnf("RPC: RollbackBatches ok bu Commit DB fails: %v", err)
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, err.Error())
	}
	return *res, nil
}

func (b *DebugEndpoints) ForceReorg(firstL1BlockNumberToKeep uint64) (interface{}, rpc.Error) {
	req := model.ReorgRequest{
		FirstL1BlockNumberToKeep: firstL1BlockNumberToKeep,
		ReasonError:              fmt.Errorf("forced reorg by RPC"),
	}
	ctx := context.Background()
	dbTx, err := b.State.BeginTransaction(ctx)
	if err != nil {
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, err.Error())
	}
	log.Warnf("RPC: Execute fake ExecuteReorg %v", req)
	res := b.State.ExecuteReorg(ctx, req, dbTx)
	if res.ExecutionError != nil {
		errRollback := dbTx.Rollback(ctx)
		if errRollback != nil {
			log.Warnf("RPC: ExecuteReorg fails %v. Rollback DB result: %v", res.ExecutionError, errRollback)
		}
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, res.ExecutionError.Error())
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		log.Warnf("RPC: ExecuteReorg ok bu Commit DB fails: %v", err)
		return nil, rpc.NewRPCError(rpc.DefaultErrorCode, err.Error())
	}
	return res, nil
}
