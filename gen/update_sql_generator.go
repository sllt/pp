package gen

import (
	"github.com/sllt/pp/exp"
	"github.com/sllt/pp/internal/builder"
	"github.com/sllt/pp/internal/errors"
)

type (
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	UpdateSQLGenerator interface {
		Dialect() string
		Generate(b builder.SQLBuilder, clauses exp.UpdateClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (github.com/sllt/pp/dialect/postgres)
	updateSQLGenerator struct {
		CommonSQLGenerator
	}
)

var (
	ErrNoSourceForUpdate    = errors.New("no source found when generating update sql")
	ErrNoSetValuesForUpdate = errors.New("no set values found when generating UPDATE sql")
)

func NewUpdateSQLGenerator(dialect string, do *SQLDialectOptions) UpdateSQLGenerator {
	return &updateSQLGenerator{NewCommonSQLGenerator(dialect, do)}
}

func (usg *updateSQLGenerator) Generate(b builder.SQLBuilder, clauses exp.UpdateClauses) {
	if !clauses.HasTable() {
		b.SetError(ErrNoSourceForUpdate)
		return
	}
	if !clauses.HasSetValues() {
		b.SetError(ErrNoSetValuesForUpdate)
		return
	}
	if !usg.DialectOptions().SupportsMultipleUpdateTables && clauses.HasFrom() {
		b.SetError(errors.New("%s dialect does not support multiple tables in UPDATE", usg.Dialect()))
	}
	updates, err := exp.NewUpdateExpressions(clauses.SetValues())
	if err != nil {
		b.SetError(err)
		return
	}
	for _, f := range usg.DialectOptions().UpdateSQLOrder {
		if b.Error() != nil {
			return
		}
		switch f {
		case CommonTableSQLFragment:
			usg.ExpressionSQLGenerator().Generate(b, clauses.CommonTables())
		case UpdateBeginSQLFragment:
			usg.UpdateBeginSQL(b)
		case SourcesSQLFragment:
			usg.updateTableSQL(b, clauses)
		case UpdateSQLFragment:
			usg.UpdateExpressionsSQL(b, updates...)
		case UpdateFromSQLFragment:
			usg.updateFromSQL(b, clauses.From())
		case WhereSQLFragment:
			usg.WhereSQL(b, clauses.Where())
		case OrderSQLFragment:
			if usg.DialectOptions().SupportsOrderByOnUpdate {
				usg.OrderSQL(b, clauses.Order())
			}
		case LimitSQLFragment:
			if usg.DialectOptions().SupportsLimitOnUpdate {
				usg.LimitSQL(b, clauses.Limit())
			}
		case ReturningSQLFragment:
			usg.ReturningSQL(b, clauses.Returning())
		default:
			b.SetError(ErrNotSupportedFragment("UPDATE", f))
		}
	}
}

// Adds the correct fragment to being an UPDATE statement
func (usg *updateSQLGenerator) UpdateBeginSQL(b builder.SQLBuilder) {
	b.Write(usg.DialectOptions().UpdateClause)
}

// Adds column setters in an update SET clause
func (usg *updateSQLGenerator) UpdateExpressionsSQL(b builder.SQLBuilder, updates ...exp.UpdateExpression) {
	b.Write(usg.DialectOptions().SetFragment)
	usg.UpdateExpressionSQL(b, updates...)
}

func (usg *updateSQLGenerator) updateTableSQL(b builder.SQLBuilder, uc exp.UpdateClauses) {
	b.WriteRunes(usg.DialectOptions().SpaceRune)
	usg.ExpressionSQLGenerator().Generate(b, uc.Table())
	if uc.HasFrom() {
		if !usg.DialectOptions().UseFromClauseForMultipleUpdateTables {
			b.WriteRunes(usg.DialectOptions().CommaRune)
			usg.ExpressionSQLGenerator().Generate(b, uc.From())
		}
	}
}

func (usg *updateSQLGenerator) updateFromSQL(b builder.SQLBuilder, ce exp.ColumnListExpression) {
	if ce == nil || ce.IsEmpty() {
		return
	}
	if usg.DialectOptions().UseFromClauseForMultipleUpdateTables {
		usg.FromSQL(b, ce)
	}
}
