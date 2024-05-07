package model

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type Tx = entities.Tx
type storageTxManager interface {
	BeginTransaction(ctx context.Context) (Tx, error)
}

type TxManager struct {
	storage storageTxManager
}

func NewTxManager(storage storageTxManager) *TxManager {
	return &TxManager{
		storage: storage,
	}
}

func (t *TxManager) BeginTransaction(ctx context.Context) (Tx, error) {
	return t.storage.BeginTransaction(ctx)
}
