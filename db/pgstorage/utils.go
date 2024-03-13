package pgstorage

import (
	"embed"
	"os"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

// RunMigrationsUp migrate up.
func RunMigrationsUp(cfg Config) error {
	return runMigrations(cfg, migrate.Up)
}

// RunMigrationsDown migrate down.
func RunMigrationsDown(cfg Config) error {
	return runMigrations(cfg, migrate.Down)
}

// ResetDB.
func ResetDB(cfg Config) error {
	c, err := pgx.ParseConfig("postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.Name)
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)
	_, err = db.Exec("DROP SCHEMA IF EXISTS mt CASCADE;")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP SCHEMA IF EXISTS sync CASCADE;")
	if err != nil {
		return err
	}
	_, err = db.Exec("DROP TABLE IF EXISTS public.gorp_migrations;")
	if err != nil {
		return err
	}
	return nil
}

//go:embed migrations/*
var dbMigrations embed.FS

// runMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes in either direction of up or down.
func runMigrations(cfg Config, direction migrate.MigrationDirection) error {
	c, err := pgx.ParseConfig("postgres://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + ":" + cfg.Port + "/" + cfg.Name)
	if err != nil {
		return err
	}
	db := stdlib.OpenDB(*c)
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "migrations",
	}

	nMigrations, err := migrate.Exec(db, "postgres", migrations, direction)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations")
	return nil
}

// InitOrReset will initializes the db running the migrations or
// will reset all the known data and rerun the migrations
func InitOrReset(cfg Config) error {
	// connect to database
	_, err := NewPostgresStorage(cfg)
	if err != nil {
		return err
	}

	// run migrations
	if err := ResetDB(cfg); err != nil {
		return err
	}
	return RunMigrationsUp(cfg)
}

// NewConfigFromEnv creates config from standard postgres environment variables,
func NewConfigFromEnv() Config {
	maxConns, _ := strconv.Atoi(getEnv("ZKEVM_BRIDGE_DATABASE_MAXCONNS", "500"))
	return Config{
		User:     getEnv("ZKEVM_BRIDGE_DATABASE_USER", "test_user"),
		Password: getEnv("ZKEVM_BRIDGE_DATABASE_PASSWORD", "test_password"),
		Name:     getEnv("ZKEVM_BRIDGE_DATABASE_NAME", "test_db"),
		Host:     getEnv("ZKEVM_BRIDGE_DATABASE_HOST", "localhost"),
		Port:     getEnv("ZKEVM_BRIDGE_DATABASE_PORT", "5435"),
		MaxConns: maxConns,
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

// nolint
var (
	String, _ = abi.NewType("string", "", nil)
	Uint8, _  = abi.NewType("uint8", "", nil)
)
