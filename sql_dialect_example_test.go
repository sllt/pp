package pp

import (
	"fmt"
)

func ExampleRegisterDialect() {
	opts := DefaultDialectOptions()
	opts.QuoteRune = '`'
	RegisterDialect("custom-dialect", opts)

	dialect := Dialect("custom-dialect")

	ds := dialect.From("test")

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM `test` []
}
