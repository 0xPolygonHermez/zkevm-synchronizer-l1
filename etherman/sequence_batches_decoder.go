package etherman

import "github.com/ethereum/go-ethereum/common"

// SequenceBatchesDecoder is an interface that defines the methods that a sequence batches decoder should implement
type SequenceBatchesDecoder interface {
	// MatchMethodId returns true if the class can decode this method
	MatchMethodId(methodId []byte) bool
	// NameMethodID returns the name of the methodID
	// if doesnt match the decoder = ""
	NameMethodID(methodId []byte) string
	DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error)
}
