package syncinterfaces

import "context"

type StorageCompatibilityChecker interface {
	CheckAndUpdateStorage(ctx context.Context) error
}
