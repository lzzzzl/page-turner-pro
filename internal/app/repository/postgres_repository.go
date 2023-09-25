package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
	"github.com/lzzzzl/page-turner-pro/internal/domain/common"
)

type PostgresRepository struct {
	db   *sqlx.DB
	pgsq sq.StatementBuilderType
}

func NewPostgresRepository(ctx context.Context, db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
		// set the default placeholder as $ instead of ? because postgres uses $
		pgsq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (r *PostgresRepository) beginTx() (*sqlx.Tx, common.Error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, common.NewError(common.ErrorCodeRemoteProcess, err)
	}
	return tx, nil
}

func (r *PostgresRepository) finishTx(err common.Error, tx *sqlx.Tx) common.Error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			warpError := multierror.Append(err, rollbackErr)
			return common.NewError(common.ErrorCodeRemoteProcess, warpError)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return common.NewError(common.ErrorCodeRemoteProcess, commitErr)
		}

		return nil
	}
}
