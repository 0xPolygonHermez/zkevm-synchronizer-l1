package synchronizer

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer/common/syncinterfaces"
	"github.com/ethereum/go-ethereum/common"
)

type stateSyncQueries interface {
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx entities.Tx) ([]entities.L1InfoTreeLeaf, error)
	GetL1InfoRootPerLeafIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx entities.Tx) (common.Hash, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx entities.Tx) (*entities.L1InfoTreeLeaf, error)
	GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32, dbTx entities.Tx) (map[uint32]entities.L1InfoTreeLeaf, error)
}

type storageSyncQueries interface {
	syncinterfaces.StorageBlockReaderInterface
	syncinterfaces.StorageSequenceBatchesInterface
	syncinterfaces.StorageVirtualBatchInterface
	syncinterfaces.StorageBlockReaderInterface
}

type SyncrhronizerQueries struct {
	state   stateSyncQueries
	storage storageSyncQueries
	ctx     context.Context
}

func NewSyncrhronizerQueries(state stateSyncQueries, storage storageSyncQueries, ctx context.Context) *SyncrhronizerQueries {
	return &SyncrhronizerQueries{
		state:   state,
		storage: storage,
		ctx:     ctx,
	}
}

func (s *SyncrhronizerQueries) GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash) ([]L1InfoTreeLeaf, error) {
	leaves, err := s.state.GetLeafsByL1InfoRoot(ctx, l1InfoRoot, nil)
	if err != nil {
		log.Error("error getting leaves by L1InfoRoot. Error: ", err)
		return nil, err
	}
	var res []L1InfoTreeLeaf
	for _, leaf := range leaves {
		tmp := L1InfoTreeLeaf(leaf)
		res = append(res, tmp)
	}
	return res, nil
}

func (s *SyncrhronizerQueries) GetL1InfoRootPerIndex(ctx context.Context, L1InfoTreeIndex uint32) (common.Hash, error) {
	root, err := s.state.GetL1InfoRootPerLeafIndex(ctx, L1InfoTreeIndex, nil)
	if errors.Is(err, entities.ErrNotFound) {
		return common.Hash{}, ErrNotFound
	}
	return root, err
}

func (s *SyncrhronizerQueries) GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32) (map[uint32]L1InfoTreeLeaf, error) {
	leaves, err := s.state.GetL1InfoTreeLeaves(ctx, indexLeaves, nil)
	if err != nil {
		return nil, err
	}
	// Convert type state.L1InfoTreeLeaf to type L1InfoTreeLeaf
	returnLeaves := make(map[uint32]L1InfoTreeLeaf)
	for _, idx := range indexLeaves {
		returnLeaves[idx] = L1InfoTreeLeaf(leaves[idx])
	}
	return returnLeaves, nil
}

func (s *SyncrhronizerQueries) GetSequenceByBatchNumber(ctx context.Context, batchNumber uint64) (*SequencedBatches, error) {
	sequence, err := s.storage.GetSequenceByBatchNumber(ctx, batchNumber, nil)
	if sequence == nil {
		return nil, err
	}
	res := SequencedBatches(*sequence)
	return &res, err
}

func (s *SyncrhronizerQueries) GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64) (*VirtualBatch, error) {
	virtualBatch, err := s.storage.GetVirtualBatchByBatchNumber(ctx, batchNumber, nil)
	if virtualBatch == nil {
		return nil, err
	}
	res := VirtualBatch(*virtualBatch)
	return &res, err
}

func (s *SyncrhronizerQueries) GetLastestVirtualBatchNumber(ctx context.Context) (uint64, error) {
	lastBatchNumber, err := s.storage.GetLastestVirtualBatchNumber(ctx, nil, nil)
	if err != nil {
		return 0, err
	}
	return lastBatchNumber, nil
}

func (s *SyncrhronizerQueries) GetL1BlockByNumber(ctx context.Context, blockNumber uint64) (*L1Block, error) {
	block, err := s.storage.GetBlockByNumber(ctx, blockNumber, nil)
	if block == nil {
		return nil, err
	}
	res := L1Block(*block)
	return &res, err
}

func (s *SyncrhronizerQueries) GetLastL1Block(ctx context.Context) (*L1Block, error) {
	block, err := s.storage.GetLastBlock(ctx, nil)
	if block == nil {
		return nil, err
	}
	res := L1Block(*block)
	return &res, err
}
