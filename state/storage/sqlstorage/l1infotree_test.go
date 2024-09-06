package sqlstorage_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	testLeafData = &sqlstorage.L1InfoTreeLeaf{
		L1InfoTreeIndex:   1,
		L1InfoTreeRoot:    common.HexToHash("0x1"),
		PreviousBlockHash: common.HexToHash("0x2"),
		BlockNumber:       3,
		Timestamp:         time.Now().UTC(),
		MainnetExitRoot:   common.HexToHash("0x4"),
		RollupExitRoot:    common.HexToHash("0x5"),
		GlobalExitRoot:    common.HexToHash("0x6"),
	}
)

func TestAddL1InfoTreeLeafFailsForeignBlockKey(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	leaf := testLeafData
	err = storage.AddL1InfoTreeLeaf(ctx, leaf, dbTx)
	require.Error(t, err)
}

func TestGetLatestL1InfoTreeLeafErrorNotFoundReturnsNil(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	leaf, err := storage.GetLatestL1InfoTreeLeaf(ctx, dbTx)
	require.NoError(t, err)
	require.Nil(t, leaf)
}

func TestAddL1InfoTreeLeafOk(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block := &sqlstorage.L1Block{
		BlockNumber: 3,
	}
	err = storage.AddBlock(ctx, block, dbTx)
	require.NoError(t, err)
	leaf := testLeafData
	err = storage.AddL1InfoTreeLeaf(ctx, leaf, dbTx)
	require.NoError(t, err)
	leafRead, err := storage.GetL1InfoLeafPerIndex(ctx, leaf.L1InfoTreeIndex, dbTx)
	require.NoError(t, err)
	require.Equal(t, leaf, leafRead)
}

func TestGetAllL1InfoTreeLeaves(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block := &sqlstorage.L1Block{
		BlockNumber: 3,
	}
	err = storage.AddBlock(ctx, block, dbTx)
	require.NoError(t, err)
	leafs := []sqlstorage.L1InfoTreeLeaf{newTestRandomLeaf(1, 3), newTestRandomLeaf(2, 3), newTestRandomLeaf(3, 3)}
	for _, leaf := range leafs {
		err = storage.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
		require.NoError(t, err)
	}
	leafsRead, err := storage.GetAllL1InfoTreeLeaves(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, leafs, leafsRead)
}

func TestGetLatestL1InfoTreeLeaf(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block := &sqlstorage.L1Block{
		BlockNumber: 3,
	}
	err = storage.AddBlock(ctx, block, dbTx)
	require.NoError(t, err)
	leafs := []sqlstorage.L1InfoTreeLeaf{newTestRandomLeaf(3, 3), newTestRandomLeaf(1, 3), newTestRandomLeaf(2, 3)}
	for _, leaf := range leafs {
		err = storage.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
		require.NoError(t, err)
	}
	leafRead, err := storage.GetLatestL1InfoTreeLeaf(ctx, dbTx)
	require.NoError(t, err)
	require.NotNil(t, leafRead)
	require.Equal(t, leafs[0], *leafRead)
}

func TestGetLeafsByL1InfoRoot(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)
	defer func() { _ = dbTx.Commit(ctx) }()
	block := &sqlstorage.L1Block{
		BlockNumber: 3,
	}
	err = storage.AddBlock(ctx, block, dbTx)
	require.NoError(t, err)
	leafs := []sqlstorage.L1InfoTreeLeaf{newTestRandomLeaf(1, 3), newTestRandomLeaf(2, 3), newTestRandomLeaf(3, 3)}
	for _, leaf := range leafs {
		err = storage.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
		require.NoError(t, err)
	}
	leafsRead, err := storage.GetLeafsByL1InfoRoot(ctx, leafs[1].L1InfoTreeRoot, dbTx)
	require.NoError(t, err)
	require.NotNil(t, leafsRead)
	require.Equal(t, len(leafsRead), 2)
}

func TestL1InfoRootRemovedOnCascade(t *testing.T) {
	skipDatabaseTestIfNeeded(t)
	ctx := context.TODO()
	dbConfig := getStorageConfig()
	storage, err := sqlstorage.NewSqlStorage(dbConfig, true)
	require.NoError(t, err)
	dbTx, err := storage.BeginTransaction(ctx)
	require.NoError(t, err)

	blocks := []sqlstorage.L1Block{
		{BlockNumber: 1}, {BlockNumber: 2}, {BlockNumber: 3},
	}
	for _, block := range blocks {
		err = storage.AddBlock(ctx, &block, dbTx)
		require.NoError(t, err)
	}
	leafs := []sqlstorage.L1InfoTreeLeaf{newTestRandomLeaf(1, 1), newTestRandomLeaf(2, 2), newTestRandomLeaf(3, 3)}
	for _, leaf := range leafs {
		err = storage.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
		require.NoError(t, err)
	}
	err = storage.ResetToL1BlockNumber(ctx, 1, dbTx)
	require.NoError(t, err)
	leafsRead, err := storage.GetAllL1InfoTreeLeaves(ctx, dbTx)
	require.NoError(t, err)
	require.NotNil(t, leafsRead)
	require.Equal(t, 1, len(leafsRead))
	require.Equal(t, leafs[0], leafsRead[0])
}

func newTestRandomLeaf(index uint32, blockNumber uint64) sqlstorage.L1InfoTreeLeaf {
	return sqlstorage.L1InfoTreeLeaf{
		L1InfoTreeIndex:   index,
		L1InfoTreeRoot:    generateRandomHash(),
		PreviousBlockHash: generateRandomHash(),
		BlockNumber:       blockNumber,
		Timestamp:         time.Now().UTC(),
		MainnetExitRoot:   generateRandomHash(),
		RollupExitRoot:    generateRandomHash(),
		GlobalExitRoot:    generateRandomHash(),
	}
}

// GenerateRandomHash generates a random common.Hash
func generateRandomHash() common.Hash {
	var hash common.Hash
	_, err := rand.Read(hash[:])
	if err != nil {
		return common.Hash{}
	}
	return hash
}
