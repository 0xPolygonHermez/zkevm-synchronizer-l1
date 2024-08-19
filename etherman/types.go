package etherman

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmetrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/oldpolygonzkevm"
	"github.com/ethereum/go-ethereum/common"
)

// Block struct
type Block struct {
	BlockNumber           uint64
	BlockHash             common.Hash
	ParentHash            common.Hash
	ForcedBatches         []ForcedBatch
	SequencedBatches      [][]SequencedBatch
	UpdateEtrogSequence   UpdateEtrogSequence
	VerifiedBatches       []VerifiedBatch
	SequencedForceBatches [][]SequencedForceBatch
	ForkIDs               []ForkID
	ReceivedAt            time.Time
	// GER data
	GlobalExitRoots, L1InfoTree []GlobalExitRoot
}

func (b *Block) HasEvents() bool {
	return len(b.ForcedBatches) > 0 || len(b.SequencedBatches) > 0 || b.UpdateEtrogSequence.BatchNumber > 0 ||
		len(b.VerifiedBatches) > 0 || len(b.SequencedForceBatches) > 0 || len(b.ForkIDs) > 0 || len(b.GlobalExitRoots) > 0 || len(b.L1InfoTree) > 0
}

// GlobalExitRoot struct
type GlobalExitRoot struct {
	BlockNumber       uint64
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
	Timestamp         time.Time
	PreviousBlockHash common.Hash
}

// SequencedBatchElderberryData represents an Elderberry sequenced batch data
type SequencedBatchElderberryData struct {
	MaxSequenceTimestamp     uint64
	InitSequencedBatchNumber uint64 // Last sequenced batch number
}

type SourceBatchDataEnum = string

const (
	SourceBatchDataCalldata           = "calldata"
	SourceBatchDataValidiumDAExternal = "DA/External"
	SourceBatchDataValidiumDATrusted  = "DA/Trusted"
)

type RollupFlavorEnum = string

const (
	RollupFlavorZkEVM    = "kEVM"
	RollupFlavorValidium = "Validium"
)

type SequencedBatchMetadata struct {
	// SourceBatchData
	SourceBatchData  SourceBatchDataEnum
	RollupFlavor     RollupFlavorEnum
	CallFunctionName string // Call function
	ForkName         string // Fork name (elderberry / etrog)

}

func (s *SequencedBatchMetadata) String() string {
	return fmt.Sprintf("SourceBatchData: %s RollupFlavor: %s CallFunctionName: %s ForkName: %s", s.SourceBatchData, s.RollupFlavor, s.CallFunctionName, s.ForkName)
}

// SequencedBatch represents virtual batch
type SequencedBatch struct {
	BatchNumber   uint64
	L1InfoRoot    *common.Hash
	SequencerAddr common.Address
	TxHash        common.Hash
	Nonce         uint64
	Coinbase      common.Address
	// Struct used in preEtrog forks
	*oldpolygonzkevm.PolygonZkEVMBatchData
	// Struct used in Etrog + Elderberry
	*polygonzkevmetrog.PolygonRollupBaseEtrogBatchData
	// Struct used in Elderberry
	*SequencedBatchElderberryData
	Metadata *SequencedBatchMetadata
}

func (s *SequencedBatch) String() string {
	res := fmt.Sprintf("BatchNumber: %d\n", s.BatchNumber)
	res += fmt.Sprintf("L1InfoRoot: %s\n", s.L1InfoRoot.String())
	res += fmt.Sprintf("SequencerAddr: %s\n", s.SequencerAddr.String())
	res += fmt.Sprintf("TxHash: %s\n", s.TxHash.String())
	res += fmt.Sprintf("Nonce: %d\n", s.Nonce)
	res += fmt.Sprintf("Coinbase: %s\n", s.Coinbase.String())
	if s.PolygonZkEVMBatchData != nil {
		res += fmt.Sprintf("PolygonZkEVMBatchData: %v\n", *s.PolygonZkEVMBatchData)
	} else {
		res += "PolygonZkEVMBatchData: nil\n"
	}
	if s.PolygonRollupBaseEtrogBatchData != nil {
		res += fmt.Sprintf("___PolygonRollupBaseEtrogBatchData:ForcedTimestamp: %d\n", s.PolygonRollupBaseEtrogBatchData.ForcedTimestamp)
		res += fmt.Sprintf("___PolygonRollupBaseEtrogBatchData:ForcedGlobalExitRoot: %s\n", hex.EncodeToString(s.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot[:]))
		res += fmt.Sprintf("___PolygonRollupBaseEtrogBatchData:ForcedBlockHashL1: %s\n", hex.EncodeToString(s.PolygonRollupBaseEtrogBatchData.ForcedBlockHashL1[:]))
		res += fmt.Sprintf("___PolygonRollupBaseEtrogBatchData:Transactions: %s\n", hex.EncodeToString(s.PolygonRollupBaseEtrogBatchData.Transactions))

	} else {
		res += "PolygonRollupBaseEtrogBatchData: nil\n"
	}
	if s.SequencedBatchElderberryData != nil {
		res += fmt.Sprintf("___SequencedBatchElderberryData:MaxSequenceTimestamp %d\n", s.SequencedBatchElderberryData.MaxSequenceTimestamp)
		res += fmt.Sprintf("___SequencedBatchElderberryData:InitSequencedBatchNumber %d\n", s.SequencedBatchElderberryData.InitSequencedBatchNumber)

	} else {
		res += "SequencedBatchElderberryData: nil\n"
	}
	if s.Metadata != nil {
		res += fmt.Sprintf("Metadata: %s\n", s.Metadata.String())
	} else {
		res += "Metadata: nil\n"

	}
	return res
}

func (s *SequencedBatch) BatchL2Data() []byte {
	if s.PolygonZkEVMBatchData != nil {
		return s.PolygonZkEVMBatchData.Transactions
	}
	if s.PolygonRollupBaseEtrogBatchData != nil {
		return s.PolygonRollupBaseEtrogBatchData.Transactions
	}
	return nil
}

// UpdateEtrogSequence represents the first etrog sequence
type UpdateEtrogSequence struct {
	BatchNumber   uint64
	SequencerAddr common.Address
	TxHash        common.Hash
	Nonce         uint64
	// Struct used in Etrog
	*polygonzkevmetrog.PolygonRollupBaseEtrogBatchData
}

// ForcedBatch represents a ForcedBatch
type ForcedBatch struct {
	BlockNumber       uint64
	ForcedBatchNumber uint64
	Sequencer         common.Address
	GlobalExitRoot    common.Hash
	RawTxsData        []byte
	ForcedAt          time.Time
}

// VerifiedBatch represents a VerifiedBatch
type VerifiedBatch struct {
	BlockNumber uint64
	BatchNumber uint64
	Aggregator  common.Address
	StateRoot   common.Hash
	TxHash      common.Hash
}

// SequencedForceBatch is a sturct to track the ForceSequencedBatches event.
type SequencedForceBatch struct {
	BatchNumber uint64
	Coinbase    common.Address
	TxHash      common.Hash
	Timestamp   time.Time
	Nonce       uint64
	polygonzkevmetrog.PolygonRollupBaseEtrogBatchData
}

// ForkID is a sturct to track the ForkID event.
type ForkID struct {
	BatchNumber uint64
	ForkID      uint64
	Version     string
}
