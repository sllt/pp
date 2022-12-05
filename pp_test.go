package pp_test

import (
	"github.com/sllt/pp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type (
	dialectWrapperSuite struct {
		suite.Suite
	}
)

func (dws *dialectWrapperSuite) SetupSuite() {
	testDialect := pp.DefaultDialectOptions()
	// override to some value to ensure correct dialect is set
	pp.RegisterDialect("test", testDialect)
}

func (dws *dialectWrapperSuite) TearDownSuite() {
	pp.DeregisterDialect("test")
}

func (dws *dialectWrapperSuite) TestFrom() {
	dw := pp.Dialect("test")
	dws.Equal(pp.From("table").WithDialect("test"), dw.From("table"))
}

func (dws *dialectWrapperSuite) TestSelect() {
	dw := pp.Dialect("test")
	dws.Equal(pp.Select("col").WithDialect("test"), dw.Select("col"))
}

func (dws *dialectWrapperSuite) TestInsert() {
	dw := pp.Dialect("test")
	dws.Equal(pp.Insert("table").WithDialect("test"), dw.Insert("table"))
}

func (dws *dialectWrapperSuite) TestDelete() {
	dw := pp.Dialect("test")
	dws.Equal(pp.Delete("table").WithDialect("test"), dw.Delete("table"))
}

func (dws *dialectWrapperSuite) TestTruncate() {
	dw := pp.Dialect("test")
	dws.Equal(pp.Truncate("table").WithDialect("test"), dw.Truncate("table"))
}

func (dws *dialectWrapperSuite) TestDB() {
	mDB, _, err := sqlmock.New()
	dws.Require().NoError(err)
	dw := pp.Dialect("test")
	dws.Equal(pp.New("test", mDB), dw.DB(mDB))
}

func TestDialectWrapper(t *testing.T) {
	suite.Run(t, new(dialectWrapperSuite))
}
