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
	updateTestCase struct {
		ds      *pp.UpdateDataset
		clauses exp.UpdateClauses
	}
	updateDatasetSuite struct {
		suite.Suite
	}
)

func (uds *updateDatasetSuite) assertCases(cases ...updateTestCase) {
	for _, s := range cases {
		uds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (uds *updateDatasetSuite) TestUpdate() {
	ds := pp.Update("test")
	uds.IsType(&pp.UpdateDataset{}, ds)
	uds.Implements((*exp.Expression)(nil), ds)
	uds.Implements((*exp.AppendableExpression)(nil), ds)
}

func (uds *updateDatasetSuite) TestClone() {
	ds := pp.Update("test")
	uds.Equal(ds, ds.Clone())
}

func (uds *updateDatasetSuite) TestExpression() {
	ds := pp.Update("test")
	uds.Equal(ds, ds.Expression())
}

func (uds *updateDatasetSuite) TestDialect() {
	ds := pp.Update("test")
	uds.NotNil(ds.Dialect())
}

func (uds *updateDatasetSuite) TestWithDialect() {
	ds := pp.Update("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := pp.GetDialect("default")
	dialectDs := ds.WithDialect("default")
	uds.Equal(md, ds.Dialect())
	uds.Equal(dialect, dialectDs.Dialect())
}

func (uds *updateDatasetSuite) TestPrepared() {
	ds := pp.Update("test")
	preparedDs := ds.Prepared(true)
	uds.True(preparedDs.IsPrepared())
	uds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	uds.True(preparedDs.Where(pp.Ex{"a": 1}).IsPrepared())

	defer pp.SetDefaultPrepared(false)
	pp.SetDefaultPrepared(true)

	// should be prepared by default
	ds = pp.Update("test")
	uds.True(ds.IsPrepared())
}

func (uds *updateDatasetSuite) TestGetClauses() {
	ds := pp.Update("test")
	ce := exp.NewUpdateClauses().SetTable(pp.I("test"))
	uds.Equal(ce, ds.GetClauses())
}

func (uds *updateDatasetSuite) TestWith() {
	from := pp.Update("cte")
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.With("test-cte", from),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestWithRecursive() {
	from := pp.Update("cte")
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.WithRecursive("test-cte", from),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				CommonTablesAppend(exp.NewCommonTableExpression(true, "test-cte", from)),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestTable() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds:      bd.Table("items2"),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items2")),
		},
		updateTestCase{
			ds:      bd.Table(pp.L("literal_table")),
			clauses: exp.NewUpdateClauses().SetTable(pp.L("literal_table")),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
	uds.PanicsWithValue(pp.ErrUnsupportedUpdateTableType, func() {
		bd.Table(true)
	})
}

func (uds *updateDatasetSuite) TestSet() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.Set(item{Name: "Test", Address: "111 Test Addr"}),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetSetValues(item{Name: "Test", Address: "111 Test Addr"}),
		},
		updateTestCase{
			ds: bd.Set(pp.Record{"name": "Test", "address": "111 Test Addr"}),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetSetValues(pp.Record{"name": "Test", "address": "111 Test Addr"}),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestFrom() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.From("other"),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetFrom(exp.NewColumnListExpression("other")),
		},
		updateTestCase{
			ds: bd.From("other").From("other2"),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetFrom(exp.NewColumnListExpression("other2")),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestWhere() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.Where(pp.Ex{"a": 1}),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}),
		},
		updateTestCase{
			ds: bd.Where(pp.Ex{"a": 1}).Where(pp.C("b").Eq("c")),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}).WhereAppend(pp.C("b").Eq("c")),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestClearWhere() {
	bd := pp.Update("items").Where(pp.Ex{"a": 1})
	uds.assertCases(
		updateTestCase{
			ds:      bd.ClearWhere(),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				WhereAppend(pp.Ex{"a": 1}),
		},
	)
}

func (uds *updateDatasetSuite) TestOrder() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.Order(pp.C("a").Desc()),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).OrderAppend(pp.C("a").Desc()),
		},
		updateTestCase{
			ds: bd.Order(pp.C("a").Desc()).Order(pp.C("b").Asc()),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("b").Asc()),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestOrderAppend() {
	bd := pp.Update("items").Order(pp.C("a").Desc())
	uds.assertCases(
		updateTestCase{
			ds: bd.OrderAppend(pp.C("b").Asc()),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("a").Desc()).
				OrderAppend(pp.C("b").Asc()),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("a").Desc()),
		},
	)
}

func (uds *updateDatasetSuite) TestOrderPrepend() {
	bd := pp.Update("items").Order(pp.C("a").Desc())
	uds.assertCases(
		updateTestCase{
			ds: bd.OrderPrepend(pp.C("b").Asc()),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("b").Asc()).
				OrderAppend(pp.C("a").Desc()),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("a").Desc()),
		},
	)
}

