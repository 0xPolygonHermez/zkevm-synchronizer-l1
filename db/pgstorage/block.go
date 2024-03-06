package pgstorage

import (
	"context"
	"errors"
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
}

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *L1Block, dbTx pgx.Tx) error {
	const addBlockSQL = "INSERT INTO sync.block (block_num, block_hash, parent_hash, received_at) VALUES ($1, $2, $3, $4)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt)
	return err
}

// GetLastBlock returns the last L1 block.
func (p *PostgresStorage) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*L1Block, error) {
	var (
		blockHash  string
		parentHash string
		block      L1Block
	)
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at FROM sync.block ORDER BY block_num DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetBlockByNumber returns the L1 block with the given number.
func (p *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*L1Block, error) {
	var (
		blockHash  string
		parentHash string
		block      L1Block
	)
	const getBlockByNumberSQL = "SELECT block_num, block_hash, parent_hash, received_at FROM sync.block WHERE block_num = $1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
func (p *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*L1Block, error) {
	var (
		blockHash  string
		parentHash string
		block      L1Block
	)
	const getPreviousBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at FROM sync.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}
