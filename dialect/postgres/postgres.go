package postgres

import (
	"manlu.org/pp"
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
