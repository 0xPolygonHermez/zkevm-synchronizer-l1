package etherman

import (
	"encoding/json"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogvalidiumpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// methodIDSequenceBatchesValidiumEtrog: MethodID for sequenceBatchesValidium in Etrog
	methodIDSequenceBatchesValidiumEtrog = []byte{0x2d, 0x72, 0xc2, 0x48} // 0x2d72c248 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],address,bytes)
)

type SequenceBatchesDecodeEtrogValidium struct {
	SequenceBatchesBase
	da dataavailability.BatchDataProvider
}

func NewDecodeSequenceBatchesEtrogValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecodeEtrogValidium, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesValidiumEtrog, "sequenceBatchesEtrogValidium", etrogvalidiumpolygonzkevm.EtrogvalidiumpolygonzkevmABI)
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecodeEtrogValidium{*base, da}, nil
}

func (s *SequenceBatchesDecodeEtrogValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	if s.da == nil {
		return nil, fmt.Errorf("data availability backend not set")
	}
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData

	var (
		// etrogvalidiumpolygonzkevm.PolygonValidiumEtrogValidiumBatchData and polygonzkevm.PolygonValidiumEtrogValidiumBatchData
		// are the same struct
		sequencesValidium   []polygonzkevm.PolygonValidiumEtrogValidiumBatchData
		dataAvailabilityMsg []byte
	)
	err = json.Unmarshal(bytedata, &sequencesValidium)
	if err != nil {
		return nil, err
	}

	coinbase := data[1].(common.Address)
	dataAvailabilityMsg = data[2].([]byte)

	batchInfos := createBatchInfo(sequencesValidium, lastBatchNumber)

	batchData, err := retrieveBatchData(s.da, batchInfos, dataAvailabilityMsg)
	if err != nil {
		return nil, err
	}

	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: s.NameMethodID(txData[:4]),
		RollupFlavor:     RollupFlavorValidium,
		ForkName:         "etrog",
	}

	sequencedBatches := createSequencedBatchList(sequencesValidium, batchInfos, batchData, l1InfoRoot, sequencer, txHash, nonce, coinbase,
		uint64(0), uint64(0), SequencedBatchMetadata)

	return sequencedBatches, nil

}

func retrieveBatchData(da dataavailability.BatchDataProvider, batchInfos []batchInfo, daMessage []byte) ([]dataavailability.BatchL2Data, error) {
	validiumData, err := getBatchL2Data(da, batchInfos, daMessage)
	if err != nil {
		return nil, err
	}

	data := make([]dataavailability.BatchL2Data, len(batchInfos))
	for i, info := range batchInfos {
		bn := info.num
		data[i] = validiumData[bn]

	}
	return data, nil
}
