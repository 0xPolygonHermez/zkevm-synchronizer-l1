package pgstorage

import (
	"context"
	"fmt"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

var (
	tableVirtualBatch           = "sync.virtual_batch"
	mandatoryFieldsVirtualBatch = []string{"batch_num", "fork_id", "raw_txs_data", "vlog_tx_hash", "coinbase", "sequence_from_batch_num", "block_num",
		"sequencer_addr", "received_at", "sync_version"}
	optionalFieldsVirtualBatch = []string{"l1_info_root", "extra_info", "batch_timestamp"}
)

// AddVirtualBatch adds a new virtual batch to the storage.
func (p *PostgresStorage) AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx dbTxType) error {
	mandatoryArguments := []interface{}{virtualBatch.BatchNumber, virtualBatch.ForkID, virtualBatch.BatchL2Data, virtualBatch.VlogTxHash.String(),
		virtualBatch.Coinbase.String(), virtualBatch.SequenceFromBatchNumber, virtualBatch.BlockNumber, virtualBatch.SequencerAddr.String(), virtualBatch.ReceivedAt, zkevm_synchronizer_l1.Version}

	var l1inforoot *string
	if virtualBatch.L1InfoRoot != nil {
		tmp := virtualBatch.L1InfoRoot.String()
		l1inforoot = &tmp
	}
	optionalArguments := []interface{}{l1inforoot, virtualBatch.ExtraInfo, virtualBatch.BatchTimestamp}
	fields := append(mandatoryFieldsVirtualBatch, optionalFieldsVirtualBatch...)
	arguments := append(mandatoryArguments, optionalArguments...)
	sql := composeInsertSql(fields, tableVirtualBatch)
	e := p.getExecQuerier(getPgTx(dbTx))
	_, err := e.Exec(ctx, sql, arguments...)
	err = translatePgxError(err, fmt.Sprintf("AddVirtualBatch %d", virtualBatch.BatchNumber))
	return err

}

func (p *PostgresStorage) GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*VirtualBatch, error) {
	fields := append(mandatoryFieldsVirtualBatch, optionalFieldsVirtualBatch...)
	sql := composeSelectSql(fields, tableVirtualBatch, "batch_num = $1")
	e := p.getExecQuerier(getPgTx(dbTx))
	row := e.QueryRow(ctx, sql, batchNumber)
	return scanVirtualBatch(row, fmt.Sprintf("GetVirtualBatchByBatchNumber %d", batchNumber))
}

func scanVirtualBatch(row pgx.Row, contextDescription string) (*VirtualBatch, error) {
	virtualBatch := &VirtualBatch{}
	var l1InfoRootStr *string
	var batchTimestamp *time.Time
	var syncVersion string
	var vlogTxHash string
	var coinbase string
	var sequencerAddr string
	err := row.Scan(&virtualBatch.BatchNumber, &virtualBatch.ForkID, &virtualBatch.BatchL2Data, &vlogTxHash, &coinbase,
		&virtualBatch.SequenceFromBatchNumber, &virtualBatch.BlockNumber, &sequencerAddr, &virtualBatch.ReceivedAt, &syncVersion,
		&l1InfoRootStr, &virtualBatch.ExtraInfo, &batchTimestamp)
	err = translatePgxError(err, contextDescription)
	if err != nil {
		return nil, err
	}
	virtualBatch.VlogTxHash = common.HexToHash(vlogTxHash)
	virtualBatch.Coinbase = common.HexToAddress(coinbase)
	virtualBatch.SequencerAddr = common.HexToAddress(sequencerAddr)
	if l1InfoRootStr != nil {
		l1InfoRoot := common.HexToHash(*l1InfoRootStr)
		virtualBatch.L1InfoRoot = &l1InfoRoot
	}
	if batchTimestamp != nil {
		virtualBatch.BatchTimestamp = batchTimestamp
	}
	return virtualBatch, nil
}

func composeSelectSql(fields []string, tableName string, whereStaments string) string {
	sql := "SELECT "
	for i, field := range fields {
		sql += field
		if i < len(fields)-1 {
			sql += ", "
		}
	}
	sql += " FROM " + tableName
	if whereStaments != "" {
		sql += " WHERE " + whereStaments
	}
	return sql
}

func composeInsertSql(fields []string, tableName string) string {
	sql := "INSERT INTO " + tableName + " ("
	for i, field := range fields {
		sql += field
		if i < len(fields)-1 {
			sql += ", "
		}
	}
	sql += ") VALUES ("
	for i := range fields {
		sql += fmt.Sprintf("$%d", i+1)
		if i < len(fields)-1 {
			sql += ", "
		}
	}
	sql += ")"
	return sql
}
