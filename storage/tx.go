package storage

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func completeTx(tx *sqlx.Tx, err error) error {
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			if rbErr == sql.ErrTxDone {
				return errors.Wrap(err, "detected error but transaction is already being complete")
			}
			return errors.Wrap(rbErr, "detected error but rollback failed due to error")
		}
		return errors.Wrap(err, "detected error and rolled transaction back")
	}

	if err := tx.Commit(); err != nil {
		if err == sql.ErrTxDone {
			return nil
		}
		return errors.Wrap(err, "commit failed due to error")
	}

	return nil
}
