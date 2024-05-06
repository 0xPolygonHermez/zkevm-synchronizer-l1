package synchronizer

import (
	"context"
)

type txProvider interface {
	BeginTransaction(ctx context.Context) (stateTxType, error)
}

type ReuseStateTx struct {
	stateTxType
}

type ReuseStateTxProvider struct {
	isReusingTx    bool
	forcedTx       stateTxType
	realTxProvider txProvider
}

func NewReuseStateTxProvider(realTxProvider txProvider, forcedTx stateTxType) *ReuseStateTxProvider {
	return &ReuseStateTxProvider{
		isReusingTx:    true,
		forcedTx:       forcedTx,
		realTxProvider: realTxProvider,
	}
}

func (r *ReuseStateTxProvider) BeginTransaction(ctx context.Context) (stateTxType, error) {
	if r.isReusingTx {
		return r.forcedTx, nil
	} else {
		return r.realTxProvider.BeginTransaction(ctx)
	}
}

func (r *ReuseStateTx) Commit(ctx context.Context) error {
	return nil
}

func (r *ReuseStateTx) Rollback(ctx context.Context) error {
	return nil
}
