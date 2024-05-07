package entities

const (
	// FORKID_ZERO is the fork id 0 (no forkid)
	FORKID_ZERO = uint64(0)
	// FORKID_BLUEBERRY is the fork id 4
	FORKID_BLUEBERRY = uint64(4)
	// FORKID_DRAGONFRUIT is the fork id 5
	FORKID_DRAGONFRUIT = uint64(5)
	// FORKID_INCABERRY is the fork id 6
	FORKID_INCABERRY = uint64(6)
	// FORKID_ETROG is the fork id 7
	FORKID_ETROG = uint64(7)
	// FORKID_ELDERBERRY is the fork id 8
	FORKID_ELDERBERRY = uint64(8)
)

// ForkIDInterval is a fork id interval
type ForkIDInterval struct {
	FromBatchNumber uint64
	ToBatchNumber   uint64
	ForkId          uint64
	Version         string
	BlockNumber     uint64
}
