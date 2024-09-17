package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ArgUint64 helps to marshal uint64 values provided in the RPC requests
type ArgUint64 uint64

// MarshalText marshals into text
func (b ArgUint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2) //nolint:gomnd
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(b), hex.Base)
	return buf, nil
}

// UnmarshalText unmarshals from text
func (b *ArgUint64) UnmarshalText(input []byte) error {
	str := strings.TrimPrefix(string(input), "0x")
	num, err := strconv.ParseUint(str, hex.Base, hex.BitSize64)
	if err != nil {
		return err
	}
	*b = ArgUint64(num)
	return nil
}

// Hex returns a hexadecimal representation
func (b ArgUint64) Hex() string {
	bb, _ := b.MarshalText()
	return string(bb)
}

// ArgUint64Ptr returns the pointer of the provided ArgUint64
func ArgUint64Ptr(a ArgUint64) *ArgUint64 {
	return &a
}

// ArgBytes helps to marshal byte array values provided in the RPC requests
type ArgBytes []byte

// MarshalText marshals into text
func (b ArgBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

// UnmarshalText unmarshals from text
func (b *ArgBytes) UnmarshalText(input []byte) error {
	hh, err := decodeToHex(input)
	if err != nil {
		return nil
	}
	aux := make([]byte, len(hh))
	copy(aux, hh)
	*b = aux
	return nil
}

// Hex returns a hexadecimal representation
func (b ArgBytes) Hex() string {
	bb, _ := b.MarshalText()
	return string(bb)
}

// ArgBytesPtr helps to marshal byte array values provided in the RPC requests
func ArgBytesPtr(b []byte) *ArgBytes {
	bb := ArgBytes(b)

	return &bb
}

// ArgBig helps to marshal big number values provided in the RPC requests
type ArgBig big.Int

// UnmarshalText unmarshals an instance of ArgBig into an array of bytes
func (a *ArgBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}

	b := new(big.Int)
	b.SetBytes(buf)
	*a = ArgBig(*b)

	return nil
}

// MarshalText marshals an array of bytes into an instance of ArgBig
func (a ArgBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)

	return []byte("0x" + b.Text(hex.Base)), nil
}

// Hex returns a hexadecimal representation
func (b ArgBig) Hex() string {
	bb, _ := b.MarshalText()
	return string(bb)
}

func decodeToHex(b []byte) ([]byte, error) {
	str := string(b)
	str = strings.TrimPrefix(str, "0x")
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return hex.DecodeString(str)
}

func encodeToHex(b []byte) []byte {
	str := hex.EncodeToString(b)
	if len(str)%2 != 0 {
		str = "0" + str
	}
	return []byte("0x" + str)
}

// ArgHash represents a common.Hash that accepts strings
// shorter than 64 bytes, like 0x00
type ArgHash common.Hash

// UnmarshalText unmarshals from text
func (arg *ArgHash) UnmarshalText(input []byte) error {
	if !hex.IsValid(string(input)) {
		return fmt.Errorf("invalid hash, it needs to be a hexadecimal value")
	}

	str := strings.TrimPrefix(string(input), "0x")
	*arg = ArgHash(common.HexToHash(str))
	return nil
}

// Hash returns an instance of common.Hash
func (arg *ArgHash) Hash() common.Hash {
	result := common.Hash{}
	if arg != nil {
		result = common.Hash(*arg)
	}
	return result
}

// ArgAddress represents a common.Address that accepts strings
// shorter than 32 bytes, like 0x00
type ArgAddress common.Address

// UnmarshalText unmarshals from text
func (b *ArgAddress) UnmarshalText(input []byte) error {
	if !hex.IsValid(string(input)) {
		return fmt.Errorf("invalid address, it needs to be a hexadecimal value")
	}

	str := strings.TrimPrefix(string(input), "0x")
	*b = ArgAddress(common.HexToAddress(str))
	return nil
}

// Address returns an instance of common.Address
func (arg *ArgAddress) Address() common.Address {
	result := common.Address{}
	if arg != nil {
		result = common.Address(*arg)
	}
	return result
}

// TxArgs is the transaction argument for the rpc endpoints
type TxArgs struct {
	From     *common.Address
	To       *common.Address
	Gas      *ArgUint64
	GasPrice *ArgBytes
	Value    *ArgBytes
	Data     *ArgBytes
	Input    *ArgBytes
	Nonce    *ArgUint64
}

// Block structure
type Block struct {
	ParentHash      common.Hash         `json:"parentHash"`
	Sha3Uncles      common.Hash         `json:"sha3Uncles"`
	Miner           *common.Address     `json:"miner"`
	StateRoot       common.Hash         `json:"stateRoot"`
	TxRoot          common.Hash         `json:"transactionsRoot"`
	ReceiptsRoot    common.Hash         `json:"receiptsRoot"`
	LogsBloom       types.Bloom         `json:"logsBloom"`
	Difficulty      ArgUint64           `json:"difficulty"`
	TotalDifficulty *ArgUint64          `json:"totalDifficulty"`
	Size            ArgUint64           `json:"size"`
	Number          ArgUint64           `json:"number"`
	GasLimit        ArgUint64           `json:"gasLimit"`
	GasUsed         ArgUint64           `json:"gasUsed"`
	Timestamp       ArgUint64           `json:"timestamp"`
	ExtraData       ArgBytes            `json:"extraData"`
	MixHash         common.Hash         `json:"mixHash"`
	Nonce           *ArgBytes           `json:"nonce"`
	Hash            *common.Hash        `json:"hash"`
	Transactions    []TransactionOrHash `json:"transactions"`
	Uncles          []common.Hash       `json:"uncles"`
	GlobalExitRoot  *common.Hash        `json:"globalExitRoot,omitempty"`
	BlockInfoRoot   *common.Hash        `json:"blockInfoRoot,omitempty"`
}

