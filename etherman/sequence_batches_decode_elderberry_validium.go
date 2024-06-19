package etherman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// methodIDSequenceBatchesValidiumElderberry: MethodID for sequenceBatchesValidium in Elderberry
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/main/contracts/v2/consensus/validium/PolygonValidiumEtrog.sol
	// function sequenceBatchesValidium(
	//     ValidiumBatchData[] calldata batches,
	//     uint64 maxSequenceTimestamp,
	//     uint64 initSequencedBatch,
	//     address l2Coinbase,
	//     bytes calldata dataAvailabilityMessage
	// ) external onlyTrustedSequencer {
	//
	//struct ValidiumBatchData {
	//     bytes32 transactionsHash;
	//     bytes32 forcedGlobalExitRoot;
	//     uint64 forcedTimestamp;
	//     bytes32 forcedBlockHashL1;
	// }
	methodIDSequenceBatchesValidiumElderberry = []byte{0xdb, 0x5b, 0x0e, 0xd7} // 0xdb5b0ed7 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],uint64,uint64,address,bytes)
)

type SequenceBatchesDecodeElderberryValidium struct {
	da     dataavailability.BatchDataProvider
	SmcABI abi.ABI
}

func NewDecodeSequenceBatchesElderberryValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecodeElderberryValidium, error) {
	smcAbi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecodeElderberryValidium{da, smcAbi}, nil
}

// MatchMethodId returns true if the methodId is the one for the sequenceBatchesEtrog method
func (s *SequenceBatchesDecodeElderberryValidium) MatchMethodId(methodId []byte) bool {
	return bytes.Equal(methodId, methodIDSequenceBatchesValidiumElderberry)
}

func (s *SequenceBatchesDecodeElderberryValidium) NameMethodID(methodId []byte) string {
	if s.MatchMethodId(methodId) {
		return "sequenceBatchesElderberryValidium"
	}
	return ""
}

func (s *SequenceBatchesDecodeElderberryValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	if s.da == nil {
		return nil, fmt.Errorf("data availability backend not set")
	}
	decoded, err := decodeSequenceCallData(s.SmcABI, txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData
	var sequencesValidium []polygonzkevm.PolygonValidiumEtrogValidiumBatchData
	err = json.Unmarshal(bytedata, &sequencesValidium)
	if err != nil {
		return nil, err
	}

	maxSequenceTimestamp := data[1].(uint64)
	initSequencedBatchNumber := data[2].(uint64)
	coinbase := data[3].(common.Address)
	dataAvailabilityMsg := data[4].([]byte)

	batchInfos := createBatchInfo(sequencesValidium, lastBatchNumber)

	batchData, err := retrieveBatchData(s.da, batchInfos, dataAvailabilityMsg)
	if err != nil {
		return nil, err
	}
	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: "sequenceBatchesElderberryValidium",
		ForkName:         "elderberry",
		RollupFlavor:     RollupFlavorValidium,
	}

	sequencedBatches := createSequencedBatchList(sequencesValidium, batchInfos, batchData, l1InfoRoot, sequencer, txHash, nonce, coinbase, maxSequenceTimestamp, initSequencedBatchNumber, SequencedBatchMetadata)
	return sequencedBatches, nil
}
