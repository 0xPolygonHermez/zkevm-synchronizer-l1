package l1_check_block

import (
	"context"
	"fmt"
)

type IterationResult struct {
	Err           error
	ReorgDetected bool
	BlockNumber   uint64
	ReorgMessage  string
}

func (ir *IterationResult) String() string {
	if ir.Err == nil {
		if ir.ReorgDetected {
			return fmt.Sprintf("IterationResult{ReorgDetected: %v, BlockNumber: %d ReorgMessage:%s}", ir.ReorgDetected, ir.BlockNumber, ir.ReorgMessage)
		} else {
			return "IterationResult{None}"
		}
	} else {
		return fmt.Sprintf("IterationResult{Err: %s, ReorgDetected: %v, BlockNumber: %d ReorgMessage:%s}", ir.Err.Error(), ir.ReorgDetected, ir.BlockNumber, ir.ReorgMessage)
	}
}

type AsyncL1BlockChecker interface {
	Run(ctx context.Context, onFinish func())
	RunSynchronous(ctx context.Context) IterationResult
	Stop()
	GetResult() *IterationResult
}

type L1BlockCheckerIntegrator interface {
	OnStart(ctx context.Context) error
	OnResetState(ctx context.Context)
	CheckReorgWrapper(ctx context.Context, reorgFirstBlockOk *L1Block, badBlockNumber uint64, errReportedByReorgFunc error) (*L1Block, uint64, error)
}
