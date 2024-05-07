package model

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

type storageTxType = entities.Tx
type stateTxType = entities.Tx

type Keyer interface {
	IsEqual(other interface{}) bool
	Key() uint64
}

// SetStorageHelper helper that add element, if already exists check that is
// the same
func SetStorageHelper[T Keyer](ctx context.Context, obj T, tx storageTxType,
	addFunc func(ctx context.Context, obj T, tx storageTxType) error,
	getFunc func(ctx context.Context, key uint64, tx storageTxType) (T, error),
) error {
	err := addFunc(ctx, obj, tx)
	if err != nil && errors.Is(err, entities.ErrAlreadyExists) {
		// Check if is the same batch on DB
		objRead, errRead := getFunc(ctx, obj.Key(), tx)
		if errRead != nil {
			return err
		}
		if obj.IsEqual(objRead) {
			return nil
		}

		return fmt.Errorf("element %d already exists but is different: %w", obj.Key(), err)
	}
	return err
}
