package etherman

import (
	"encoding/json"
	"fmt"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonvalidiumetrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
)

var (
	/*
		// contract: v8.0.0-rc.1-fork.12
			https://github.com/0xPolygonHermez/zkevm-contracts/blob/a5eacc6e51d7456c12efcabdfc1c37457f2219b2/contracts/v2/consensus/validium/PolygonValidiumEtrog.sol#L29
			  struct ValidiumBatchData {
				bytes32 transactionsHash;
				bytes32 forcedGlobalExitRoot;
				uint64 forcedTimestamp;
				bytes32 forcedBlockHashL1;
			}
			//https://github.com/0xPolygonHermez/zkevm-contracts/blob/a9b4f742f66bd4f3bcd98a3a188422480ffe0d4e/contracts/v2/consensus/validium/PolygonValidiumEtrog.sol#L91
			function sequenceBatchesValidium(
				ValidiumBatchData[] calldata batches,
				uint32 indexL1InfoRoot,
				uint64 maxSequenceTimestamp,
				bytes32 expectedFinalAccInputHash,
				address l2Coinbase,
				bytes calldata dataAvailabilityMessage
			)
		165e8a8d50cd47dabdf9bde8bf707c673d2379d465bd458693e23b75ab3a4424
		sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],uint32,uint64,bytes32,address,bytes)

	*/
	methodIDSequenceBatchesBananaValidium     = []byte{0x16, 0x5e, 0x8a, 0x8d} // 165e8a8d sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],uint32,uint64,bytes32,address,bytes)
	methodIDSequenceBatchesBananaValidiumName = "sequenceBatchesBananaValidium"
)

type SequenceBatchesDecoderBananaValidium struct {
	SequenceBatchesBase
	da dataavailability.BatchDataProvider
}

func NewSequenceBatchesDecoderBananaValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecoderBananaValidium, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesBananaValidium,
		methodIDSequenceBatchesBananaValidiumName, polygonvalidiumetrog.PolygonvalidiumetrogABI)
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecoderBananaValidium{*base, da}, nil
}

func (s *SequenceBatchesDecoderBananaValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {

	if s.da == nil {
		return nil, fmt.Errorf("data availability backend not set")
	}
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData
	var sequencesValidium []polygonvalidiumetrog.PolygonValidiumEtrogValidiumBatchData
	err = json.Unmarshal(bytedata, &sequencesValidium)
	if err != nil {
		return nil, err
	}

	counterL1InfoRoot := data[1].(uint32)
	maxSequenceTimestamp := data[2].(uint64)
	//expectedFinalAccInputHash := data[3].(common.Hash)
	expectedFinalAccInputHashraw := data[3].([common.HashLength]byte)
	expectedFinalAccInputHash := common.Hash(expectedFinalAccInputHashraw)
	coinbase := data[4].(common.Address)
	dataAvailabilityMsg := data[5].([]byte)

	bananaData := BananaSequenceData{
		CounterL1InfoRoot:         counterL1InfoRoot,
		MaxSequenceTimestamp:      maxSequenceTimestamp,
		ExpectedFinalAccInputHash: expectedFinalAccInputHash,
		DataAvailabilityMsg:       dataAvailabilityMsg,
	}

	batchInfos := createBatchInfoBanana(sequencesValidium, lastBatchNumber)

	batchData, err := retrieveBatchData(s.da, batchInfos, dataAvailabilityMsg)
	if err != nil {
		return nil, err
	}

	log.Debugf("Decoded Banana sequenceBatchesValidium: counterL1InfoRoot:%d maxSequenceTimestamp:%d expectedFinalAccInputHash:%s coinbase:%s dataAvailabilityMsg:%s",
		counterL1InfoRoot, maxSequenceTimestamp, expectedFinalAccInputHash, coinbase.Hex(), dataAvailabilityMsg)
	log.Debugf("%s batchNum: %d Data:%s", methodIDSequenceBatchesBananaValidiumName, lastBatchNumber+1, common.Bytes2Hex(txData))
	for i, d := range batchData {
		log.Debugf("%s    BatchData[%d]: %s", methodIDSequenceBatchesBananaValidiumName, i, common.Bytes2Hex(d.Data))
	}
	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: s.NameMethodID(txData[:4]),
		ForkName:         "banana",
		RollupFlavor:     RollupFlavorValidium,
	}
	sequencedBatches := createSequencedBatchListBanana(sequencesValidium, batchInfos, batchData, l1InfoRoot, sequencer, txHash, nonce, coinbase, maxSequenceTimestamp, bananaData, SequencedBatchMetadata)

	return sequencedBatches, nil

}

func createBatchInfoBanana(sequencesValidium []polygonvalidiumetrog.PolygonValidiumEtrogValidiumBatchData, lastBatchNumber uint64) []batchInfo {
	// Pair the batch number, hash, and if it is forced. This will allow
	// retrieval from different sources, and keep them in original order.

	var batchInfos []batchInfo
	for i, d := range sequencesValidium {
		bn := lastBatchNumber - uint64(len(sequencesValidium)-(i+1))
		forced := d.ForcedTimestamp > 0
		h := d.TransactionsHash
		batchInfos = append(batchInfos, batchInfo{num: bn, hash: h, isForced: forced})
	}
	return batchInfos
}

func createSequencedBatchListBanana(sequencesValidium []polygonvalidiumetrog.PolygonValidiumEtrogValidiumBatchData, batchInfos []batchInfo, batchData []dataavailability.BatchL2Data,
	l1InfoRoot common.Hash, sequencer common.Address, txHash common.Hash, nonce uint64, coinbase common.Address,
	maxSequenceTimestamp uint64,
	bananaSequenceData BananaSequenceData,
	metaData *SequencedBatchMetadata) []SequencedBatch {
	sequencedBatches := make([]SequencedBatch, len(sequencesValidium))
	for i, info := range batchInfos {
		bn := info.num
		s := EtrogSequenceData{
			Transactions:         batchData[i].Data,
			ForcedGlobalExitRoot: sequencesValidium[i].ForcedGlobalExitRoot,
			ForcedTimestamp:      sequencesValidium[i].ForcedTimestamp,
			ForcedBlockHashL1:    sequencesValidium[i].ForcedBlockHashL1,
		}
		if metaData != nil {
			switch batchData[i].Source {
			case dataavailability.External:
				metaData.SourceBatchData = SourceBatchDataValidiumDAExternal
			case dataavailability.Trusted:
				metaData.SourceBatchData = SourceBatchDataValidiumDATrusted

				metaData.SourceBatchData = string(batchData[i].Source)
			}
		}
		batch := SequencedBatch{
			BatchNumber:       bn,
			L1InfoRoot:        &l1InfoRoot,
			SequencerAddr:     sequencer,
			TxHash:            txHash,
			Nonce:             nonce,
			Coinbase:          coinbase,
			EtrogSequenceData: &s,
			BananaData:        &bananaSequenceData,
			Metadata:          metaData,
		}
		sequencedBatches[i] = batch
	}
	return sequencedBatches
}
