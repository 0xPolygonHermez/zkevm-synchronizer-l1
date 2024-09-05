package sqlstorage

import (
	"context"
	"fmt"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type VirtualBatch = entities.VirtualBatch

var (
	virtualBatchTable           = "virtual_batch"
	mandatoryFieldsVirtualBatch = []string{"batch_num", "fork_id", "raw_txs_data", "vlog_tx_hash", "coinbase", "sequence_from_batch_num", "block_num",
		"sequencer_addr", "received_at", "sync_version"}
	optionalFieldsVirtualBatch = []string{"l1_info_root", "extra_info", "batch_timestamp"}
)

// AddVirtualBatch adds a new virtual batch to the storage.
func (p *SqlStorage) AddVirtualBatch(ctx context.Context, virtualBatch *VirtualBatch, dbTx dbTxType) error {
	mandatoryArguments := []interface{}{virtualBatch.BatchNumber, virtualBatch.ForkID, virtualBatch.BatchL2Data, virtualBatch.VlogTxHash.String(),
		virtualBatch.Coinbase.String(), virtualBatch.SequenceFromBatchNumber, virtualBatch.BlockNumber, virtualBatch.SequencerAddr.String(), virtualBatch.ReceivedAt.UTC(), zkevm_synchronizer_l1.Version}

	var l1inforoot *string
	if virtualBatch.L1InfoRoot != nil {
		tmp := virtualBatch.L1InfoRoot.String()
		l1inforoot = &tmp
	}
	var tmpBatchTimestamp *time.Time
	if virtualBatch.BatchTimestamp != nil {
		utcTime := virtualBatch.BatchTimestamp.UTC()
		tmpBatchTimestamp = &utcTime
	}
	optionalArguments := []interface{}{l1inforoot, virtualBatch.ExtraInfo, tmpBatchTimestamp}
	fields := append(mandatoryFieldsVirtualBatch, optionalFieldsVirtualBatch...)
	arguments := append(mandatoryArguments, optionalArguments...)
	sql := composeInsertSql(fields, p.BuildTableName(virtualBatchTable))
	e := p.getExecQuerier(getSqlTx(dbTx))
	_, err := e.ExecContext(ctx, sql, arguments...)
	err = translateSqlError(err, fmt.Sprintf("AddVirtualBatch %d", virtualBatch.BatchNumber))
	return err

}

func (p *SqlStorage) GetVirtualBatchByBatchNumber(ctx context.Context, batchNumber uint64, dbTx dbTxType) (*VirtualBatch, error) {
	fields := append(mandatoryFieldsVirtualBatch, optionalFieldsVirtualBatch...)
	sql := composeSelectSql(fields, p.BuildTableName(virtualBatchTable), "batch_num = $1")
	e := p.getExecQuerier(getSqlTx(dbTx))
	row := e.QueryRowContext(ctx, sql, batchNumber)
	return scanVirtualBatch(row, fmt.Sprintf("GetVirtualBatchByBatchNumber %d", batchNumber))
}

func (p *SqlStorage) GetLastestVirtualBatchNumber(ctx context.Context, constrains *VirtualBatchConstraints, dbTx dbTxType) (uint64, error) {
	whereClause := ""
	if constrains != nil {
		whereClause = constrains.WhereClause()
		if whereClause != "" {
			whereClause = "WHERE " + whereClause
		}
	}
	sql := "SELECT batch_num FROM " + p.BuildTableName(virtualBatchTable) + " ORDER BY batch_num " + whereClause + " DESC LIMIT 1"
	e := p.getExecQuerier(getSqlTx(dbTx))
	row := e.QueryRowContext(ctx, sql)
	var batchNumber uint64
	err := row.Scan(&batchNumber)
	err = translateSqlError(err, "GetLastestVirtualBatchNumber")
	if err != nil {
		return 0, err
	}
	return batchNumber, nil
}

// VirtualBatchConstraints is a struct that contains the constraints to filter the virtual batches.
// is ready to add constraints to the query.
type VirtualBatchConstraints struct {
	batchNumberEqual *uint64
	batchNumberGt    *uint64
	batchNumberLt    *uint64
}

func (c *VirtualBatchConstraints) BatchNumberEqual(batchNumber uint64) {
	c.batchNumberEqual = &batchNumber
}

func (c *VirtualBatchConstraints) BatchNumberGt(batchNumber uint64) {
	c.batchNumberGt = &batchNumber
}

func (c *VirtualBatchConstraints) BatchNumberLt(batchNumber uint64) {
	c.batchNumberLt = &batchNumber
}

func (c *VirtualBatchConstraints) WhereClause() string {
	res := ""
	if c.batchNumberEqual != nil {
		res += fmt.Sprintf("batch_num = %d ", *c.batchNumberEqual)
	}
	if c.batchNumberGt != nil {
		res += fmt.Sprintf("batch_num>%d ", *c.batchNumberEqual)
	}
	if c.batchNumberLt != nil {
		res += fmt.Sprintf("batch_num<%d ", *c.batchNumberEqual)
	}
	return res
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
	err = translateSqlError(err, contextDescription)
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
