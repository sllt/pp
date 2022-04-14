package pp

import (
	"manlu.org/pp/internal/builder"
	"testing"

	"github.com/stretchr/testify/suite"
	"manlu.org/pp/exp"
	"manlu.org/pp/gen/mocks"
)

type dialectTestSuite struct {
	suite.Suite
}

func (dts *dialectTestSuite) TestDialect() {
	opts := DefaultDialectOptions()
	sm := new(mocks.SelectSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, selectGen: sm}

	dts.Equal("test", d.Dialect())
}

func (dts *dialectTestSuite) TestToSelectSQL() {
	opts := DefaultDialectOptions()
	sm := new(mocks.SelectSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, selectGen: sm}

	b := builder.NewSQLBuilder(true)
	sc := exp.NewSelectClauses()
	sm.On("Generate", b, sc).Return(nil).Once()

	d.ToSelectSQL(b, sc)
	sm.AssertExpectations(dts.T())
}

func (dts *dialectTestSuite) TestToUpdateSQL() {
	opts := DefaultDialectOptions()
	um := new(mocks.UpdateSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, updateGen: um}

	b := builder.NewSQLBuilder(true)
	uc := exp.NewUpdateClauses()
	um.On("Generate", b, uc).Return(nil).Once()

	d.ToUpdateSQL(b, uc)
	um.AssertExpectations(dts.T())
}

func (dts *dialectTestSuite) TestToInsertSQL() {
	opts := DefaultDialectOptions()
	im := new(mocks.InsertSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, insertGen: im}

	b := builder.NewSQLBuilder(true)
	ic := exp.NewInsertClauses()
	im.On("Generate", b, ic).Return(nil).Once()

	d.ToInsertSQL(b, ic)
	im.AssertExpectations(dts.T())
}

func (dts *dialectTestSuite) TestToDeleteSQL() {
	opts := DefaultDialectOptions()
	dm := new(mocks.DeleteSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, deleteGen: dm}

	b := builder.NewSQLBuilder(true)
	dc := exp.NewDeleteClauses()
	dm.On("Generate", b, dc).Return(nil).Once()

	d.ToDeleteSQL(b, dc)
	dm.AssertExpectations(dts.T())
}

func (dts *dialectTestSuite) TestToTruncateSQL() {
	opts := DefaultDialectOptions()
	tm := new(mocks.TruncateSQLGenerator)
	d := sqlDialect{dialect: "test", dialectOptions: opts, truncateGen: tm}

	b := builder.NewSQLBuilder(true)
	tc := exp.NewTruncateClauses()
	tm.On("Generate", b, tc).Return(nil).Once()

	d.ToTruncateSQL(b, tc)
	tm.AssertExpectations(dts.T())
}

func TestSQLDialect(t *testing.T) {
	suite.Run(t, new(dialectTestSuite))
}
