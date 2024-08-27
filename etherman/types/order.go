package types

import "fmt"

// EventOrder is the the type used to identify the events order
type EventOrder string

const (
	// GlobalExitRootsOrder identifies a GlobalExitRoot event
	GlobalExitRootsOrder EventOrder = "GlobalExitRoots"
	// L1InfoTreeOrder identifies a L1InTree event
	L1InfoTreeOrder EventOrder = "L1InfoTreeOrder"
	// SequenceBatchesOrder identifies a VerifyBatch event
	SequenceBatchesOrder EventOrder = "SequenceBatches"
	// UpdateEtrogSequenceOrder identifies a VerifyBatch event
	UpdateEtrogSequenceOrder EventOrder = "UpdateEtrogSequence"
	// ForcedBatchesOrder identifies a ForcedBatches event
	ForcedBatchesOrder EventOrder = "ForcedBatches"
	// TrustedVerifyBatchOrder identifies a TrustedVerifyBatch event
	TrustedVerifyBatchOrder EventOrder = "TrustedVerifyBatch"
	// VerifyBatchOrder identifies a VerifyBatch event
	VerifyBatchOrder EventOrder = "VerifyBatch"
	// SequenceForceBatchesOrder identifies a SequenceForceBatches event
	SequenceForceBatchesOrder EventOrder = "SequenceForceBatches"
	// ForkIDsOrder identifies an updateZkevmVersion event
	ForkIDsOrder EventOrder = "forkIDs"
	// InitialSequenceBatchesOrder identifies a VerifyBatch event
	InitialSequenceBatchesOrder EventOrder = "InitialSequenceBatches"
	// UpdateL1InfoTreeOrder identifies a rollbackBatchesSignatureHash event
	UpdateL1InfoTreeV2Order EventOrder = "UpdateL1InfoTreeV2"
	// RollbackBatchesOrder identifies a rollbackBatchesManagerSignatureHash event
	RollbackBatchesOrder EventOrder = "RollbackBatches"
)

// Order contains the event order to let the synchronizer store the information following this order.
type Order struct {
	Name EventOrder
	Pos  int
}

func (o Order) String() string {
	return fmt.Sprintf("Name: %s, Pos: %d", o.Name, o.Pos)
}
