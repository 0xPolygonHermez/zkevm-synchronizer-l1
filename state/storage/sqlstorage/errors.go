package sqlstorage

import (
	"database/sql"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
)

func translateSqlError(err error, contextDescription string) error {
	if err == nil {
		return nil
	}
	newErr := err
	if err == sql.ErrNoRows {
		newErr = entities.ErrNotFound
	}

	return fmt.Errorf("storage error: %s: Err: %w", contextDescription, newErr)
}
