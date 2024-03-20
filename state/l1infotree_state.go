package state

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/l1infotree"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	// SkipL1InfoTreeLeaf is special  index that skip the change of GlobalExitRoot, so the value of this leaf is never used
	SkipL1InfoTreeLeaf = uint32(0)
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

type StorageL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx pgx.Tx) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
}

type L1InfoTreeState struct {
	storage    StorageL1InfoTreeInterface
	l1InfoTree *l1infotree.L1InfoTree
}

func NewL1InfoTreeManager(storage StorageL1InfoTreeInterface) *L1InfoTreeState {
	return &L1InfoTreeState{
		storage: storage,
	}
}

func HashLeaf(leaf *L1InfoTreeLeaf) common.Hash {
	timestamp := uint64(leaf.Timestamp.Unix())
	return l1infotree.HashLeafData(leaf.GlobalExitRoot, leaf.PreviousBlockHash, timestamp)
}

func (s *L1InfoTreeState) BuildL1InfoTreeCacheIfNeed(ctx context.Context, dbTx pgx.Tx) error {
	if s.l1InfoTree != nil {
		return nil
	}
	log.Debugf("Building L1InfoTree cache")
	allLeaves, err := s.storage.GetAllL1InfoTreeLeaves(ctx, dbTx)
	if err != nil {
		log.Error("error getting all leaves. Error: ", err)
		return fmt.Errorf("error getting all leaves. Error: %w", err)
	}
	var leaves [][32]byte
	for _, leaf := range allLeaves {
		tmp := L1InfoTreeLeaf(leaf)
		leaves = append(leaves, HashLeaf(&tmp))
	}
	mt, err := l1infotree.NewL1InfoTree(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTree = mt
	return nil
}

func (s *L1InfoTreeState) AddL1InfoTreeLeaf(ctx context.Context, exitRoot *L1InfoTreeLeaf, dbTx pgx.Tx) (*L1InfoTreeLeaf, error) {
	var newIndex uint32
	lastLeaf, err := s.storage.GetLatestL1InfoTreeLeaf(ctx, dbTx)
	if err != nil {
		log.Error("error getting latest l1InfoTree index. Error: ", err)
		return nil, err
	}
	if lastLeaf != nil {
		newIndex = lastLeaf.L1InfoTreeIndex + 1
	} else {
		newIndex = 0
	}
	err = s.BuildL1InfoTreeCacheIfNeed(ctx, dbTx)
	if err != nil {
		log.Error("error building L1InfoTree cache. Error: ", err)
		return nil, err
	}
	log.Debug("newIndex: ", newIndex)
	root, err := s.l1InfoTree.AddLeaf(newIndex, HashLeaf(exitRoot))
	if err != nil {
		log.Error("error add new leaf to the L1InfoTree. Error: ", err)
		return nil, err
	}
	entry := pgstorage.L1InfoTreeLeaf(*exitRoot)
	entry.L1InfoTreeRoot = root
	entry.L1InfoTreeIndex = newIndex
	err = s.storage.AddL1InfoTreeLeaf(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	tmp := L1InfoTreeLeaf(entry)
	return &tmp, nil
}

func (s *L1InfoTreeState) GetL1InfoRootPerLeafIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (common.Hash, error) {
	leaf, err := s.storage.GetL1InfoLeafPerIndex(ctx, L1InfoTreeIndex, dbTx)
	if err != nil {
		log.Error("error getting L1InfoRoot per leaf index. Error: ", err)
		return common.Hash{}, err
	}
	if leaf == nil {
		return common.Hash{}, ErrNotFound
	}
	return leaf.L1InfoTreeRoot, nil
}

// GetL1InfoTreeLeaves returns the required leaves for the L1InfoTree
func (s *L1InfoTreeState) GetL1InfoTreeLeaves(ctx context.Context, indexLeaves []uint32, dbTx pgx.Tx) (map[uint32]L1InfoTreeLeaf, error) {
	res := map[uint32]L1InfoTreeLeaf{}
	for _, idx := range indexLeaves {
		if idx == SkipL1InfoTreeLeaf {
			// Skip this value
			continue
		}
		if _, found := res[idx]; found {
			// Is already in the result map
			continue
		}
		leaf, err := s.storage.GetL1InfoLeafPerIndex(ctx, idx, dbTx)
		if err != nil {
			err = fmt.Errorf("error getting L1InfoTree leaf %d. Error: %w", idx, err)
			log.Errorf(err.Error())
			return nil, err
		}
		if leaf == nil {
			return nil, ErrNotFound
		}
		res[idx] = L1InfoTreeLeaf(*leaf)
	}
	return res, nil
}

func (s *L1InfoTreeState) GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]L1InfoTreeLeaf, error) {
	leaves, err := s.storage.GetLeafsByL1InfoRoot(ctx, l1InfoRoot, dbTx)
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
