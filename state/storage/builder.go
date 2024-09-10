package storage

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/log"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/pgstorage"
	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
)

func NewStorage(config Config) (Storer, error) {
	if config.DriverName == "sqlite3" {
		return sqlstorage.NewSqlStorage(sqlstorage.Config{
			DriverName: config.DriverName,
			DataSource: config.DataSource,
		}, true)
	}
	if config.DriverName == "postgres" {
		log.Warnf("Deprecated driver %s, please use sqlite3", config.DriverName)
		parsedURL, err := url.Parse(config.DataSource)
		if err != nil {
			return nil, fmt.Errorf("error parsing datasource %s: %w", config.DataSource, err)
		}
		password, _ := parsedURL.User.Password()
		maxConns, err := strconv.Atoi(parsedURL.Query().Get("pool_max_conns"))
		if err != nil {
			return nil, fmt.Errorf("error getting pool_max_conns %s: %w", config.DataSource, err)
		}
		pgCfg := pgstorage.Config{
			User:     parsedURL.User.Username(),
			Password: password,
			Host:     parsedURL.Hostname(),
			Port:     parsedURL.Port(),
			Name:     strings.Trim(parsedURL.Path, "/"),
			MaxConns: maxConns,
		}
		return pgstorage.NewPostgresStorage(pgCfg)
	}
	return nil, fmt.Errorf("unknown driver %s", config.DriverName)
}