// Batch structure
type Batch struct {
	Number              ArgUint64           `json:"number"`
	ForcedBatchNumber   *ArgUint64          `json:"forcedBatchNumber,omitempty"`
	Coinbase            common.Address      `json:"coinbase"`
	StateRoot           common.Hash         `json:"stateRoot"`
	GlobalExitRoot      common.Hash         `json:"globalExitRoot"`
	MainnetExitRoot     common.Hash         `json:"mainnetExitRoot"`
	RollupExitRoot      common.Hash         `json:"rollupExitRoot"`
	LocalExitRoot       common.Hash         `json:"localExitRoot"`
	AccInputHash        common.Hash         `json:"accInputHash"`
	Timestamp           ArgUint64           `json:"timestamp"`
	SendSequencesTxHash *common.Hash        `json:"sendSequencesTxHash"`
	VerifyBatchTxHash   *common.Hash        `json:"verifyBatchTxHash"`
	Closed              bool                `json:"closed"`
	Blocks              []BlockOrHash       `json:"blocks"`
	Transactions        []TransactionOrHash `json:"transactions"`
	BatchL2Data         ArgBytes            `json:"batchL2Data"`
}

// BatchFilter is a list of batch numbers to retrieve
type BatchFilter struct {
	Numbers []BatchNumber `json:"numbers"`
}

// BatchData is an abbreviated structure that only contains the number and L2 batch data
type BatchData struct {
	Number      ArgUint64 `json:"number"`
	BatchL2Data ArgBytes  `json:"batchL2Data,omitempty"`
	Empty       bool      `json:"empty"`
}

// BatchDataResult is a list of BatchData for a BatchFilter
type BatchDataResult struct {
	Data []*BatchData
}

// TransactionOrHash for union type of transaction and types.Hash
type TransactionOrHash struct {
	Hash *common.Hash
	Tx   *Transaction
}

// MarshalJSON marshals into json
func (th TransactionOrHash) MarshalJSON() ([]byte, error) {
	if th.Hash != nil {
		return json.Marshal(th.Hash)
	}
	return json.Marshal(th.Tx)
}

// UnmarshalJSON unmarshals from json
func (th *TransactionOrHash) UnmarshalJSON(input []byte) error {
	v := string(input)
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "\"0x") {
		var h common.Hash
		err := json.Unmarshal(input, &h)
		if err != nil {
			return err
		}
		*th = TransactionOrHash{Hash: &h}
		return nil
	}

	var t Transaction
	err := json.Unmarshal(input, &t)
	if err != nil {
		return err
	}
	*th = TransactionOrHash{Tx: &t}
	return nil
}

// BlockOrHash for union type of block and types.Hash
type BlockOrHash struct {
	Hash  *common.Hash
	Block *Block
}

// MarshalJSON marshals into json
func (bh BlockOrHash) MarshalJSON() ([]byte, error) {
	if bh.Hash != nil {
		return json.Marshal(bh.Hash)
	}
	return json.Marshal(bh.Block)
}

// UnmarshalJSON unmarshals from json
func (bh *BlockOrHash) UnmarshalJSON(input []byte) error {
	v := string(input)
	if strings.HasPrefix(v, "0x") || strings.HasPrefix(v, "\"0x") {
		var h common.Hash
		err := json.Unmarshal(input, &h)
		if err != nil {
			return err
		}
		*bh = BlockOrHash{Hash: &h}
		return nil
	}

	var b Block
	err := json.Unmarshal(input, &b)
	if err != nil {
		return err
	}
	*bh = BlockOrHash{Block: &b}
	return nil
}

// Transaction structure
type Transaction struct {
	Nonce       ArgUint64       `json:"nonce"`
	GasPrice    ArgBig          `json:"gasPrice"`
	Gas         ArgUint64       `json:"gas"`
	To          *common.Address `json:"to"`
	Value       ArgBig          `json:"value"`
	Input       ArgBytes        `json:"input"`
	V           ArgBig          `json:"v"`
	R           ArgBig          `json:"r"`
	S           ArgBig          `json:"s"`
	Hash        common.Hash     `json:"hash"`
	From        common.Address  `json:"from"`
	BlockHash   *common.Hash    `json:"blockHash"`
	BlockNumber *ArgUint64      `json:"blockNumber"`
	TxIndex     *ArgUint64      `json:"transactionIndex"`
	ChainID     ArgBig          `json:"chainId"`
	Type        ArgUint64       `json:"type"`
	Receipt     *Receipt        `json:"receipt,omitempty"`
	L2Hash      *common.Hash    `json:"l2Hash,omitempty"`
}

