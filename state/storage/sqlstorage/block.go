package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type L1Block = entities.L1Block

const blockTable = "block"

// AddBlock adds a new block to the State Store
func (p *SqlStorage) AddBlock(ctx context.Context, block *L1Block, dbTx dbTxType) error {
	addBlockSQL := "INSERT INTO " + p.BuildTableName(blockTable) + " (block_num, block_hash, parent_hash, received_at,checked,has_events, sync_version) VALUES ($1, $2, $3, $4, $5, $6,$7)"

	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt.UTC(), block.Checked, block.HasEvents, block.SyncVersion)
	return translateSqlError(err, fmt.Sprintf("AddBlock %d", block.Key()))
}

// UpdateCheckedBlockByNumber update checked flag for a block
func (p *SqlStorage) UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx dbTxType) error {
	query := "UPDATE " + p.BuildTableName(blockTable) + " SET checked = $1 WHERE block_num = $2"

	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, query, newCheckedStatus, blockNumber)
	return err
}

// -- READ FUNCTIONS ---------------------------------
const selectSQLAllFieldsBlock = "SELECT block_num, block_hash, parent_hash, received_at,checked,has_events,sync_version FROM "

// GetLastBlock returns the last L1 block.
func (p *SqlStorage) GetLastBlock(ctx context.Context, dbTx dbTxType) (*L1Block, error) {
	getLastBlockSQL := selectSQLAllFieldsBlock + p.BuildTableName(blockTable) + " ORDER BY block_num DESC LIMIT 1"
	return p.queryBlock(ctx, "GetLastBlock", getLastBlockSQL, dbTx)
}

// GetBlockByNumber returns the L1 block with the given number.
func (p *SqlStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx dbTxType) (*L1Block, error) {
	getBlockByNumberSQL := selectSQLAllFieldsBlock + p.BuildTableName(blockTable) + " WHERE block_num = $1"
	return p.queryBlock(ctx, fmt.Sprintf("GetBlockByNumber %d", blockNumber), getBlockByNumberSQL, dbTx, blockNumber)
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
// 0 is latest or fromBlockNumber
// 1 is the previous block to latest or to fromBlockNumber
// so on...
func (p *SqlStorage) GetPreviousBlock(ctx context.Context, offset uint64, dbTx dbTxType) (*L1Block, error) {
	getPreviousBlockSQL := selectSQLAllFieldsBlock + p.BuildTableName(blockTable) + "  ORDER BY block_num DESC LIMIT 1 OFFSET $1"
	return p.queryBlock(ctx, fmt.Sprintf("GetPreviousBlock %d", offset), getPreviousBlockSQL, dbTx, offset)
}

// GetFirstUncheckedBlock returns the first L1 block that has not been checked from a given block number.
func (p *SqlStorage) GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx dbTxType) (*L1Block, error) {
	getLastBlockSQL := selectSQLAllFieldsBlock + p.BuildTableName(blockTable) + "  WHERE block_num>=$1 AND  checked=false ORDER BY block_num LIMIT 1"
	return p.queryBlock(ctx, "GetFirstUncheckedBlock", getLastBlockSQL, dbTx, fromBlockNumber)
}

// GetUncheckedBlocks returns all the unchecked blocks between fromBlockNumber and toBlockNumber (both included).
func (p *SqlStorage) GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx dbTxType) (*[]L1Block, error) {
	getUncheckedBlocksSQL := selectSQLAllFieldsBlock + p.BuildTableName(blockTable) + " WHERE block_num>=$1 AND block_num<=$2 AND checked=false ORDER BY block_num"
	return p.queryBlocks(ctx, "GetUncheckedBlocks", getUncheckedBlocksSQL, getSqlTx(dbTx), fromBlockNumber, toBlockNumber)
}

func (p *SqlStorage) queryBlocks(ctx context.Context, desc string, sql string, dbTx *sql.Tx, args ...interface{}) (*[]L1Block, error) {
	q := p.getExecQuerier(dbTx)
	rows, err := q.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []L1Block
	for rows.Next() {
		block, err := scanBlock(rows)
		if err != nil {
			err = translateSqlError(err, desc)
			return nil, err
		}
		blocks = append(blocks, block)
	}
	return &blocks, nil
}

func (p *SqlStorage) queryBlock(ctx context.Context, desc string, sql string, dbTx dbTxType, args ...interface{}) (*L1Block, error) {
	q := p.getExecQuerier(getSqlTx(dbTx))
	row := q.QueryRowContext(ctx, sql, args...)
	block, err := scanBlock(row)
	err = translateSqlError(err, desc)
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
