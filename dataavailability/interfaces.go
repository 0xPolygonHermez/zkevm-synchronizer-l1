package dataavailability

import (
	"context"
	"math/big"

	jsonrpcclienttypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/jsonrpcclient/types"
	"github.com/ethereum/go-ethereum/common"
)

// DABackender is an interface for components that store and retrieve batch data
type DABackender interface {
	SequenceRetriever
	SequenceSender
	// Init initializes the DABackend
	Init() error
}

// SequenceSender is used to send provided sequence of batches
type SequenceSender interface {
	// PostSequence sends the sequence data to the data availability backend, and returns the dataAvailabilityMessage
	// as expected by the contract
	PostSequence(ctx context.Context, batchesData [][]byte) ([]byte, error)
}

// SequenceRetriever is used to retrieve batch data
type SequenceRetriever interface {
	// GetSequence retrieves the sequence data from the data availability backend
	GetSequence(ctx context.Context, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([][]byte, error)
}

type BatchL2Data struct {
	Data   []byte
	Source DataSourcePriority
}

// BatchDataProvider is used to retrieve batch data
type BatchDataProvider interface {
	// GetBatchL2Data retrieve the data of a batch from the DA backend. The returned data must be the pre-image of the hash
	GetBatchL2Data(batchNum []uint64, batchHashes []common.Hash, dataAvailabilityMessage []byte) ([]BatchL2Data, error)
}

// DataManager is an interface for components that send and retrieve batch data
type DataManager interface {
	BatchDataProvider
	SequenceSender
}

// ZKEVMClientTrustedBatchesGetter contains the methods required to interact with zkEVM-RPC
type ZKEVMClientTrustedBatchesGetter interface {
	BatchByNumber(ctx context.Context, number *big.Int) (*jsonrpcclienttypes.Batch, error)
	BatchesByNumbers(ctx context.Context, numbers []*big.Int) ([]*jsonrpcclienttypes.BatchData, error)
	ForcedBatchesByNumbers(ctx context.Context, numbers []*big.Int) ([]*jsonrpcclienttypes.BatchData, error)
}
