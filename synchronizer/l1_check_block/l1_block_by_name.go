package l1_check_block

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	L1BlockPointWithOffsetDelimiter = "/"
)

// L1BlockPoint is an enum that represents the point of the L1 block
type L1BlockPoint int

type L1BlockPointWithOffset struct {
	BlockPoint L1BlockPoint
	Offset     int
}

func (v L1BlockPointWithOffset) String() string {
	return fmt.Sprintf("%s%s%d", v.BlockPoint.ToString(), L1BlockPointWithOffsetDelimiter, v.Offset)
}

func (p *L1BlockPointWithOffset) GreaterThan(other *L1BlockPointWithOffset) bool {
	if p == nil || other == nil {
		return false
	}
	if p.BlockPoint < other.BlockPoint {
		return true
	}
	if p.BlockPoint == other.BlockPoint {
		return p.Offset > other.Offset
	}
	return false
}

const (
	// FinalizedBlockNumber is the finalized block number
	FinalizedBlockNumber L1BlockPoint = 3
	// SafeBlockNumber is the safe block number
	SafeBlockNumber L1BlockPoint = 2
	// PendingBlockNumber is the pending block number
	PendingBlockNumber L1BlockPoint = 1
	// LastBlockNumber is the last block number
	LastBlockNumber L1BlockPoint = 0
)

var (
	// L1SafeFetch object preconfigured to fetch the safe block number
	L1SafeFetch = L1BlockNumberByNameFetch{L1BlockPointWithOffset{BlockPoint: SafeBlockNumber, Offset: 0}, false}
	// L1FinalizedFetch object preconfigured to fetch the finalized block number
	L1FinalizedFetch = L1BlockNumberByNameFetch{L1BlockPointWithOffset{BlockPoint: FinalizedBlockNumber, Offset: 0}, false}
	// L1PendingFetch object preconfigured to fetch the pending block number
	L1PendingFetch = L1BlockNumberByNameFetch{L1BlockPointWithOffset{BlockPoint: PendingBlockNumber, Offset: 0}, false}
	// L1LastestFetch object preconfigured to fetch the latest block number
	L1LastestFetch = L1BlockNumberByNameFetch{L1BlockPointWithOffset{BlockPoint: LastBlockNumber, Offset: 0}, false}
)

// ToString converts a L1BlockPoint to a string
func (v L1BlockPoint) ToString() string {
	switch v {
	case FinalizedBlockNumber:
		return "finalized"
	case SafeBlockNumber:
		return "safe"
	case PendingBlockNumber:
		return "pending"
	case LastBlockNumber:
		return "latest"
	}
	return "Unknown"
}

// StringToL1BlockPoint converts a string to a L1BlockPoint
func StringToL1BlockPointWithOffset(s string) (L1BlockPointWithOffset, error) {
	result := L1BlockPointWithOffset{}
	splitted := strings.Split(s, L1BlockPointWithOffsetDelimiter)
	if len(splitted) > 2 {
		return result, fmt.Errorf("invalid L1BlockPointWithOffset string: %s (only 1 delimiter permitted)", s)
	}
	l1block, err := StringToL1BlockPoint(splitted[0])
	if err != nil {
		return result, err
	}
	result.BlockPoint = l1block
	if len(splitted) == 2 {
		offset, err := fmt.Sscanf(splitted[1], "%d", &result.Offset)
		if err != nil {
			return result, fmt.Errorf("invalid L1BlockPointWithOffset string: %s (offset must be an integer)", s)
		}
		if offset != 1 {
			return result, fmt.Errorf("invalid L1BlockPointWithOffset string: %s (only 1 offset permitted)", s)
		}
	}
	return result, nil
}

// StringToL1BlockPoint converts a string to a L1BlockPoint
func StringToL1BlockPoint(s string) (L1BlockPoint, error) {
	switch s {
	case "finalized":
		return FinalizedBlockNumber, nil
	case "safe":
		return SafeBlockNumber, nil
	case "pending":
		return PendingBlockNumber, nil
	case "latest":
		return LastBlockNumber, nil
	default:
		return FinalizedBlockNumber, fmt.Errorf("invalid L1BlockPoint string: %s", s)
	}
}

// ToGethRequest converts a L1BlockPoint to a big.Int used for request to GETH
func (v L1BlockPoint) ToGethRequest() *big.Int {
	switch v {
	case FinalizedBlockNumber:
		return big.NewInt(int64(rpc.FinalizedBlockNumber))
	case PendingBlockNumber:
		return big.NewInt(int64(rpc.PendingBlockNumber))
	case SafeBlockNumber:
		return big.NewInt(int64(rpc.SafeBlockNumber))
	case LastBlockNumber:
		return nil
	}
	return big.NewInt(int64(v))
}

// L1BlockNumberByNameFetch  implements a safe L1 block number fetch
type L1BlockNumberByNameFetch struct {
	L1BlockPointWithOffset
	IfNotFoundReturnsZeroFlag bool
}

// NewL1BlockNumberByNameFetch creates a new SafeL1BlockNumberFetch
func NewL1BlockNumberByNameFetch(safeBlockPointWithOffset L1BlockPointWithOffset) *L1BlockNumberByNameFetch {
	res := L1BlockNumberByNameFetch{L1BlockPointWithOffset: safeBlockPointWithOffset}
	return &res
}

func (p *L1BlockNumberByNameFetch) SetIfNotFoundReturnsZero() *L1BlockNumberByNameFetch {
	if p == nil {
		return p
	}
	p.IfNotFoundReturnsZeroFlag = true
	return p
}

// Description returns a string representation of SafeL1BlockNumberFetch
func (p *L1BlockNumberByNameFetch) Description() string {
	if p == nil {
		return "nil"
	}
	return (*p).String()
}

// BlockNumber gets the safe block number from L1
func (p *L1BlockNumberByNameFetch) BlockNumber(ctx context.Context, requester L1Requester) (uint64, error) {
	l1SafePointBlock, err := requester.HeaderByNumber(ctx, p.BlockPoint.ToGethRequest())
	blockNumber := uint64(0)
	if err != nil {
		if strings.Contains(err.Error(), "block not found") && p.IfNotFoundReturnsZeroFlag {
			log.Warnf("block %s not found, assuming 0", p.String())
			return blockNumber, nil
		} else {
			log.Errorf("%s: Error getting L1 block %d. err: %s", logPrefix, p.String(), err.Error())
			return uint64(0), err
		}
	}
	blockNumber = l1SafePointBlock.Number.Uint64()
	if p.Offset < 0 {
		if blockNumber < uint64(-p.Offset) {
			blockNumber = 0
		} else {
			blockNumber += uint64(p.Offset)
		}
	} else {
		blockNumber = l1SafePointBlock.Number.Uint64() + uint64(p.Offset)
	}
	if p.BlockPoint == LastBlockNumber {
		blockNumber = min(blockNumber, l1SafePointBlock.Number.Uint64())
	}

	return blockNumber, nil
}
