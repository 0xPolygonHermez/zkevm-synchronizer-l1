package etherman

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
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
	methodIDSequenceBatchesValidiumElderberry     = []byte{0xdb, 0x5b, 0x0e, 0xd7} // 0xdb5b0ed7 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],uint64,uint64,address,bytes)
	methodIDSequenceBatchesValidiumElderberryName = "sequenceBatchesElderberryValidium"
)

type SequenceBatchesDecodeElderberryValidium struct {
	SequenceBatchesBase
	da dataavailability.BatchDataProvider
}

func NewDecodeSequenceBatchesElderberryValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecodeElderberryValidium, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesValidiumElderberry, methodIDSequenceBatchesValidiumElderberryName, polygonzkevm.PolygonzkevmABI)
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecodeElderberryValidium{*base, da}, nil
}

func (s *SequenceBatchesDecodeElderberryValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	hexStr := hex.EncodeToString(txData)
	log.Debug("txData=", hexStr)
	log.Debug("lastBatchNumber=", lastBatchNumber)
	log.Debug("sequencer=", sequencer.String())
	log.Debug("txHash=", txHash.String())
	log.Debug("nonce=", nonce)
	log.Debug("l1InfoRoot=", l1InfoRoot.String())
	if s.da == nil {
		return nil, fmt.Errorf("data availability backend not set")
	}
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
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
		CallFunctionName: s.NameMethodID(txData[:4]),
		ForkName:         "elderberry",
		RollupFlavor:     RollupFlavorValidium,
	}

	sequencedBatches := createSequencedBatchList(sequencesValidium, batchInfos, batchData, l1InfoRoot, sequencer, txHash, nonce, coinbase, maxSequenceTimestamp, initSequencedBatchNumber, SequencedBatchMetadata)

	for i, _ := range sequencedBatches {
		log.Debug("sequencedBatches[", i, "]=", sequencedBatches[i].String())
	}

	return sequencedBatches, nil
}
