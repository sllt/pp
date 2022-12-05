package gen

import (
	"strings"

	"github.com/sllt/pp/exp"
	"github.com/sllt/pp/internal/builder"
	"github.com/sllt/pp/internal/errors"
)

type (
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	InsertSQLGenerator interface {
		Dialect() string
		Generate(b builder.SQLBuilder, clauses exp.InsertClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (github.com/sllt/pp/dialect/postgres)
	insertSQLGenerator struct {
		CommonSQLGenerator
	}
)

var (
	ErrConflictUpdateValuesRequired = errors.New("values are required for on conflict update expression")
	ErrNoSourceForInsert            = errors.New("no source found when generating insert sql")
)

func errMisMatchedRowLength(expectedL, actualL int) error {
	return errors.New("rows with different value length expected %d got %d", expectedL, actualL)
}

func errUpsertWithWhereNotSupported(dialect string) error {
	return errors.New("dialect does not support upsert with where clause [dialect=%s]", dialect)
}

func NewInsertSQLGenerator(dialect string, do *SQLDialectOptions) InsertSQLGenerator {
	return &insertSQLGenerator{NewCommonSQLGenerator(dialect, do)}
}

func (isg *insertSQLGenerator) Generate(
	b builder.SQLBuilder,
	clauses exp.InsertClauses,
) {
	if !clauses.HasInto() {
		b.SetError(ErrNoSourceForInsert)
		return
	}
	for _, f := range isg.DialectOptions().InsertSQLOrder {
		if b.Error() != nil {
			return
		}
		switch f {
		case CommonTableSQLFragment:
			isg.ExpressionSQLGenerator().Generate(b, clauses.CommonTables())
		case InsertBeingSQLFragment:
			isg.InsertBeginSQL(b, clauses.OnConflict())
		case IntoSQLFragment:
			b.WriteRunes(isg.DialectOptions().SpaceRune)
			isg.ExpressionSQLGenerator().Generate(b, clauses.Into())
		case InsertSQLFragment:
			isg.InsertSQL(b, clauses)
		case ReturningSQLFragment:
			isg.ReturningSQL(b, clauses.Returning())
		default:
			b.SetError(ErrNotSupportedFragment("INSERT", f))
		}
	}
}

// Adds the correct fragment to being an INSERT statement
func (isg *insertSQLGenerator) InsertBeginSQL(b builder.SQLBuilder, o exp.ConflictExpression) {
	if isg.DialectOptions().SupportsInsertIgnoreSyntax && o != nil {
		b.Write(isg.DialectOptions().InsertIgnoreClause)
	} else {
		b.Write(isg.DialectOptions().InsertClause)
	}
}

// Adds the columns list to an insert statement
func (isg *insertSQLGenerator) InsertSQL(b builder.SQLBuilder, ic exp.InsertClauses) {
	switch {
	case ic.HasRows():
		ie, err := exp.NewInsertExpression(ic.Rows()...)
		if err != nil {
			b.SetError(err)
			return
		}
		isg.InsertExpressionSQL(b, ie)
	case ic.HasCols() && ic.HasVals():
		isg.insertColumnsSQL(b, ic.Cols())
		isg.insertValuesSQL(b, ic.Vals())
	case ic.HasCols() && ic.HasFrom():
		isg.insertColumnsSQL(b, ic.Cols())
		isg.insertFromSQL(b, ic.From())
	case ic.HasFrom():
		isg.insertFromSQL(b, ic.From())
	default:
		isg.defaultValuesSQL(b)
	}
	if ic.HasAlias() {
		b.Write(isg.DialectOptions().AsFragment)
		isg.ExpressionSQLGenerator().Generate(b, ic.Alias())
	}
	isg.onConflictSQL(b, ic.OnConflict())
}

func (isg *insertSQLGenerator) InsertExpressionSQL(b builder.SQLBuilder, ie exp.InsertExpression) {
	switch {
	case ie.IsInsertFrom():
		isg.insertFromSQL(b, ie.From())
	case ie.IsEmpty():
		isg.defaultValuesSQL(b)
	default:
		isg.insertColumnsSQL(b, ie.Cols())
		isg.insertValuesSQL(b, ie.Vals())
	}
}

// Adds the DefaultValuesFragment to an SQL statement
func (isg *insertSQLGenerator) defaultValuesSQL(b builder.SQLBuilder) {
	b.Write(isg.DialectOptions().DefaultValuesFragment)
}

func (isg *insertSQLGenerator) insertFromSQL(b builder.SQLBuilder, ae exp.AppendableExpression) {
	b.WriteRunes(isg.DialectOptions().SpaceRune)
	ae.AppendSQL(b)
}

// Adds the columns list to an insert statement
func (isg *insertSQLGenerator) insertColumnsSQL(b builder.SQLBuilder, cols exp.ColumnListExpression) {
	b.WriteRunes(isg.DialectOptions().SpaceRune, isg.DialectOptions().LeftParenRune)
	isg.ExpressionSQLGenerator().Generate(b, cols)
	b.WriteRunes(isg.DialectOptions().RightParenRune)
}

// Adds the values clause to an SQL statement
func (isg *insertSQLGenerator) insertValuesSQL(b builder.SQLBuilder, values [][]interface{}) {
	b.Write(isg.DialectOptions().ValuesFragment)
	rowLen := len(values[0])
	valueLen := len(values)
	for i, row := range values {
		if len(row) != rowLen {
			b.SetError(errMisMatchedRowLength(rowLen, len(row)))
			return
		}
		isg.ExpressionSQLGenerator().Generate(b, row)
		if i < valueLen-1 {
			b.WriteRunes(isg.DialectOptions().CommaRune, isg.DialectOptions().SpaceRune)
		}
	}
}

// Adds the DefaultValuesFragment to an SQL statement
func (isg *insertSQLGenerator) onConflictSQL(b builder.SQLBuilder, o exp.ConflictExpression) {
	if o == nil {
		return
	}
	b.Write(isg.DialectOptions().ConflictFragment)
	switch t := o.(type) {
	case exp.ConflictUpdateExpression:
		target := t.TargetColumn()
		if isg.DialectOptions().SupportsConflictTarget && target != "" {
			wrapParens := !strings.HasPrefix(strings.ToLower(target), "on constraint")

			b.WriteRunes(isg.DialectOptions().SpaceRune)
			if wrapParens {
				b.WriteRunes(isg.DialectOptions().LeftParenRune).
					WriteStrings(target).
					WriteRunes(isg.DialectOptions().RightParenRune)
			} else {
				b.Write([]byte(target))
			}
		}
		isg.onConflictDoUpdateSQL(b, t)
	default:
		b.Write(isg.DialectOptions().ConflictDoNothingFragment)
	}
}

func (isg *insertSQLGenerator) onConflictDoUpdateSQL(b builder.SQLBuilder, o exp.ConflictUpdateExpression) {
	b.Write(isg.DialectOptions().ConflictDoUpdateFragment)
	update := o.Update()
	if update == nil {
		b.SetError(ErrConflictUpdateValuesRequired)
		return
	}
	ue, err := exp.NewUpdateExpressions(update)
	if err != nil {
		b.SetError(err)
		return
	}
	isg.UpdateExpressionSQL(b, ue...)
	if b.Error() == nil && o.WhereClause() != nil {
		if !isg.DialectOptions().SupportsConflictUpdateWhere {
			b.SetError(errUpsertWithWhereNotSupported(isg.Dialect()))
			return
		}
		isg.WhereSQL(b, o.WhereClause())
	}
}
