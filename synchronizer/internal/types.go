package internal

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
)

type stateTxType = entities.Tx

type ReorgExecutionResult struct {
	// FirstL1BlockNumberValidAfterReorg is the first block or nil if the reorg have delete all blocks
	FirstL1BlockNumberValidAfterReorg *uint64
	ReasonError                       error
}

type RollbackBatchesData struct {
	LastBatchNumber       uint64
	LastBatchAccInputHash common.Hash
}
