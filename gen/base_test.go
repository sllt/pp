package gen

import (
	"github.com/stretchr/testify/suite"
	"manlu.org/pp/internal/builder"
)

type baseSQLGeneratorSuite struct {
	suite.Suite
}

func (bsgs *baseSQLGeneratorSuite) assertNotPreparedSQL(b builder.SQLBuilder, expectedSQL string) {
	actualSQL, actualArgs, err := b.Build()
	bsgs.NoError(err)
	bsgs.Equal(expectedSQL, actualSQL)
	bsgs.Empty(actualArgs)
}

func (bsgs *baseSQLGeneratorSuite) assertPreparedSQL(
	b builder.SQLBuilder,
	expectedSQL string,
	expectedArgs []interface{},
) {
	actualSQL, actualArgs, err := b.Build()
	bsgs.NoError(err)
	bsgs.Equal(expectedSQL, actualSQL)
	if len(actualArgs) == 0 {
		bsgs.Empty(expectedArgs)
	} else {
		bsgs.Equal(expectedArgs, actualArgs)
	}
}

func (bsgs *baseSQLGeneratorSuite) assertErrorSQL(b builder.SQLBuilder, errMsg string) {
	actualSQL, actualArgs, err := b.Build()
	bsgs.EqualError(err, errMsg)
	bsgs.Empty(actualSQL)
	bsgs.Empty(actualArgs)
}
