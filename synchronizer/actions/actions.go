package actions

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

var (
	// ErrInvalidParams is used when the object is not found
	ErrInvalidParams = errors.New("invalid params")
)

// L1EventProcessor is the interface for a processor of L1 events
// The main function is Process that must execute the event
type L1EventProcessor interface {
	// Name of the processor
	Name() string
	// SupportedForkIds list of forkId that support (you could use WildcardForkId)
	SupportedForkIds() []ForkIdType
	// SupportedEvents list of events that support (typically one)
	SupportedEvents() []etherman.EventOrder
	// Process a incomming event
	Process(ctx context.Context, forkId ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error
}
