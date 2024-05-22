package entities

import "time"

type KVKey = string
type KVMetadataEntry struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SyncVersion string
}
