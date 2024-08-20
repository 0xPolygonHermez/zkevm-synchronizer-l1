package etherman

import (
	"fmt"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonvalidiumetrog"
	"github.com/ethereum/go-ethereum/common"
)

var (
	/*
			struct BatchData {
		        bytes transactions;
		        bytes32 forcedGlobalExitRoot;
		        uint64 forcedTimestamp;
		        bytes32 forcedBlockHashL1;
		    }

			function sequenceBatches(
		        BatchData[] calldata batches,
		        uint32 indexL1InfoRoot,
		        uint64 maxSequenceTimestamp,
		        bytes32 expectedFinalAccInputHash,
		        address l2Coinbase
		    )
		b910e0f97e3128707f02262ec5c90481c982744e80bd892a4449261f745e95d0
		sequenceBatches((bytes,bytes32,uint64,bytes32)[],uint32,uint64,bytes32,address)

	*/
	methodIDSequenceBatchesBanana     = []byte{0xdb, 0x5b, 0x0e, 0xd7} // 165e8a8d sequenceBatches((bytes,bytes32,uint64,bytes32)[],uint32,uint64,bytes32,address)
	methodIDSequenceBatchesBananaName = "sequenceBatchesBanana"
)

type DecodeSequenceBatchesDecodeBanana struct {
	SequenceBatchesBase
}

func NewDecodeSequenceBatchesDecodeBanana() (*DecodeSequenceBatchesDecodeBanana, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesBanana, methodIDSequenceBatchesBananaName, polygonvalidiumetrog.PolygonvalidiumetrogABI)
	if err != nil {
		return nil, err
	}
	return &DecodeSequenceBatchesDecodeBanana{*base}, nil
}

func (s *DecodeSequenceBatchesDecodeBanana) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	//decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	_, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("not implemented")
	/*
		data := decoded.Data
		bytedata := decoded.InputByteData
		var sequences []polygonvalidiumetrog.PolygonRollupBaseEtrogBatchData
		err = json.Unmarshal(bytedata, &sequences)
		if err != nil {
			return nil, err
		}
		maxSequenceTimestamp := data[1].(uint64)
		initSequencedBatchNumber := data[2].(uint64)
		coinbase := (data[3]).(common.Address)
		sequencedBatches := make([]SequencedBatch, len(sequences))

		SequencedBatchMetadata := &SequencedBatchMetadata{
			CallFunctionName: s.NameMethodID(txData[:4]),
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
	*/
}
