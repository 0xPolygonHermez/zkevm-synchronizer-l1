package entities

import "time"

type KVMetadataEntry struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SyncVersion string
}
