package state

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/l1infotree"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type StorageL1InfoTreeInterface interface {
	AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx pgx.Tx) error
	GetAllL1InfoTreeLeaves(ctx context.Context, dbTx pgx.Tx) ([]pgstorage.L1InfoTreeLeaf, error)
	GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
	GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error)
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

func HashLeaf(leaf *pgstorage.L1InfoTreeLeaf) common.Hash {
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
		leaves = append(leaves, HashLeaf(&leaf))
	}
	mt, err := l1infotree.NewL1InfoTree(uint8(32), leaves) //nolint:gomnd
	if err != nil {
		log.Error("error creating L1InfoTree. Error: ", err)
		return fmt.Errorf("error creating L1InfoTree. Error: %w", err)
	}
	s.l1InfoTree = mt
	return nil
}

func (s *L1InfoTreeState) AddL1InfoTreeLeaf(ctx context.Context, exitRoot *pgstorage.L1InfoTreeLeaf, dbTx pgx.Tx) (*pgstorage.L1InfoTreeLeaf, error) {
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
	entry := *exitRoot
	entry.L1InfoTreeRoot = root
	entry.L1InfoTreeIndex = newIndex
	err = s.storage.AddL1InfoTreeLeaf(ctx, &entry, dbTx)
	if err != nil {
		log.Error("error adding L1InfoRoot to ExitRoot. Error: ", err)
		return nil, err
	}
	return &entry, nil
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
