package etherman

import (
	"bytes"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/dataavailability"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// methodIDSequenceBatchesValidiumElderberry: MethodID for sequenceBatchesValidium in Elderberry
	// https://github.com/0xPolygonHermez/zkevm-contracts/blob/main/contracts/v2/previousVersions/PolygonValidiumEtrogPrevious.sol
	methodIDSequenceBatchesValidiumElderberry = []byte{0xdb, 0x5b, 0x0e, 0xd7} // 0xdb5b0ed7 sequenceBatchesValidium((bytes32,bytes32,uint64,bytes32)[],uint64,uint64,address,bytes)
)

type SequenceBatchesDecodeElderberryValidium struct {
	da dataavailability.BatchDataProvider
}

func NewDecodeSequenceBatchesElderberryValidium(da dataavailability.BatchDataProvider) (*SequenceBatchesDecodeElderberryValidium, error) {
	return &SequenceBatchesDecodeElderberryValidium{da}, nil
}

// MatchMethodId returns true if the methodId is the one for the sequenceBatchesEtrog method
func (s *SequenceBatchesDecodeElderberryValidium) MatchMethodId(methodId []byte) bool {
	return bytes.Equal(methodId, methodIDSequenceBatchesValidiumElderberry)
}

func (s *SequenceBatchesDecodeElderberryValidium) NameMethodID(methodId []byte) string {
	if s.MatchMethodId(methodId) {
		return "sequenceBatchesElderberryValidium"
	}
	return ""
}

func (s *SequenceBatchesDecodeElderberryValidium) DecodeSequenceBatches(txData []byte, lastBatchNumber uint64, sequencer common.Address, txHash common.Hash, nonce uint64, l1InfoRoot common.Hash) ([]SequencedBatch, error) {
	return nil, fmt.Errorf("not implemented elderberry validium")
}
