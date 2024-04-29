package pgstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// L1Block struct
type L1Block struct {
	BlockNumber uint64
	BlockHash   common.Hash
	ParentHash  common.Hash
	ReceivedAt  time.Time
	Checked     bool // The block is safe (have past the safe point, e.g. Finalized in L1)
	SyncVersion string
}

func (b *L1Block) String() string {
	if b == nil {
		return "nil"
	}
	return fmt.Sprintf("BlockNumber: %d, BlockHash: %s, ParentHash: %s, ReceivedAt: %s, Checked: %t, SyncVersion: %s",
		b.BlockNumber, b.BlockHash.String(), b.ParentHash.String(), b.ReceivedAt.String(), b.Checked, b.SyncVersion)
}

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *L1Block, dbTx pgx.Tx) error {
	const addBlockSQL = "INSERT INTO sync.block (block_num, block_hash, parent_hash, received_at,checked, sync_version) VALUES ($1, $2, $3, $4, $5, $6)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt, block.Checked, block.SyncVersion)
	return err
}

// GetLastBlock returns the last L1 block.
func (p *PostgresStorage) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*L1Block, error) {
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked, sync_version FROM sync.block ORDER BY block_num DESC LIMIT 1"
	return p.queryBlock(ctx, getLastBlockSQL, dbTx)
}

// GetBlockByNumber returns the L1 block with the given number.
func (p *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*L1Block, error) {
	const getBlockByNumberSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked,sync_version FROM sync.block WHERE block_num = $1"
	return p.queryBlock(ctx, getBlockByNumberSQL, dbTx, blockNumber)
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
func (p *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*L1Block, error) {
	const getPreviousBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked,sync_version FROM sync.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	return p.queryBlock(ctx, getPreviousBlockSQL, dbTx, offset)
}

func (p *PostgresStorage) queryBlock(ctx context.Context, sql string, dbTx pgx.Tx, args ...interface{}) (*L1Block, error) {
	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, sql, args...)
	block, err := scanBlock(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &block, err
}

func scanBlock(row pgx.Row) (L1Block, error) {
	var (
		blockHash  string
		parentHash string
	)
	block := L1Block{}
	if err := row.Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked, &block.SyncVersion); err != nil {
		return block, err
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return block, nil
}
