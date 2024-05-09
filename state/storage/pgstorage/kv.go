package pgstorage

import (
	"context"
	"encoding/json"
)

func (p *PostgresStorage) KVSetString(ctx context.Context, key string, value string, dbTx dbTxType) error {
	e := p.getExecQuerier(getPgTx(dbTx))
	const setSQL = "INSERT INTO sync.kv (key, value) VALUES ($1, $2) ON CONFLICT (key) DO UPDATE SET value = $2"
	if _, err := e.Exec(ctx, setSQL, key, value); err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) KVSetJson(ctx context.Context, key string, value interface{}, dbTx dbTxType) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.KVSetString(ctx, key, string(jsonValue), dbTx)
}

func (p *PostgresStorage) KVSetUint64(ctx context.Context, key string, value uint64, dbTx dbTxType) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return p.KVSetString(ctx, key, string(jsonValue), dbTx)
}

func (p *PostgresStorage) KVGetString(ctx context.Context, key string, dbTx dbTxType) (string, error) {
	e := p.getExecQuerier(getPgTx(dbTx))
	const getSQL = "SELECT value FROM sync.kv WHERE key = $1"
	row := e.QueryRow(ctx, getSQL, key)
	var value string
	err := row.Scan(&value)
	err = translatePgxError(err, "KVGetString")
	if err != nil {
		return "", err
	}
	return value, nil
}

func (p *PostgresStorage) KVGetJson(ctx context.Context, key string, value interface{}, dbTx dbTxType) error {
	valueStr, err := p.KVGetString(ctx, key, dbTx)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(valueStr), value)
}

func (p *PostgresStorage) KVGetUint64(ctx context.Context, key string, dbTx dbTxType) (uint64, error) {
	valueStr, err := p.KVGetString(ctx, key, dbTx)
	if err != nil {
		return 0, err
	}
	value := uint64(0)
	err = json.Unmarshal([]byte(valueStr), &value)
	return value, err
}

func (p *PostgresStorage) KVExists(ctx context.Context, key string, dbTx dbTxType) (bool, error) {
	e := p.getExecQuerier(getPgTx(dbTx))
	const existsSQL = "SELECT EXISTS(SELECT 1 FROM sync.kv WHERE key = $1)"
	row := e.QueryRow(ctx, existsSQL, key)
	var exists bool
	err := row.Scan(&exists)
	err = translatePgxError(err, "KVExists")
	if err != nil {
		return false, err
	}
	return exists, nil
}
