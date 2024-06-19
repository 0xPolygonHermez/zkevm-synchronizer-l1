package etherman

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type SequenceBatchesDecodeElderberry struct {
}

func NewDecodeSequenceBatchesElderberry() *SequenceBatchesDecodeElderberry {
	return &SequenceBatchesDecodeElderberry{}
}

// MatchMethodId returns true if the methodId is the one for the sequenceBatchesEtrog method
func (s *SequenceBatchesDecodeElderberry) MatchMethodId(methodId []byte) bool {
	return bytes.Equal(methodId, methodIDSequenceBatchesElderberry)
}

func (s *SequenceBatchesDecodeElderberry) NameMethodID(methodId []byte) string {
	if s.MatchMethodId(methodId) {
		return "sequenceBatchesElderberry"
	}
	return ""
}

func (s *SequenceBatchesDecodeElderberry) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi, err := abi.JSON(strings.NewReader(polygonzkevm.PolygonzkevmABI))
	if err != nil {
		return nil, err
	}

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
	var sequences []polygonzkevm.PolygonRollupBaseEtrogBatchData
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}
	maxSequenceTimestamp := data[1].(uint64)
	initSequencedBatchNumber := data[2].(uint64)
	coinbase := (data[3]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))

	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: "sequenceBatchesElderberry",
		RollupFlavor:     RollupFlavorZkEVM,
		ForkName:         "elderberry",
	}

	for i, seq := range sequences {
		elderberry := SequencedBatchElderberryData{
			MaxSequenceTimestamp:     maxSequenceTimestamp,
			InitSequencedBatchNumber: initSequencedBatchNumber,
		}
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := seq
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:                     bn,
			L1InfoRoot:                      &l1InfoRoot,
			SequencerAddr:                   sequencer,
			TxHash:                          txHash,
			Nonce:                           nonce,
			Coinbase:                        coinbase,
			PolygonRollupBaseEtrogBatchData: &s,
			SequencedBatchElderberryData:    &elderberry,
			Metadata:                        SequencedBatchMetadata,
		}
	}

	return sequencedBatches, nil
}