// CoreTx returns a geth core type Transaction
func (t Transaction) CoreTx() *types.Transaction {
	return types.NewTx(&types.LegacyTx{
		Nonce:    uint64(t.Nonce),
		GasPrice: (*big.Int)(&t.GasPrice),
		Gas:      uint64(t.Gas),
		To:       t.To,
		Value:    (*big.Int)(&t.Value),
		Data:     t.Input,
		V:        (*big.Int)(&t.V),
		R:        (*big.Int)(&t.R),
		S:        (*big.Int)(&t.S),
	})
}

// Receipt structure
type Receipt struct {
	Root              *common.Hash    `json:"root,omitempty"`
	CumulativeGasUsed ArgUint64       `json:"cumulativeGasUsed"`
	LogsBloom         types.Bloom     `json:"logsBloom"`
	Logs              []*types.Log    `json:"logs"`
	Status            ArgUint64       `json:"status"`
	TxHash            common.Hash     `json:"transactionHash"`
	TxIndex           ArgUint64       `json:"transactionIndex"`
	BlockHash         common.Hash     `json:"blockHash"`
	BlockNumber       ArgUint64       `json:"blockNumber"`
	GasUsed           ArgUint64       `json:"gasUsed"`
	FromAddr          common.Address  `json:"from"`
	ToAddr            *common.Address `json:"to"`
	ContractAddress   *common.Address `json:"contractAddress"`
	Type              ArgUint64       `json:"type"`
	EffectiveGasPrice *ArgBig         `json:"effectiveGasPrice,omitempty"`
	TxL2Hash          *common.Hash    `json:"transactionL2Hash,omitempty"`
}

// Log structure
type Log struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        ArgBytes       `json:"data"`
	BlockNumber ArgUint64      `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	TxIndex     ArgUint64      `json:"transactionIndex"`
	BlockHash   common.Hash    `json:"blockHash"`
	LogIndex    ArgUint64      `json:"logIndex"`
	Removed     bool           `json:"removed"`
}

// NewLog creates a new instance of Log
func NewLog(l types.Log) Log {
	return Log{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockNumber: ArgUint64(l.BlockNumber),
		TxHash:      l.TxHash,
		TxIndex:     ArgUint64(l.TxIndex),
		BlockHash:   l.BlockHash,
		LogIndex:    ArgUint64(l.Index),
		Removed:     l.Removed,
	}
}

// ExitRoots structure
type ExitRoots struct {
	BlockNumber     ArgUint64   `json:"blockNumber"`
	Timestamp       ArgUint64   `json:"timestamp"`
	MainnetExitRoot common.Hash `json:"mainnetExitRoot"`
	RollupExitRoot  common.Hash `json:"rollupExitRoot"`
}

// ZKCounters counters for the tx
type ZKCounters struct {
	GasUsed              ArgUint64 `json:"gasUsed"`
	UsedKeccakHashes     ArgUint64 `json:"usedKeccakHashes"`
	UsedPoseidonHashes   ArgUint64 `json:"usedPoseidonHashes"`
	UsedPoseidonPaddings ArgUint64 `json:"usedPoseidonPaddings"`
	UsedMemAligns        ArgUint64 `json:"usedMemAligns"`
	UsedArithmetics      ArgUint64 `json:"usedArithmetics"`
	UsedBinaries         ArgUint64 `json:"usedBinaries"`
	UsedSteps            ArgUint64 `json:"usedSteps"`
	UsedSHA256Hashes     ArgUint64 `json:"usedSHA256Hashes"`
}

// ZKCountersLimits used to return the zk counter limits to the user
type ZKCountersLimits struct {
	MaxGasUsed          ArgUint64 `json:"maxGasUsed"`
	MaxKeccakHashes     ArgUint64 `json:"maxKeccakHashes"`
	MaxPoseidonHashes   ArgUint64 `json:"maxPoseidonHashes"`
	MaxPoseidonPaddings ArgUint64 `json:"maxPoseidonPaddings"`
	MaxMemAligns        ArgUint64 `json:"maxMemAligns"`
	MaxArithmetics      ArgUint64 `json:"maxArithmetics"`
	MaxBinaries         ArgUint64 `json:"maxBinaries"`
	MaxSteps            ArgUint64 `json:"maxSteps"`
	MaxSHA256Hashes     ArgUint64 `json:"maxSHA256Hashes"`
}

// RevertInfo contains the reverted message and data when a tx
// is reverted during the zk counter estimation
type RevertInfo struct {
	Message string    `json:"message"`
	Data    *ArgBytes `json:"data,omitempty"`
}

// ZKCountersResponse returned when counters are estimated
type ZKCountersResponse struct {
	CountersUsed   ZKCounters       `json:"countersUsed"`
	CountersLimits ZKCountersLimits `json:"countersLimit"`
	Revert         *RevertInfo      `json:"revert,omitempty"`
	OOCError       *string          `json:"oocError,omitempty"`
}
