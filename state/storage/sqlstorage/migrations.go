package sqlstorage

import (
	"database/sql"
	"embed"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*
var dbMigrations embed.FS

func RunMigrationsUp(driverName string, db *sql.DB) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root:       "migrations",
	}
	direction := migrate.Up
	nMigrations, err := migrate.Exec(db, driverName, migrations, direction)
	if err != nil {
		return err
	}

	log.Info("successfully ran ", nMigrations, " migrations")
	return nil

}
