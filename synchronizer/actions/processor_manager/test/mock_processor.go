package processor_manager_test

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/etherman"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/actions"
)

type ProcessorStub struct {
	name             string
	supportedEvents  []etherman.EventOrder
	supportedForkIds []actions.ForkIdType
	responseProcess  error
}

func (p *ProcessorStub) Name() string {
	return p.name
}

func (p *ProcessorStub) SupportedEvents() []etherman.EventOrder {
	return p.supportedEvents
}

func (p *ProcessorStub) SupportedForkIds() []actions.ForkIdType {
	return p.supportedForkIds
}

func (p *ProcessorStub) Process(ctx context.Context, forkId actions.ForkIdType, order etherman.Order, l1Block *etherman.Block, dbTx entities.Tx) error {
	return p.responseProcess
}
