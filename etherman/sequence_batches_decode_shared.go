package etherman

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type DecodedSequenceBatchesCallData struct {
	InputByteData []byte
	Data          []interface{}
}

func decodeSequenceCallData(smcAbi abi.ABI, txData []byte) (*DecodedSequenceBatchesCallData, error) {
	// Extract coded txs.
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
	bytedata, err := json.Marshal(data[0])
	if err != nil {
		return nil, err
	}
	return &DecodedSequenceBatchesCallData{InputByteData: bytedata, Data: data}, nil
}
