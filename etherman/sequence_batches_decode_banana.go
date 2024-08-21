package etherman

import (
	"encoding/json"

	"github.com/0xPolygon/cdk-contracts-tooling/contracts/banana/polygonvalidiumetrog"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	methodIDSequenceBatchesBanana     = crypto.Keccak256Hash([]byte("sequenceBatches((bytes,bytes32,uint64,bytes32)[],uint32,uint64,bytes32,address)")).Bytes()[:4] // b910e0f9
	methodIDSequenceBatchesBananaName = "sequenceBatchesBanana"
)

type DecodeSequenceBatchesBanana struct {
	SequenceBatchesBase
}

func NewDecodeSequenceBatchesBanana() (*DecodeSequenceBatchesBanana, error) {
	base, err := NewSequenceBatchesBase(methodIDSequenceBatchesBanana, methodIDSequenceBatchesBananaName, polygonvalidiumetrog.PolygonvalidiumetrogABI)
	if err != nil {
		return nil, err
	}
	return &DecodeSequenceBatchesBanana{*base}, nil
}

func (s *DecodeSequenceBatchesBanana) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	decoded, err := decodeSequenceCallData(s.SmcABI(), txData)
	if err != nil {
		return nil, err
	}
	data := decoded.Data
	bytedata := decoded.InputByteData
	var sequences []polygonvalidiumetrog.PolygonRollupBaseEtrogBatchData
	err = json.Unmarshal(bytedata, &sequences)
	if err != nil {
		return nil, err
	}

	counterL1InfoRoot := data[1].(uint32)
	maxSequenceTimestamp := data[2].(uint64)
	//expectedFinalAccInputHash := data[3].(common.Hash)
	expectedFinalAccInputHashraw := data[3].([common.HashLength]byte)
	expectedFinalAccInputHash := common.Hash(expectedFinalAccInputHashraw)
	coinbase := data[4].(common.Address)

	bananaData := BananaSequenceData{
		CounterL1InfoRoot:         counterL1InfoRoot,
		MaxSequenceTimestamp:      maxSequenceTimestamp,
		ExpectedFinalAccInputHash: expectedFinalAccInputHash,
		DataAvailabilityMsg:       []byte{},
	}
	SequencedBatchMetadata := &SequencedBatchMetadata{
		CallFunctionName: s.NameMethodID(txData[:4]),
		ForkName:         "banana",
		RollupFlavor:     RollupFlavorZkEVM,
	}

	log.Debugf("Decoded %s: event(lastBatchNumber:%d, sequencer:%s, txHash:%s, nonce:%d, l1InfoRoot:%s)  bananaData:%s",
		methodIDSequenceBatchesBananaName,
		lastBatchNumber, sequencer.String(), txHash.String(), nonce, l1InfoRoot.String(),
		bananaData.String())
	log.Debugf("%s batchNum: %d Data:%s", methodIDSequenceBatchesBananaName, lastBatchNumber+1, common.Bytes2Hex(txData))
	for i, d := range sequences {
		log.Debugf("%s    BatchData[%d]: %s", methodIDSequenceBatchesBananaName, i, common.Bytes2Hex(d.Transactions))
	}
	sequencedBatches := make([]SequencedBatch, len(sequences))

	for i, seq := range sequences {

		bn := lastBatchNumber - uint64(len(sequences)-(i+1))
		s := EtrogSequenceData{
			Transactions:         seq.Transactions,
			ForcedGlobalExitRoot: seq.ForcedGlobalExitRoot,
			ForcedTimestamp:      seq.ForcedTimestamp,
			ForcedBlockHashL1:    seq.ForcedBlockHashL1,
		}
		sequencedBatches[i] = SequencedBatch{
			BatchNumber:       bn,
			L1InfoRoot:        &l1InfoRoot,
			SequencerAddr:     sequencer,
			TxHash:            txHash,
			Nonce:             nonce,
			Coinbase:          coinbase,
			EtrogSequenceData: &s,
			BananaData:        &bananaData,
			Metadata:          SequencedBatchMetadata,
		}
	}
	return sequencedBatches, nil

}
