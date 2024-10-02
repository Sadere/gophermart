package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func WrapTx(ctx context.Context, db *sqlx.DB, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(ctx, tx); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
