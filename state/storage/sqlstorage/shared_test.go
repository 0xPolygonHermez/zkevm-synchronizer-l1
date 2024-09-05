package sqlstorage_test

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/storage/sqlstorage"
)

func skipDatabaseTestIfNeeded(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping database test in short mode")
	}
}
func getStorageConfig() sqlstorage.Config {
	return sqlstorage.Config{
		DriverName: "sqlite3",
		DataSource: "file::memory:",
	}
}
