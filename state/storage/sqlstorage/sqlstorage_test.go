package sqlstorage_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	cfg := sqlstorage.Config{
		DriverName: "sqlite3",
		DataSource: "file::memory:?cache=shared",
	}
	sut, err := sqlstorage.NewSqlStorage(cfg, true)
	require.NoError(t, err)

	tx, err := sut.BeginTransaction(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, tx)

}
