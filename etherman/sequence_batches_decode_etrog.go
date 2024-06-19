package etherman

import (
	"encoding/json"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	"github.com/ethereum/go-ethereum/common"
)

type SequenceBatchesDecodeEtrog struct {
	SequenceBatchesBase
}

func NewDecodeSequenceBatchesEtrog() (*SequenceBatchesDecodeEtrog, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesEtrog, "sequenceBatchesEtrog", etrogpolygonzkevm.EtrogpolygonzkevmABI)
	if err != nil {
		return nil, err
	}
	return &SequenceBatchesDecodeEtrog{*base}, nil
}

func (s *SequenceBatchesDecodeEtrog) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	// Extract coded txs.
	// Load contract ABI
	smcAbi := s.SmcABI()

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
		CallFunctionName: s.NameMethodID(txData[:4]),
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
