package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/mattn/go-sqlite3"
)

type dbTxType = entities.Tx

type SqlStorage struct {
	db *sql.DB
}

const (
	SqliteDriverName = "sqlite3"
)

func NewSqlStorage(cfg Config, runMigrations bool) (*SqlStorage, error) {
	err := cfg.SanityCheck()
	if err != nil {
		return nil, fmt.Errorf("config %s dont pass sanityCheck. Err: %w", cfg.String(), err)
	}
	log.Infof("Opening sync DB: cfg=%s", cfg.String())
	db, err := sql.Open(cfg.DriverName, cfg.DataSource)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	if runMigrations {
		log.Infof("Running DB migrations")
		err := RunMigrationsUp(cfg.DriverName, db)
		if err != nil {
			err := fmt.Errorf("error executing migrations: %w", err)
			log.Errorf(err.Error())
			return nil, err
		}
	}
	if cfg.DriverName == SqliteDriverName {
		log.Debugf("Enabling foreign keys for sqlite database")
		_, err = db.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			err := fmt.Errorf("error enabling foreign keys: %w", err)
			log.Errorf(err.Error())
			return nil, err
		}
	}
	return &SqlStorage{db}, nil
}

// BeginTransaction starts a transaction
func (s *SqlStorage) BeginTransaction(ctx context.Context) (dbTxType, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	stateTx := NewTxImpl(tx)
	return stateTx, nil
}

// BuildTableName returns the table name with the database prefixed if apply
// example: sqlite: BuildTableName("table") -> "table"
// example: postgres: BuildTableName("table") -> "sync.table"
func (s *SqlStorage) BuildTableName(tableName string) string {
	return tableName
}

func (s *SqlStorage) getExecQuerier(dbTx *sql.Tx) execQuerier {
	if dbTx != nil {
		return dbTx
	}
	return s.db
}
