package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
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
