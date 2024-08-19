package etherman

import (
	"encoding/json"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonzkevmetrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogpolygonzkevm"
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
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData
	var sequences []polygonzkevmetrog.PolygonRollupBaseEtrogBatchData
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