func (uds *updateDatasetSuite) TestClearOrder() {
	bd := pp.Update("items").Order(pp.C("a").Desc())
	uds.assertCases(
		updateTestCase{
			ds:      bd.ClearOrder(),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
		updateTestCase{
			ds: bd,
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				OrderAppend(pp.C("a").Desc()),
		},
	)
}

func (uds *updateDatasetSuite) TestLimit() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds:      bd.Limit(10),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")).SetLimit(uint(10)),
		},
		updateTestCase{
			ds:      bd.Limit(0),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestLimitAll() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds:      bd.LimitAll(),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")).SetLimit(pp.L("ALL")),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestClearLimit() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds:      bd.LimitAll().ClearLimit(),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
		updateTestCase{
			ds:      bd.Limit(10).ClearLimit(),
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestReturning() {
	bd := pp.Update("items")
	uds.assertCases(
		updateTestCase{
			ds: bd.Returning("a", "b"),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetReturning(exp.NewColumnListExpression("a", "b")),
		},
		updateTestCase{
			ds: bd.Returning(),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		updateTestCase{
			ds: bd.Returning(nil),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetReturning(exp.NewColumnListExpression()),
		},
		updateTestCase{
			ds: bd.Returning("a", "b").Returning("c"),
			clauses: exp.NewUpdateClauses().
				SetTable(pp.C("items")).
				SetReturning(exp.NewColumnListExpression("c")),
		},
		updateTestCase{
			ds:      bd,
			clauses: exp.NewUpdateClauses().SetTable(pp.C("items")),
		},
	)
}

func (uds *updateDatasetSuite) TestReturnsColumns() {
	ds := pp.Update("test")
	uds.False(ds.ReturnsColumns())
	uds.True(ds.Returning("foo", "bar").ReturnsColumns())
}

func (uds *updateDatasetSuite) TestBuild() {
	md := new(mocks.SQLDialect)
	ds := pp.Update("test").SetDialect(md)
	r := pp.Record{"c": "a"}
	c := ds.GetClauses().SetSetValues(r)
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToUpdateSQL", sqlB, c).Return(nil).Once()
	updateSQL, args, err := ds.Set(r).Build()
	uds.Empty(updateSQL)
	uds.Empty(args)
	uds.Nil(err)
	md.AssertExpectations(uds.T())
}

func (uds *updateDatasetSuite) TestBuild_Prepared() {
	md := new(mocks.SQLDialect)
	ds := pp.Update("test").Prepared(true).SetDialect(md)
	r := pp.Record{"c": "a"}
	c := ds.GetClauses().SetSetValues(r)
	sqlB := builder.NewSQLBuilder(true)
	md.On("ToUpdateSQL", sqlB, c).Return(nil).Once()
	updateSQL, args, err := ds.Set(pp.Record{"c": "a"}).Build()
	uds.Empty(updateSQL)
	uds.Empty(args)
	uds.Nil(err)
	md.AssertExpectations(uds.T())
}

func (uds *updateDatasetSuite) TestBuild_WithError() {
	md := new(mocks.SQLDialect)
	ds := pp.Update("test").SetDialect(md)
	r := pp.Record{"c": "a"}
	c := ds.GetClauses().SetSetValues(r)
	sqlB := builder.NewSQLBuilder(false)
	ee := errors.New("expected error")
	md.On("ToUpdateSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(ee)
	}).Once()

	updateSQL, args, err := ds.Set(pp.Record{"c": "a"}).Build()
	uds.Empty(updateSQL)
	uds.Empty(args)
	uds.Equal(ee, err)
	md.AssertExpectations(uds.T())
}

func (uds *updateDatasetSuite) TestExecutor() {
	mDB, _, err := sqlmock.New()
	uds.NoError(err)
	ds := pp.New("mock", mDB).
		Update("items").
		Set(pp.Record{"address": "111 Test Addr", "name": "Test1"}).
		Where(pp.C("name").IsNull())

	updateSQL, args, err := ds.Executor().Build()
	uds.NoError(err)
	uds.Empty(args)
	uds.Equal(`UPDATE "items" SET "address"='111 Test Addr',"name"='Test1' WHERE ("name" IS NULL)`, updateSQL)

	updateSQL, args, err = ds.Prepared(true).Executor().Build()
	uds.NoError(err)
	uds.Equal([]interface{}{"111 Test Addr", "Test1"}, args)
	uds.Equal(`UPDATE "items" SET "address"=?,"name"=? WHERE ("name" IS NULL)`, updateSQL)

	defer pp.SetDefaultPrepared(false)
	pp.SetDefaultPrepared(true)

	updateSQL, args, err = ds.Executor().Build()
	uds.NoError(err)
	uds.Equal([]interface{}{"111 Test Addr", "Test1"}, args)
	uds.Equal(`UPDATE "items" SET "address"=?,"name"=? WHERE ("name" IS NULL)`, updateSQL)
}

func (uds *updateDatasetSuite) TestSetError() {
	err1 := errors.New("error #1")
	err2 := errors.New("error #2")
	err3 := errors.New("error #3")

	// Verify initial error set/get works properly
	md := new(mocks.SQLDialect)
	ds := pp.Update("test").SetDialect(md)
	ds = ds.SetError(err1)
	uds.Equal(err1, ds.Error())
	sql, args, err := ds.Build()
	uds.Empty(sql)
	uds.Empty(args)
	uds.Equal(err1, err)

	// Repeated SetError calls on Dataset should not overwrite the original error
	ds = ds.SetError(err2)
	uds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	uds.Empty(sql)
	uds.Empty(args)
	uds.Equal(err1, err)

	// Builder functions should not lose the error
	ds = ds.ClearLimit()
	uds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	uds.Empty(sql)
	uds.Empty(args)
	uds.Equal(err1, err)

	// Deeper errors inside SQL generation should still return original error
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToUpdateSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(err3)
	}).Once()

	sql, args, err = ds.Build()
	uds.Empty(sql)
	uds.Empty(args)
	uds.Equal(err1, err)
}

func TestUpdateDataset(t *testing.T) {
	suite.Run(t, new(updateDatasetSuite))
}
