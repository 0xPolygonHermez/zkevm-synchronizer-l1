package storage

import (
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/db/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/utils/gerror"
)

// Storage interface
type Storage interface{}

// NewStorage creates a new Storage
func NewStorage(cfg Config) (Storage, error) {
	if cfg.Database == "postgres" {
		pg, err := pgstorage.NewPostgresStorage(pgstorage.Config{
			Name:     cfg.Name,
			User:     cfg.User,
			Password: cfg.Password,
			Host:     cfg.Host,
			Port:     cfg.Port,
			MaxConns: cfg.MaxConns,
		})
		return pg, err
	}
	return nil, gerror.ErrStorageNotRegister
}

// RunMigrations will execute pending migrations if needed to keep
// the database updated with the latest changes
func RunMigrations(cfg Config) error {
	config := pgstorage.Config{
		Name:     cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
		Host:     cfg.Host,
		Port:     cfg.Port,
	}
	return pgstorage.RunMigrationsUp(config)
}
