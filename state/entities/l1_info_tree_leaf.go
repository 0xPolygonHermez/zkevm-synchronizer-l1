package entities

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type L1InfoTreeLeaf struct {
	L1InfoTreeRoot    common.Hash
	L1InfoTreeIndex   uint32
	PreviousBlockHash common.Hash
	BlockNumber       uint64
	Timestamp         time.Time
	MainnetExitRoot   common.Hash
	RollupExitRoot    common.Hash
	GlobalExitRoot    common.Hash
}
