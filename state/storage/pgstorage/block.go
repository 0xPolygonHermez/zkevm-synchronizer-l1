package pgstorage

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *L1Block, dbTx dbTxType) error {
	const addBlockSQL = "INSERT INTO sync.block (block_num, block_hash, parent_hash, received_at,checked,has_events, sync_version) VALUES ($1, $2, $3, $4, $5, $6,$7)"

	e := p.getExecQuerier(getPgTx(dbTx))
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt, block.Checked, block.HasEvents, block.SyncVersion)
	return translatePgxError(err, fmt.Sprintf("AddBlock %d", block.Key()))
}

// GetLastBlock returns the last L1 block.
func (p *PostgresStorage) GetLastBlock(ctx context.Context, dbTx dbTxType) (*L1Block, error) {
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked, has_events,sync_version FROM sync.block ORDER BY block_num DESC LIMIT 1"
	return p.queryBlock(ctx, "GetLastBlock", getLastBlockSQL, dbTx)
}

// GetBlockByNumber returns the L1 block with the given number.
func (p *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx dbTxType) (*L1Block, error) {
	const getBlockByNumberSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked,has_events,sync_version FROM sync.block WHERE block_num = $1"
	return p.queryBlock(ctx, fmt.Sprintf("GetBlockByNumber %d", blockNumber), getBlockByNumberSQL, dbTx, blockNumber)
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
// 0 is latest or fromBlockNumber
// 1 is the previous block to latest or to fromBlockNumber
// so on...
func (p *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, fromBlockNumber *uint64, dbTx dbTxType) (*L1Block, error) {
	whereClause := ""
	if fromBlockNumber != nil {
		whereClause = "WHERE block_num <= " + fmt.Sprintf("%d", *fromBlockNumber)
	}
	getPreviousBlockSQL := "SELECT block_num, block_hash, parent_hash, received_at,checked,has_events,sync_version FROM sync.block " + whereClause + " ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	return p.queryBlock(ctx, fmt.Sprintf("GetPreviousBlock %d", offset), getPreviousBlockSQL, dbTx, offset)
}

// GetFirstUncheckedBlock returns the first L1 block that has not been checked from a given block number.
func (p *PostgresStorage) GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx dbTxType) (*L1Block, error) {
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at, has_events,checked FROM sync.block  WHERE block_num>=$1 AND  checked=false ORDER BY block_num LIMIT 1"
	return p.queryBlock(ctx, "GetFirstUncheckedBlock", getLastBlockSQL, dbTx, fromBlockNumber)
}

// GetUncheckedBlocks returns all the unchecked blocks between fromBlockNumber and toBlockNumber (both included).
func (p *PostgresStorage) GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx dbTxType) (*[]L1Block, error) {
	const getUncheckedBlocksSQL = "SELECT block_num, block_hash, parent_hash, received_at,has_events, checked FROM sync.block WHERE block_num>=$1 AND block_num<=$2 AND checked=false ORDER BY block_num"
	return p.queryBlocks(ctx, "GetUncheckedBlocks", getUncheckedBlocksSQL, getPgTx(dbTx), fromBlockNumber, toBlockNumber)
}
func (p *PostgresStorage) queryBlocks(ctx context.Context, desc string, sql string, dbTx pgx.Tx, args ...interface{}) (*[]L1Block, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []L1Block
	for rows.Next() {
		block, err := scanBlock(rows)
		if err != nil {
			err = translatePgxError(err, desc)
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return &blocks, nil
}

func (p *PostgresStorage) queryBlock(ctx context.Context, desc string, sql string, dbTx dbTxType, args ...interface{}) (*L1Block, error) {
	q := p.getExecQuerier(getPgTx(dbTx))
	row := q.QueryRow(ctx, sql, args...)
	block, err := scanBlock(row)
	err = translatePgxError(err, desc)
	return &block, err
}

func scanBlock(row pgx.Row) (L1Block, error) {
	var (
		blockHash  string
		parentHash string
	)
	block := L1Block{}
	if err := row.Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked, &block.HasEvents, &block.SyncVersion); err != nil {
		return block, err
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return block, nil
}
