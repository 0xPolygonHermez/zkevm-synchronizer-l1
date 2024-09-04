package sqlstorage

import (
	"context"
	"encoding/json"
	"time"

	zkevm_synchronizer_l1 "github.com/0xPolygonHermez/zkevm-synchronizer-l1"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

const kvTable = "kv"

type KVMetadataEntry = entities.KVMetadataEntry
type KVKey = entities.KVKey

func (p *SqlStorage) KVSetString(ctx context.Context, key KVKey, value string, metadata *KVMetadataEntry, dbTx dbTxType) error {
	if metadata == nil {
		timeNow := time.Now()
		metadata = &KVMetadataEntry{
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
			SyncVersion: zkevm_synchronizer_l1.Version,
		}
	}
	e := p.getExecQuerier(getSqlTx(dbTx))
	setSQL := "INSERT INTO " + p.BuildTableName(kvTable) + " (key, value, created_at, updated_at, sync_version) VALUES ($1, $2, $3,$4,$5) ON CONFLICT (key) " +
		"DO UPDATE SET value = $2,  updated_at=$4, sync_version=$5"
	if _, err := e.ExecContext(ctx, setSQL, key, value, metadata.CreatedAt.UTC(), metadata.UpdatedAt.UTC(), metadata.SyncVersion); err != nil {
		return err
	}
	return nil
}

func (p *SqlStorage) KVSetJson(ctx context.Context, key KVKey, value interface{}, metadata *KVMetadataEntry, dbTx dbTxType) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.KVSetString(ctx, key, string(jsonValue), metadata, dbTx)
}

func (p *SqlStorage) KVGetString(ctx context.Context, key KVKey, metadata *KVMetadataEntry, dbTx dbTxType) (string, error) {
	e := p.getExecQuerier(getSqlTx(dbTx))
	getSQL := "SELECT value, created_at, updated_at, sync_version FROM " + p.BuildTableName(kvTable) + " WHERE key = $1"
	storageMetaData := &KVMetadataEntry{}
	row := e.QueryRowContext(ctx, getSQL, key)
	var value string
	err := row.Scan(&value, &storageMetaData.CreatedAt, &storageMetaData.UpdatedAt, &storageMetaData.SyncVersion)
	err = translateSqlError(err, "KVGetString")
	if err != nil {
		return "", err
	}
	if metadata != nil {
		*metadata = *storageMetaData
	}
	return value, nil
}

func (p *SqlStorage) KVGetJson(ctx context.Context, key KVKey, value interface{}, metadata *KVMetadataEntry, dbTx dbTxType) error {
	valueStr, err := p.KVGetString(ctx, key, metadata, dbTx)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(valueStr), value)
}

func (p *SqlStorage) KVExists(ctx context.Context, key KVKey, dbTx dbTxType) (bool, error) {
	e := p.getExecQuerier(getSqlTx(dbTx))
	existsSQL := "SELECT EXISTS(SELECT 1 FROM " + p.BuildTableName(kvTable) + " WHERE key = $1)"
	row := e.QueryRowContext(ctx, existsSQL, key)
	var exists bool
	err := row.Scan(&exists)
	err = translateSqlError(err, "KVExists")
	if err != nil {
		return false, err
	}
	return exists, nil
}
