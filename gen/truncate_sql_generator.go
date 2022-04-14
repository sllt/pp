package gen

import (
	"strings"

	"manlu.org/pp/exp"
	"manlu.org/pp/internal/builder"
	"manlu.org/pp/internal/errors"
)

type (
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	TruncateSQLGenerator interface {
		Dialect() string
		Generate(b builder.SQLBuilder, clauses exp.TruncateClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (manlu.org/pp/dialect/postgres)
	truncateSQLGenerator struct {
		CommonSQLGenerator
	}
)

var errNoSourceForTruncate = errors.New("no source found when generating truncate sql")

func NewTruncateSQLGenerator(dialect string, do *SQLDialectOptions) TruncateSQLGenerator {
	return &truncateSQLGenerator{NewCommonSQLGenerator(dialect, do)}
}

func (tsg *truncateSQLGenerator) Generate(b builder.SQLBuilder, clauses exp.TruncateClauses) {
	if !clauses.HasTable() {
		b.SetError(errNoSourceForTruncate)
		return
	}
	for _, f := range tsg.DialectOptions().TruncateSQLOrder {
		if b.Error() != nil {
			return
		}
		switch f {
		case TruncateSQLFragment:
			tsg.TruncateSQL(b, clauses.Table(), clauses.Options())
		default:
			b.SetError(ErrNotSupportedFragment("TRUNCATE", f))
		}
	}
}

// Generates a TRUNCATE statement
func (tsg *truncateSQLGenerator) TruncateSQL(b builder.SQLBuilder, from exp.ColumnListExpression, opts exp.TruncateOptions) {
	b.Write(tsg.DialectOptions().TruncateClause)
	tsg.SourcesSQL(b, from)
	if opts.Identity != tsg.DialectOptions().EmptyString {
		b.WriteRunes(tsg.DialectOptions().SpaceRune).
			WriteStrings(strings.ToUpper(opts.Identity)).
			Write(tsg.DialectOptions().IdentityFragment)
	}
	if opts.Cascade {
		b.Write(tsg.DialectOptions().CascadeFragment)
	} else if opts.Restrict {
		b.Write(tsg.DialectOptions().RestrictFragment)
	}
}
