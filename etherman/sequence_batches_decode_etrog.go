package etherman

import (
	"encoding/json"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/etrogpolygonzkevm"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/smartcontracts/polygonzkevm"
	ethtypes "github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman/types"
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

func (s *SequenceBatchesDecodeEtrog) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]ethtypes.SequencedBatch, error) {
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData
	var sequences []polygonzkevm.PolygonRollupBaseEtrogBatchData
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}

	SequencedBatchMetadata := &ethtypes.SequencedBatchMetadata{
		CallFunctionName: s.NameMethodID(txData[:4]),
		RollupFlavor:     ethtypes.RollupFlavorZkEVM,
		ForkName:         "etrog",
	}

	coinbase := (data[1]).(common.Address)
	sequencedBatches := make([]ethtypes.SequencedBatch, len(sequences))
	for i, seq := range sequences {
		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := ethtypes.EtrogSequenceData{
			Transactions:         seq.Transactions,
			ForcedGlobalExitRoot: seq.ForcedGlobalExitRoot,
			ForcedTimestamp:      seq.ForcedTimestamp,
			ForcedBlockHashL1:    seq.ForcedBlockHashL1,
		}
		sequencedBatches[i] = ethtypes.SequencedBatch{
			BatchNumber:       bn,
			L1InfoRoot:        &l1InfoRoot,
			SequencerAddr:     sequencer,
			TxHash:            txHash,
			Nonce:             nonce,
			Coinbase:          coinbase,
			EtrogSequenceData: &s,
			Metadata:          SequencedBatchMetadata,
		}
	}
	return sequencedBatches, nil
}
