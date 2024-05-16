package synchronizer

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
)

// object: ReuseStateTxProvider
// intent:
//   for a process that create a new tx each time
//   you can bypass this behavious returning a 'fake tx' ReuseStateTx
//   that ignore commit. The rollback is real
//   The usage implies that the level that use this
//   will commit / rollback the real Tx
// It substitute the state class that implements
//    BeginTransaction

type txProvider interface {
	BeginTransaction(ctx context.Context) (stateTxType, error)
}

type ReuseStateTx struct {
	stateTxType
	hasRollback bool
}

type ReuseStateTxProvider struct {
	forcedTx       *ReuseStateTx
	realTxProvider txProvider
}

func NewReuseStateTxProvider(realTxProvider txProvider, forcedTx stateTxType) *ReuseStateTxProvider {
	return &ReuseStateTxProvider{
		forcedTx:       &ReuseStateTx{forcedTx, false},
		realTxProvider: realTxProvider,
	}
}

func (r *ReuseStateTxProvider) BeginTransaction(ctx context.Context) (stateTxType, error) {
	return r.forcedTx, nil
}

func (r *ReuseStateTx) Commit(ctx context.Context) error {
	if r.hasRollback {
		return fmt.Errorf("this transaction have been rollbacked")
	}
	log.Debugf("Ignoring commit because resuing tx")
	return nil
}

func (r *ReuseStateTx) Rollback(ctx context.Context) error {
	r.hasRollback = true
	return r.stateTxType.Rollback(ctx)
}
