package postgres

import (
	"github.com/sllt/pp"
)

func DialectOptions() *pp.SQLDialectOptions {
	do := pp.DefaultDialectOptions()
	do.PlaceHolderFragment = []byte("$")
	do.IncludePlaceholderNum = true
	return do
}

func init() {
	pp.RegisterDialect("postgres", DialectOptions())
}
