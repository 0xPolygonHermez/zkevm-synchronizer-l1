package etherman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogvalidiumpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// methodIDSequenceBatchesValidiumEtrog: MethodID for sequenceBatchesValidium in Etrog
	// function sequenceBatchesValidium(
	//     ValidiumBatchData[] calldata batches,
	//     uint64 maxSequenceTimestamp,
	//     uint64 initSequencedBatch,
	//     address l2Coinbase,
	//     bytes calldata dataAvailabilityMessage
	// ) external onlyTrustedSequencer {
	methodIDSequenceBatchesValidiumEtrog = []byte{0x2d, 0x72, 0xc2, 0x48} // 0x2d72c248 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],address,bytes)
)

type SequenceBatchesDecodeEtrogValidium struct {
	da     dataavailability.BatchDataProvider
	SmcABI abi.ABI
}

type batchInfo struct {
	num      uint64
	hash     common.Hash
	isForced bool
}

func NewDecodeSequenceBatchesEtrogValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecodeEtrogValidium, error) {
	smcAbi, err := abi.JSON(strings.NewReader(etrogvalidiumpolygonzkevm.EtrogvalidiumpolygonzkevmABI))
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecodeEtrogValidium{da, smcAbi}, nil
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
	if s.da == nil {
		return nil, fmt.Errorf("data availability backend not set")
	}
	decoded, err := decodeSequenceCallData(s.SmcABI, txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData

	var (
		maxSequenceTimestamp     uint64
		initSequencedBatchNumber uint64
		coinbase                 common.Address
	)

	var (
		sequencesValidium   []etrogvalidiumpolygonzkevm.PolygonValidiumEtrogValidiumBatchData
		dataAvailabilityMsg []byte
	)
	err = json.Unmarshal(bytedata, &sequencesValidium)
	if err != nil {
		return nil, err
	}

	coinbase = data[1].(common.Address)
	dataAvailabilityMsg = data[2].([]byte)

	// Pair the batch number, hash, and if it is forced. This will allow
	// retrieval from different sources, and keep them in original order.
	var batchInfos []batchInfo
	for i, d := range sequencesValidium {
		bn := lastBatchNumber - uint64(len(sequencesValidium)-(i+1))
		forced := d.ForcedTimestamp > 0
		h := d.TransactionsHash
		batchInfos = append(batchInfos, batchInfo{num: bn, hash: h, isForced: forced})
	}

	batchData, err := retrieveBatchData(s.da, batchInfos, dataAvailabilityMsg)
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

		elderberry := &SequencedBatchElderberryData{
			MaxSequenceTimestamp:     maxSequenceTimestamp,
			InitSequencedBatchNumber: initSequencedBatchNumber,
		}
		batch.SequencedBatchElderberryData = elderberry
		sequencedBatches[i] = batch
	}

	return sequencedBatches, nil

}

func retrieveBatchData(da dataavailability.BatchDataProvider, batchInfos []batchInfo, daMessage []byte) ([][]byte, error) {
	validiumData, err := getBatchL2Data(da, batchInfos, daMessage)
	if err != nil {
		return nil, err
	}

	data := make([][]byte, len(batchInfos))
	for i, info := range batchInfos {
		bn := info.num
		data[i] = validiumData[bn]

	}
	return data, nil
}

func getBatchL2Data(da dataavailability.BatchDataProvider, batchInfos []batchInfo, daMessage []byte) (map[uint64][]byte, error) {
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

	data := make(map[uint64][]byte)
	for i, bn := range batchNums {
		data[bn] = batchL2Data[i]
	}

	return data, nil
}
