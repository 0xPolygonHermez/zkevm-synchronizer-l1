package etherman

import (
	"encoding/json"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type batchInfo struct {
	num      uint64
	hash     common.Hash
	isForced bool
}

type DecodedSequenceBatchesCallData struct {
	InputByteData []byte
	Data          []interface{}
}

func decodeSequenceCallData(smcAbi abi.ABI, txData []byte) (*DecodedSequenceBatchesCallData, error) {
	// Extract coded txs.
	// Recover Method from signature and ABI
	method, err := smcAbi.MethodById(txData[:4])
	if err != nil {
		return nil, err
	}

	// Unpack method inputs
	data, err := method.Inputs.Unpack(txData[4:])
	if err != nil {
		return nil, err
	}
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	return &DecodedSequenceBatchesCallData{InputByteData: bytedata, Data: data}, nil
}

func createBatchInfo(sequencesValidium []polygonzkevm.PolygonValidiumEtrogValidiumBatchData, lastBatchNumber uint64) []batchInfo {
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

func createSequencedBatchList(sequencesValidium []polygonzkevm.PolygonValidiumEtrogValidiumBatchData, batchInfos []batchInfo, batchData []dataavailability.BatchL2Data,
	l1InfoRoot common.Hash, sequencer common.Address, txHash common.Hash, nonce uint64, coinbase common.Address,
	maxSequenceTimestamp uint64, initSequencedBatchNumber uint64,
	metaData *SequencedBatchMetadata) []SequencedBatch {
	sequencedBatches := make([]SequencedBatch, len(sequencesValidium))
	for i, info := range batchInfos {
		bn := info.num
		s := polygonzkevm.PolygonRollupBaseEtrogBatchData{
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
			BatchNumber:                     bn,
			L1InfoRoot:                      &l1InfoRoot,
			SequencerAddr:                   sequencer,
			TxHash:                          txHash,
			Nonce:                           nonce,
			Coinbase:                        coinbase,
			PolygonRollupBaseEtrogBatchData: &s,
			Metadata:                        metaData,
		}

		elderberry := &SequencedBatchElderberryData{
			MaxSequenceTimestamp:     maxSequenceTimestamp,
			InitSequencedBatchNumber: initSequencedBatchNumber,
		}
		batch.SequencedBatchElderberryData = elderberry
		sequencedBatches[i] = batch
	}
	return sequencedBatches
}

func getBatchL2Data(da dataavailability.BatchDataProvider, batchInfos []batchInfo, daMessage []byte) (map[uint64]dataavailability.BatchL2Data, error) {
	var batchNums []uint64
	var batchHashes []common.Hash
	for _, info := range batchInfos {
		if !info.isForced {
			batchNums = append(batchNums, info.num)
			batchHashes = append(batchHashes, info.hash)
		}
	}
	if len(batchNums) == 0 {
		return nil, nil
	}

	batchL2Data, err := da.GetBatchL2Data(batchNums, batchHashes, daMessage)
	if err != nil {
		return nil, err
	}

	if len(batchL2Data) != len(batchNums) {
		return nil,
			fmt.Errorf("failed to retrieve all batch data. Expected %d, got %d", len(batchNums), len(batchL2Data))
	}

	data := make(map[uint64]dataavailability.BatchL2Data)
	for i, bn := range batchNums {
		data[bn] = batchL2Data[i]
	}

	return data, nil
}
