package pp_test

import (
	"fmt"
	"github.com/sllt/pp"
)

func ExampleRegisterDialect() {
	opts := pp.DefaultDialectOptions()
	opts.QuoteRune = '`'
	pp.RegisterDialect("custom-dialect", opts)

	dialect := pp.Dialect("custom-dialect")

	ds := dialect.From("test")

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM `test` []
}
