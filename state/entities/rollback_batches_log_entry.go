package entities

import (
	"encoding/binary"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

type RollbackBatchesLogEntry struct {
	BlockNumber           uint64
	LastBatchNumber       uint64
	LastBatchAccInputHash common.Hash
	L1EventAt             time.Time
	ReceivedAt            time.Time
	UndoFirstBlockNumber  uint64
	Description           string
	SequencesDeleted      []SequencedBatches
	syncVersion           *string
}

func (r *RollbackBatchesLogEntry) ID() common.Hash {
	var res [32]byte
	hash := sha3.NewLegacyKeccak256()
	blockNumberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(blockNumberBytes, r.BlockNumber)
	hash.Write(blockNumberBytes)
	lastBatchNumberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastBatchNumberBytes, r.LastBatchNumber)
	hash.Write(lastBatchNumberBytes)
	hash.Write(r.LastBatchAccInputHash.Bytes())
	copy(res[:], hash.Sum(nil))
	return res
}

func (r *RollbackBatchesLogEntry) SetSyncVersion(syncVersion string) {
	r.syncVersion = &syncVersion
}

func (r *RollbackBatchesLogEntry) SyncVersion() string {
	return zkevm_synchronizer_l1.Version
}
