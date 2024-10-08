package entities

import (
	"fmt"
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

func (l *L1InfoTreeLeaf) String() string {

	if l == nil {
		return "nil"
	}
	return fmt.Sprintf("L1InfoTreeRoot:%s L1InfoTreeIndex:%d PreviousBlockHash:%s "+
		"BlockNumber:%d Timestamp:%s MainnetExitRoot:%s RollupExitRoot:%s GlobalExitRoot:%s",
		l.L1InfoTreeRoot.String(), l.L1InfoTreeIndex, l.PreviousBlockHash.String(),
		l.BlockNumber, l.Timestamp.String(), l.MainnetExitRoot.String(),
		l.RollupExitRoot.String(), l.GlobalExitRoot.String())

}
