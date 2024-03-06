package pgstorage

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
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

func (p *PostgresStorage) AddL1InfoTreeLeaf(ctx context.Context, exitRoot *L1InfoTreeLeaf, dbTx pgx.Tx) error {
	const addGlobalExitRootSQL = `
		INSERT INTO sync.exit_root(block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGlobalExitRootSQL,
		exitRoot.BlockNumber, exitRoot.Timestamp, exitRoot.MainnetExitRoot.String(), exitRoot.RollupExitRoot.String(),
		exitRoot.GlobalExitRoot.String(), exitRoot.PreviousBlockHash.String(), exitRoot.L1InfoTreeRoot.String(), exitRoot.L1InfoTreeIndex)
	return err
}

func (p *PostgresStorage) GetAllL1InfoTreeLeaves(ctx context.Context, dbTx pgx.Tx) ([]L1InfoTreeLeaf, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM sync.exit_root 
		WHERE l1_info_tree_index IS NOT NULL
		ORDER BY l1_info_tree_index`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getL1InfoRootSQL)
	if errors.Is(err, pgx.ErrNoRows) {
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

func (p *PostgresStorage) GetLatestL1InfoTreeLeaf(ctx context.Context, dbTx pgx.Tx) (*L1InfoTreeLeaf, error) {
	const getLatestL1InfoTreeLeafSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM sync.exit_root 
		WHERE l1_info_tree_index IS NOT NULL
		ORDER BY l1_info_tree_index DESC LIMIT 1`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getLatestL1InfoTreeLeafSQL)
	entry, err := scanL1InfoTreeExitRootStorageEntry(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &entry, err
}

func (p *PostgresStorage) GetL1InfoLeafPerIndex(ctx context.Context, L1InfoTreeIndex uint32, dbTx pgx.Tx) (*L1InfoTreeLeaf, error) {
	const getL1InfoLeafPerIndexSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM sync.exit_root 
		WHERE l1_info_tree_index = $1`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getL1InfoLeafPerIndexSQL, L1InfoTreeIndex)
	entry, err := scanL1InfoTreeExitRootStorageEntry(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &entry, err
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
		entry.L1InfoTreeRoot = common.HexToHash(L1InfoTreeRoot)
		entry.PreviousBlockHash = common.HexToHash(PreviousBlockHash)
		entry.MainnetExitRoot = common.HexToHash(MainnetExitRoot)
		entry.RollupExitRoot = common.HexToHash(RollupExitRoot)
		entry.GlobalExitRoot = common.HexToHash(GlobalExitRoot)
		return entry, err
	}
	return entry, nil
}
