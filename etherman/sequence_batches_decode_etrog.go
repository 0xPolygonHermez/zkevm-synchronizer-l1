package etherman

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type SequenceBatchesDecodeEtrog struct {
}

func NewDecodeSequenceBatchesEtrog() *SequenceBatchesDecodeEtrog {
	return &SequenceBatchesDecodeEtrog{}
}

// MatchMethodId returns true if the methodId is the one for the sequenceBatchesEtrog method
func (s *SequenceBatchesDecodeEtrog) MatchMethodId(methodId []byte) bool {
	return bytes.Equal(methodId, methodIDSequenceBatchesEtrog)
}

func (s *SequenceBatchesDecodeEtrog) NameMethodID(methodId []byte) string {
	if s.MatchMethodId(methodId) {
		return "sequenceBatchesEtrog"
	}
	return ""
}

func (s *SequenceBatchesDecodeEtrog) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi, err := abi.JSON(strings.NewReader(etrogpolygonzkevm.EtrogpolygonzkevmABI))
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

	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: "sequenceBatchesEtrog",
		RollupFlavor:     RollupFlavorZkEVM,
		ForkName:         "etrog",
	}

	coinbase := (data[1]).(common.Address)
	sequencedBatches := make([]SequencedBatch, len(sequences))
	for i, seq := range sequences {
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
			Metadata:                        SequencedBatchMetadata,
		}
	}

	return sequencedBatches, nil
}
