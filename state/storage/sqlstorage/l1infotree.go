package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type L1InfoTreeLeaf = entities.L1InfoTreeLeaf

const exitRootTable = "exit_root"

func (p *SqlStorage) AddL1InfoTreeLeaf(ctx context.Context, exitRoot *L1InfoTreeLeaf, dbTx dbTxType) error {
	addGlobalExitRootSQL := "INSERT INTO " + p.BuildTableName(exitRootTable) + "(block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, addGlobalExitRootSQL,
		exitRoot.BlockNumber, exitRoot.Timestamp.UTC(), exitRoot.MainnetExitRoot.String(), exitRoot.RollupExitRoot.String(),
		exitRoot.GlobalExitRoot.String(), exitRoot.PreviousBlockHash.String(), exitRoot.L1InfoTreeRoot.String(), exitRoot.L1InfoTreeIndex)
	err = translateSqlError(err, "AddL1InfoTreeLeaf")
	return err
}

func (p *SqlStorage) GetAllL1InfoTreeLeaves(ctx context.Context, dbTx dbTxType) ([]L1InfoTreeLeaf, error) {
	getL1InfoRootSQL := "SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index " +
		"FROM " + p.BuildTableName(exitRootTable) + " " +
		"WHERE l1_info_tree_index IS NOT NULL " +
		"ORDER BY l1_info_tree_index"

	e := p.getExecQuerier(getSqlTx(dbTx))
	rows, err := e.QueryContext(ctx, getL1InfoRootSQL)
	err = translateSqlError(err, "GetAllL1InfoTreeLeaves")
	if errors.Is(err, entities.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []L1InfoTreeLeaf
	for rows.Next() {
		entry, err := scanL1InfoTreeExitRootStorageEntry(rows)

		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (p *SqlStorage) GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx dbTxType) (*L1InfoTreeLeaf, error) {
	getLatestL1InfoTreeLeafSQL := "SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index " +
		"FROM " + p.BuildTableName(exitRootTable) + " " +
		"WHERE l1_info_tree_index IS NOT NULL " +
		"ORDER BY l1_info_tree_index DESC LIMIT 1"
	e := p.getExecQuerier(getSqlTx(dbTx))
	row := e.QueryRowContext(ctx, getLatestL1InfoTreeLeafSQL)
	entry, err := scanL1InfoTreeExitRootStorageEntry(row)
	err = translateSqlError(err, "GetLatestL1InfoTreeLeaf")
	if errors.Is(err, entities.ErrNotFound) {
		return nil, nil
	}
	return &entry, err
}

func (p *SqlStorage) GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx dbTxType) (*L1InfoTreeLeaf, error) {
	getL1InfoLeafPerIndexSQL := "SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index " +
		"FROM " + p.BuildTableName(exitRootTable) + " " +
		"WHERE l1_info_tree_index = $1"
	e := p.getExecQuerier(getSqlTx(dbTx))
	row := e.QueryRowContext(ctx, getL1InfoLeafPerIndexSQL, L1InfoTreeIndex)
	entry, err := scanL1InfoTreeExitRootStorageEntry(row)
	err = translateSqlError(err, "GetL1InfoLeafPerIndex")
	if errors.Is(err, entities.ErrNotFound) {
		return nil, nil
	}
	return &entry, err
}

func (p *SqlStorage) GetLeafsByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx dbTxType) ([]L1InfoTreeLeaf, error) {
	getLeafsByL1InfoRootSQL := "SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index " +
		"FROM " + p.BuildTableName(exitRootTable) + " " +
		"WHERE l1_info_tree_index IS NOT NULL AND l1_info_tree_index <= (SELECT l1_info_tree_index FROM " + p.BuildTableName(exitRootTable) + " WHERE l1_info_root=$1) " +
		"ORDER BY l1_info_tree_index ASC"
	e := p.getExecQuerier(getSqlTx(dbTx))
	rows, err := e.QueryContext(ctx, getLeafsByL1InfoRootSQL, l1InfoRoot.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]L1InfoTreeLeaf, 0)

	for rows.Next() {
		entry, err := scanL1InfoTreeExitRootStorageEntry(rows)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("l1InfoRoot:%s  Err: %w", l1InfoRoot.String(), entities.ErrNotFound)
	}

	return entries, nil
}

func scanL1InfoTreeExitRootStorageEntry(row pgx.Row) (L1InfoTreeLeaf, error) {
	var (
		L1InfoTreeRoot    string
		PreviousBlockHash string
		MainnetExitRoot   string
		RollupExitRoot    string
		GlobalExitRoot    string
	)
	entry := L1InfoTreeLeaf{}

	if err := row.Scan(
		&entry.BlockNumber, &entry.Timestamp, &MainnetExitRoot, &RollupExitRoot, &GlobalExitRoot,
		&PreviousBlockHash, &L1InfoTreeRoot, &entry.L1InfoTreeIndex); err != nil {
		return entry, err
	}
	entry.L1InfoTreeRoot = common.HexToHash(L1InfoTreeRoot)
	entry.PreviousBlockHash = common.HexToHash(PreviousBlockHash)
	entry.MainnetExitRoot = common.HexToHash(MainnetExitRoot)
	entry.RollupExitRoot = common.HexToHash(RollupExitRoot)
	entry.GlobalExitRoot = common.HexToHash(GlobalExitRoot)
	return entry, nil
}
