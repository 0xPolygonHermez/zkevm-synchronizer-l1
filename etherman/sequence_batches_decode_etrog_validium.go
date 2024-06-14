package etherman

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// methodIDSequenceBatchesValidiumEtrog: MethodID for sequenceBatchesValidium in Etrog
	methodIDSequenceBatchesValidiumEtrog = []byte{0x2d, 0x72, 0xc2, 0x48} // 0x2d72c248 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],address,bytes)
)

type SequenceBatchesDecodeEtrogValidium struct {
}

func NewDecodeSequenceBatchesEtrogValidium() *SequenceBatchesDecodeEtrogValidium {
	return &SequenceBatchesDecodeEtrogValidium{}
}

// MatchMethodId returns true if the methodId is the one for the sequenceBatchesEtrog method
func (s *SequenceBatchesDecodeEtrogValidium) MatchMethodId(methodId []byte) bool {
	return bytes.Equal(methodId, methodIDSequenceBatchesValidiumEtrog)
}

func (s *SequenceBatchesDecodeEtrogValidium) NameMethodID(methodId []byte) string {
	if s.MatchMethodId(methodId) {
		return "sequenceBatchesEtrogValidium"
	}
	return ""
}

func (s *SequenceBatchesDecodeEtrogValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	/*var (
		sequencesValidium   []polygonzkevmvalidium.PolygonValidiumEtrogValidiumBatchData
		dataAvailabilityMsg []byte
	)
	err := json.Unmarshal(txData, &sequencesValidium)
	if err != nil {
		return nil, err
	}

	switch forkID {
	case state.FORKID_ETROG:
		coinbase = data[1].(common.Address)
		dataAvailabilityMsg = data[2].([]byte)

	case state.FORKID_ELDERBERRY:
		maxSequenceTimestamp = data[1].(uint64)
		initSequencedBatchNumber = data[2].(uint64)
		coinbase = data[3].(common.Address)
		dataAvailabilityMsg = data[4].([]byte)
	}

	// Pair the batch number, hash, and if it is forced. This will allow
	// retrieval from different sources, and keep them in original order.
	var batchInfos []batchInfo
	for i, d := range sequencesValidium {
		bn := lastBatchNumber - uint64(len(sequencesValidium)-(i+1))
		forced := d.ForcedTimestamp > 0
		h := d.TransactionsHash
		batchInfos = append(batchInfos, batchInfo{num: bn, hash: h, isForced: forced})
	}

	batchData, err := retrieveBatchData(da, st, batchInfos, dataAvailabilityMsg)
	if err != nil {
		return nil, err
	}

	sequencedBatches := make([]SequencedBatch, len(sequencesValidium))
	for i, info := range batchInfos {
		bn := info.num
		s := polygonzkevm.PolygonRollupBaseEtrogBatchData{
			Transactions:         batchData[i],
			ForcedGlobalExitRoot: sequencesValidium[i].ForcedGlobalExitRoot,
			ForcedTimestamp:      sequencesValidium[i].ForcedTimestamp,
			ForcedBlockHashL1:    sequencesValidium[i].ForcedBlockHashL1,
		}
		batch := SequencedBatch{
			BatchNumber:                     bn,
			L1InfoRoot:                      &l1InfoRoot,
			SequencerAddr:                   sequencer,
			TxHash:                          txHash,
			Nonce:                           nonce,
			Coinbase:                        coinbase,
			PolygonRollupBaseEtrogBatchData: &s,
		}
		if forkID >= state.FORKID_ELDERBERRY {
			elderberry := &SequencedBatchElderberryData{
				MaxSequenceTimestamp:     maxSequenceTimestamp,
				InitSequencedBatchNumber: initSequencedBatchNumber,
			}
			batch.SequencedBatchElderberryData = elderberry
		}
		sequencedBatches[i] = batch
	}

	return sequencedBatches, nil
	*/
	return nil, fmt.Errorf("Not implemented")
}
