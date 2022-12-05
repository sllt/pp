package exec

import (
	"context"
	"database/sql"

	"github.com/sllt/pp/internal/builder"
)

type (
	// nolint:stylecheck // keep name for backwards compatibility
	DbExecutor interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	}
	QueryFactory interface {
		FromSQL(sql string, args ...interface{}) QueryExecutor
		FromSQLBuilder(b builder.SQLBuilder) QueryExecutor
	}
	querySupport struct {
		de DbExecutor
	}
)

func NewQueryFactory(de DbExecutor) QueryFactory {
	return &querySupport{de}
}

func (qs *querySupport) FromSQL(query string, args ...interface{}) QueryExecutor {
	return newQueryExecutor(qs.de, nil, query, args...)
}

func (qs *querySupport) FromSQLBuilder(b builder.SQLBuilder) QueryExecutor {
	query, args, err := b.Build()
	return newQueryExecutor(qs.de, err, query, args...)
}
