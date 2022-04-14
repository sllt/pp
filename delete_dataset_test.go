package pp_test

import (
	"manlu.org/pp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"manlu.org/pp/exp"
	"manlu.org/pp/internal/builder"
	"manlu.org/pp/internal/errors"
	"manlu.org/pp/mocks"
)

type (
	deleteTestCase struct {
		ds      *pp.DeleteDataset
		clauses exp.DeleteClauses
	}
	deleteDatasetSuite struct {
		suite.Suite
	}
)

func (dds *deleteDatasetSuite) assertCases(cases ...deleteTestCase) {
	for _, s := range cases {
		dds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (dds *deleteDatasetSuite) SetupSuite() {
	noReturn := pp.DefaultDialectOptions()
	noReturn.SupportsReturn = false
	pp.RegisterDialect("no-return", noReturn)

	limitOnDelete := pp.DefaultDialectOptions()
	limitOnDelete.SupportsLimitOnDelete = true
	pp.RegisterDialect("limit-on-delete", limitOnDelete)

	orderOnDelete := pp.DefaultDialectOptions()
	orderOnDelete.SupportsOrderByOnDelete = true
	pp.RegisterDialect("order-on-delete", orderOnDelete)
}

func (dds *deleteDatasetSuite) TearDownSuite() {
	pp.DeregisterDialect("no-return")
	pp.DeregisterDialect("limit-on-delete")
	pp.DeregisterDialect("order-on-delete")
}

func (dds *deleteDatasetSuite) TestDelete() {
	ds := pp.Delete("test")
	dds.IsType(&pp.DeleteDataset{}, ds)
	dds.Implements((*exp.Expression)(nil), ds)
	dds.Implements((*exp.AppendableExpression)(nil), ds)
}

func (dds *deleteDatasetSuite) TestClone() {
	ds := pp.Delete("test")
	dds.Equal(ds.Clone(), ds)
}

func (dds *deleteDatasetSuite) TestExpression() {
	ds := pp.Delete("test")
	dds.Equal(ds.Expression(), ds)
}

func (dds *deleteDatasetSuite) TestDialect() {
	ds := pp.Delete("test")
	dds.NotNil(ds.Dialect())
}

func (dds *deleteDatasetSuite) TestWithDialect() {
	ds := pp.Delete("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := pp.GetDialect("default")
	dialectDs := ds.WithDialect("default")
	dds.Equal(md, ds.Dialect())
	dds.Equal(dialect, dialectDs.Dialect())
}

func (dds *deleteDatasetSuite) TestPrepared() {
	ds := pp.Delete("test")
	preparedDs := ds.Prepared(true)
	dds.True(preparedDs.IsPrepared())
	dds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	dds.True(preparedDs.Where(pp.Ex{"a": 1}).IsPrepared())

	defer pp.SetDefaultPrepared(false)
	pp.SetDefaultPrepared(true)

	// should be prepared by default
	ds = pp.Delete("test")
	dds.True(ds.IsPrepared())
}

func (dds *deleteDatasetSuite) TestGetClauses() {
	ds := pp.Delete("test")
	ce := exp.NewDeleteClauses().SetFrom(pp.I("test"))
	dds.Equal(ce, ds.GetClauses())
}

func (dds *deleteDatasetSuite) TestWith() {
	from := pp.From("cte")
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.With("test-cte", from),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestWithRecursive() {
	from := pp.From("cte")
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.WithRecursive("test-cte", from),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(true, "test-cte", from)),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestFrom_withIdentifier() {
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds:      bd.From("items2"),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items2")),
		},
		deleteTestCase{
			ds:      bd.From(pp.C("items2")),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items2")),
		},
		deleteTestCase{
			ds:      bd.From(pp.T("items2")),
			clauses: exp.NewDeleteClauses().SetFrom(pp.T("items2")),
		},
		deleteTestCase{
			ds:      bd.From("schema.table"),
			clauses: exp.NewDeleteClauses().SetFrom(pp.I("schema.table")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)

	dds.PanicsWithValue(pp.ErrBadFromArgument, func() {
		pp.Delete("test").From(true)
	})
}

func (dds *deleteDatasetSuite) TestWhere() {
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Where(pp.Ex{"a": 1}),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}),
		},
		deleteTestCase{
			ds: bd.Where(pp.Ex{"a": 1}).Where(pp.C("b").Eq("c")),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}).
				WhereAppend(pp.C("b").Eq("c")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearWhere() {
	bd := pp.Delete("items").Where(pp.Ex{"a": 1})
	dds.assertCases(
		deleteTestCase{
			ds: bd.ClearWhere(),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrder() {
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Order(pp.C("a").Asc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc()),
		},
		deleteTestCase{
			ds: bd.Order(pp.C("a").Asc()).Order(pp.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("b").Desc()),
		},
		deleteTestCase{
			ds: bd.Order(pp.C("a").Asc(), pp.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc(), pp.C("b").Desc()),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrderAppend() {
	bd := pp.Delete("items").Order(pp.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds: bd.OrderAppend(pp.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc(), pp.C("b").Desc()),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestOrderPrepend() {
	bd := pp.Delete("items").Order(pp.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds: bd.OrderPrepend(pp.C("b").Desc()),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("b").Desc(), pp.C("a").Asc()),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearOrder() {
	bd := pp.Delete("items").Order(pp.C("a").Asc())
	dds.assertCases(
		deleteTestCase{
			ds:      bd.ClearOrder(),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
		deleteTestCase{
			ds: bd,
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetOrder(pp.C("a").Asc()),
		},
	)
}

func (dds *deleteDatasetSuite) TestLimit() {
	bd := pp.Delete("test")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Limit(10),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("test")).
				SetLimit(uint(10)),
		},
		deleteTestCase{
			ds:      bd.Limit(0),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")),
		},
		deleteTestCase{
			ds: bd.Limit(10).Limit(2),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("test")).
				SetLimit(uint(2)),
		},
		deleteTestCase{
			ds:      bd.Limit(10).Limit(0),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")),
		},
	)
}

func (dds *deleteDatasetSuite) TestLimitAll() {
	bd := pp.Delete("test")
	dds.assertCases(
		deleteTestCase{
			ds: bd.LimitAll(),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("test")).
				SetLimit(pp.L("ALL")),
		},
		deleteTestCase{
			ds: bd.Limit(10).LimitAll(),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("test")).
				SetLimit(pp.L("ALL")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")),
		},
	)
}

func (dds *deleteDatasetSuite) TestClearLimit() {
	bd := pp.Delete("test").Limit(10)
	dds.assertCases(
		deleteTestCase{
			ds:      bd.ClearLimit(),
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("test")).SetLimit(uint(10)),
		},
	)
}

func (dds *deleteDatasetSuite) TestReturning() {
	bd := pp.Delete("items")
	dds.assertCases(
		deleteTestCase{
			ds: bd.Returning("a"),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetReturning(exp.NewColumnListExpression("a")),
		},
		deleteTestCase{
			ds: bd.Returning(),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		deleteTestCase{
			ds: bd.Returning(nil),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		deleteTestCase{
			ds: bd.Returning("a").Returning("b"),
			clauses: exp.NewDeleteClauses().
				SetFrom(pp.C("items")).
				SetReturning(exp.NewColumnListExpression("b")),
		},
		deleteTestCase{
			ds:      bd,
			clauses: exp.NewDeleteClauses().SetFrom(pp.C("items")),
		},
	)
}

func (dds *deleteDatasetSuite) TestReturnsColumns() {
	ds := pp.Delete("test")
	dds.False(ds.ReturnsColumns())
	dds.True(ds.Returning("foo", "bar").ReturnsColumns())
}

func (dds *deleteDatasetSuite) TestBuild() {
	md := new(mocks.SQLDialect)
	ds := pp.Delete("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Nil(err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestBuild_Prepared() {
	md := new(mocks.SQLDialect)
	ds := pp.Delete("test").Prepared(true).SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(true)
	md.On("ToDeleteSQL", sqlB, c).Return(nil).Once()

	sql, args, err := ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Nil(err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestBuild_WithError() {
	md := new(mocks.SQLDialect)
	ds := pp.Delete("test").SetDialect(md)
	c := ds.GetClauses()
	ee := errors.New("expected error")
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(ee)
	}).Once()

	sql, args, err := ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(ee, err)
	md.AssertExpectations(dds.T())
}

func (dds *deleteDatasetSuite) TestExecutor() {
	mDB, _, err := sqlmock.New()
	dds.NoError(err)

	ds := pp.New("mock", mDB).Delete("items").Where(pp.Ex{"id": pp.Op{"gt": 10}})

	dsql, args, err := ds.Executor().Build()
	dds.NoError(err)
	dds.Empty(args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > 10)`, dsql)

	dsql, args, err = ds.Prepared(true).Executor().Build()
	dds.NoError(err)
	dds.Equal([]interface{}{int64(10)}, args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > ?)`, dsql)

	defer pp.SetDefaultPrepared(false)
	pp.SetDefaultPrepared(true)

	dsql, args, err = ds.Executor().Build()
	dds.NoError(err)
	dds.Equal([]interface{}{int64(10)}, args)
	dds.Equal(`DELETE FROM "items" WHERE ("id" > ?)`, dsql)
}

func (dds *deleteDatasetSuite) TestSetError() {
	err1 := errors.New("error #1")
	err2 := errors.New("error #2")
	err3 := errors.New("error #3")

	// Verify initial error set/get works properly
	md := new(mocks.SQLDialect)
	ds := pp.Delete("test").SetDialect(md)
	ds = ds.SetError(err1)
	dds.Equal(err1, ds.Error())
	sql, args, err := ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Repeated SetError calls on Dataset should not overwrite the original error
	ds = ds.SetError(err2)
	dds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Builder functions should not lose the error
	ds = ds.ClearLimit()
	dds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)

	// Deeper errors inside SQL generation should still return original error
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToDeleteSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(err3)
	}).Once()

	sql, args, err = ds.Build()
	dds.Empty(sql)
	dds.Empty(args)
	dds.Equal(err1, err)
}

func TestDeleteDataset(t *testing.T) {
	suite.Run(t, new(deleteDatasetSuite))
}
