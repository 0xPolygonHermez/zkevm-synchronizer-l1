package etherman

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/oldpolygonzkevm"
	"github.com/ethereum/go-ethereum/common"
)

type L1InfoTreeV2Data struct {
	CurrentL1InfoRoot common.Hash
	LeafCount         uint32
	BlockHash         common.Hash
	MinTimestamp      uint64
}

func (l *L1InfoTreeV2Data) String() string {
	return fmt.Sprintf("CurrentL1InfoRoot: %s LeafCount: %d BlockHash: %s MinTimestamp: %d", l.CurrentL1InfoRoot.String(), l.LeafCount, l.BlockHash.String(), l.MinTimestamp)
}

type RollbackBatchesData struct {
	TargetBatch            uint64
	AccInputHashToRollback common.Hash
}

func (r *RollbackBatchesData) String() string {
	return fmt.Sprintf("TargetBatch: %d AccInputHashToRollback: %s", r.TargetBatch, r.AccInputHashToRollback.String())
}

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
	L1InfoTreeV2                []L1InfoTreeV2Data
	RollbackBatches             []RollbackBatchesData
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

func (g *GlobalExitRoot) String() string {
	return fmt.Sprintf("BlockNumber: %d MainnetExitRoot: %s RollupExitRoot: %s GlobalExitRoot: %s Timestamp: %s PreviousBlockHash: %s",
		g.BlockNumber, g.MainnetExitRoot.String(), g.RollupExitRoot.String(), g.GlobalExitRoot.String(), g.Timestamp.String(), g.PreviousBlockHash.String())
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

type BananaSequenceData struct {
	CounterL1InfoRoot         uint32
	MaxSequenceTimestamp      uint64
	ExpectedFinalAccInputHash common.Hash
	DataAvailabilityMsg       []byte
}

func (b *BananaSequenceData) String() string {
	res := fmt.Sprintf("CounterL1InfoRoot: %d MaxSequenceTimestamp: %d ExpectedFinalAccInputHash: %s", b.CounterL1InfoRoot, b.MaxSequenceTimestamp, b.ExpectedFinalAccInputHash.String())
	daMsg := fmt.Sprintf("DataAvailabilityMsg(%d):", len(b.DataAvailabilityMsg))
	if len(b.DataAvailabilityMsg) > 0 {
		daMsg += " " + hex.EncodeToString(b.DataAvailabilityMsg)
	}
	return res + " " + daMsg
}

func (b *BananaSequenceData) ToJson() string {
	jsonData, err := json.Marshal(b)
	if err != nil {
		return "error"
	}
	return string(jsonData)
}

// EtrogSequenceData also apply to Elderberry
type EtrogSequenceData struct {
	Transactions         []byte
	ForcedGlobalExitRoot common.Hash
	ForcedTimestamp      uint64
	ForcedBlockHashL1    common.Hash
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
	//*polygonzkevm.PolygonRollupBaseEtrogBatchData
	*EtrogSequenceData
	// Struct used in Elderberry
	*SequencedBatchElderberryData
	BananaData *BananaSequenceData
	Metadata   *SequencedBatchMetadata
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
	if s.EtrogSequenceData != nil {
		res += fmt.Sprintf("___EtrogSequenceData:ForcedTimestamp: %d\n", s.EtrogSequenceData.ForcedTimestamp)
		res += fmt.Sprintf("___EtrogSequenceData:ForcedGlobalExitRoot: %s\n", hex.EncodeToString(s.EtrogSequenceData.ForcedGlobalExitRoot[:]))
		res += fmt.Sprintf("___EtrogSequenceData:ForcedBlockHashL1: %s\n", hex.EncodeToString(s.EtrogSequenceData.ForcedBlockHashL1[:]))
		res += fmt.Sprintf("___EtrogSequenceData:Transactions: %s\n", hex.EncodeToString(s.EtrogSequenceData.Transactions))

	} else {
		res += "EtrogSequenceData: nil\n"
	}
	if s.SequencedBatchElderberryData != nil {
		res += fmt.Sprintf("___SequencedBatchElderberryData:MaxSequenceTimestamp %d\n", s.SequencedBatchElderberryData.MaxSequenceTimestamp)
		res += fmt.Sprintf("___SequencedBatchElderberryData:InitSequencedBatchNumber %d\n", s.SequencedBatchElderberryData.InitSequencedBatchNumber)

	} else {
		res += "SequencedBatchElderberryData: nil\n"
	}
	if s.BananaData != nil {
		res += fmt.Sprintf("BananaData: %s\n", s.BananaData.String())
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
	if s.EtrogSequenceData != nil {
		return s.EtrogSequenceData.Transactions
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
	*EtrogSequenceData
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
	EtrogSequenceData
}

// ForkID is a sturct to track the ForkID event.
type ForkID struct {
	BatchNumber uint64
	ForkID      uint64
	Version     string
}
