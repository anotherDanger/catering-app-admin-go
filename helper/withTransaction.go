package helper

import (
	"database/sql"
)

func WithTransaction(tx *sql.Tx, err *error) {
	if *err != nil {
		_ = tx.Rollback()
	} else {
		_ = tx.Commit()
	}
}
