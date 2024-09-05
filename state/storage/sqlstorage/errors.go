package sqlstorage

import (
	"database/sql"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-synchronizer-l1/state/entities"
	gosqlite3 "github.com/mattn/go-sqlite3"
)

func translateSqlError(err error, contextDescription string) error {
	if err == nil {
		return nil
	}
	newErr := err
	if err == sql.ErrNoRows {
		newErr = entities.ErrNotFound
	}

	slErr, ok := err.(gosqlite3.Error)
	if ok {
		switch slErr.Code {
		case 19: // UNIQUE constraint failed
			if slErr.ExtendedCode == 1555 {
				newErr = fmt.Errorf("%w : sqlError:%w ", entities.ErrAlreadyExists, err)
			}
			if slErr.ExtendedCode == 787 {
				newErr = fmt.Errorf("%w : sqlError:%w ", entities.ErrForeignKeyViolation, err)
			}
		}
	}

	return fmt.Errorf("storage error: %s: Err: %w", contextDescription, newErr)
}
