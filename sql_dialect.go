package pp

import (
	"manlu.org/pp/exp"
	"manlu.org/pp/gen"
	"manlu.org/pp/internal/builder"
	"strings"
	"sync"
)

type (
	SQLDialectOptions = gen.SQLDialectOptions
	// An adapter interface to be used by a Dataset to generate SQL for a specific dialect.
	// See DefaultAdapter for a concrete implementation and examples.
	SQLDialect interface {
		Dialect() string
		ToSelectSQL(b builder.SQLBuilder, clauses exp.SelectClauses)
		ToUpdateSQL(b builder.SQLBuilder, clauses exp.UpdateClauses)
		ToInsertSQL(b builder.SQLBuilder, clauses exp.InsertClauses)
		ToDeleteSQL(b builder.SQLBuilder, clauses exp.DeleteClauses)
		ToTruncateSQL(b builder.SQLBuilder, clauses exp.TruncateClauses)
	}
	// The default adapter. This class should be used when building a new adapter. When creating a new adapter you can
	// either override methods, or more typically update default values.
	// See (github.com/doug-martin/pp/dialect/postgres)
	sqlDialect struct {
		dialect        string
		dialectOptions *SQLDialectOptions
		selectGen      gen.SelectSQLGenerator
		updateGen      gen.UpdateSQLGenerator
		insertGen      gen.InsertSQLGenerator
		deleteGen      gen.DeleteSQLGenerator
		truncateGen    gen.TruncateSQLGenerator
	}
)

var (
	dialects              = make(map[string]SQLDialect)
	DefaultDialectOptions = gen.DefaultDialectOptions
	dialectsMu            sync.RWMutex
)

func init() {
	RegisterDialect("default", DefaultDialectOptions())
}

func RegisterDialect(name string, do *SQLDialectOptions) {
	dialectsMu.Lock()
	defer dialectsMu.Unlock()
	lowerName := strings.ToLower(name)
	dialects[lowerName] = newDialect(lowerName, do)
}

func DeregisterDialect(name string) {
	dialectsMu.Lock()
	defer dialectsMu.Unlock()
	delete(dialects, strings.ToLower(name))
}

func GetDialect(name string) SQLDialect {
	name = strings.ToLower(name)
	if d, ok := dialects[name]; ok {
		return d
	}
	return newDialect("default", DefaultDialectOptions())
}

func newDialect(dialect string, do *SQLDialectOptions) SQLDialect {
	return &sqlDialect{
		dialect:        dialect,
		dialectOptions: do,
		selectGen:      gen.NewSelectSQLGenerator(dialect, do),
		updateGen:      gen.NewUpdateSQLGenerator(dialect, do),
		insertGen:      gen.NewInsertSQLGenerator(dialect, do),
		deleteGen:      gen.NewDeleteSQLGenerator(dialect, do),
		truncateGen:    gen.NewTruncateSQLGenerator(dialect, do),
	}
}

func (d *sqlDialect) Dialect() string {
	return d.dialect
}

func (d *sqlDialect) ToSelectSQL(b builder.SQLBuilder, clauses exp.SelectClauses) {
	d.selectGen.Generate(b, clauses)
}

func (d *sqlDialect) ToUpdateSQL(b builder.SQLBuilder, clauses exp.UpdateClauses) {
	d.updateGen.Generate(b, clauses)
}

func (d *sqlDialect) ToInsertSQL(b builder.SQLBuilder, clauses exp.InsertClauses) {
	d.insertGen.Generate(b, clauses)
}

func (d *sqlDialect) ToDeleteSQL(b builder.SQLBuilder, clauses exp.DeleteClauses) {
	d.deleteGen.Generate(b, clauses)
}

func (d *sqlDialect) ToTruncateSQL(b builder.SQLBuilder, clauses exp.TruncateClauses) {
	d.truncateGen.Generate(b, clauses)
}