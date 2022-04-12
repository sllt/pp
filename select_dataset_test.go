package pp

import (
	"manlu.org/pp/exp"
	"manlu.org/pp/internal/builder"
	"manlu.org/pp/internal/errors"
	"manlu.org/pp/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	selectTestCase struct {
		ds      *SelectDataset
		clauses exp.SelectClauses
	}
	dsTestActionItem struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	dsUntaggedTestActionItem struct {
		Address  string `db:"address"`
		Name     string `db:"name"`
		Untagged string
	}
	selectDatasetSuite struct {
		suite.Suite
	}
)

func (sds *selectDatasetSuite) assertCases(cases ...selectTestCase) {
	for _, s := range cases {
		sds.Equal(s.clauses, s.ds.GetClauses())
	}
}

func (sds *selectDatasetSuite) TestReturnsColumns() {
	ds := Select(L("NOW()"))
	sds.True(ds.ReturnsColumns())
}

func (sds *selectDatasetSuite) TestClone() {
	ds := From("test")
	sds.Equal(ds, ds.Clone())
}

func (sds *selectDatasetSuite) TestExpression() {
	ds := From("test")
	sds.Equal(ds, ds.Expression())
}

func (sds *selectDatasetSuite) TestDialect() {
	ds := From("test")
	sds.NotNil(ds.Dialect())
}

func (sds *selectDatasetSuite) TestWithDialect() {
	ds := From("test")
	md := new(mocks.SQLDialect)
	ds = ds.SetDialect(md)

	dialect := GetDialect("default")
	dialectDs := ds.WithDialect("default")
	sds.Equal(md, ds.Dialect())
	sds.Equal(dialect, dialectDs.Dialect())
}

func (sds *selectDatasetSuite) TestPrepared() {
	ds := From("test")
	preparedDs := ds.Prepared(true)
	sds.True(preparedDs.IsPrepared())
	sds.False(ds.IsPrepared())
	// should apply the prepared to any datasets created from the root
	sds.True(preparedDs.Where(Ex{"a": 1}).IsPrepared())

	defer SetDefaultPrepared(false)
	SetDefaultPrepared(true)

	// should be prepared by default
	ds = From("test")
	sds.True(ds.IsPrepared())
}

func (sds *selectDatasetSuite) TestGetClauses() {
	ds := From("test")
	ce := exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression(I("test")))
	sds.Equal(ce, ds.GetClauses())
}

func (sds *selectDatasetSuite) TestUpdate() {
	where := Ex{"a": 1}
	from := From("cte")
	limit := uint(1)
	order := []exp.OrderedExpression{C("a").Asc(), C("b").Desc()}
	ds := From("test").
		With("test-cte", from).
		Where(where).
		Limit(limit).
		Order(order...)
	ec := exp.NewUpdateClauses().
		SetTable(C("test")).
		CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)).
		WhereAppend(ds.GetClauses().Where()).
		SetLimit(limit).
		SetOrder(order...)
	sds.Equal(ec, ds.Update().GetClauses())
}

func (sds *selectDatasetSuite) TestInsert() {
	where := Ex{"a": 1}
	from := From("cte")
	limit := uint(1)
	order := []exp.OrderedExpression{C("a").Asc(), C("b").Desc()}
	ds := From("test").
		With("test-cte", from).
		Where(where).
		Limit(limit).
		Order(order...)
	ec := exp.NewInsertClauses().
		SetInto(C("test")).
		CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from))
	sds.Equal(ec, ds.Insert().GetClauses())
}

func (sds *selectDatasetSuite) TestDelete() {
	where := Ex{"a": 1}
	from := From("cte")
	limit := uint(1)
	order := []exp.OrderedExpression{C("a").Asc(), C("b").Desc()}
	ds := From("test").
		With("test-cte", from).
		Where(where).
		Limit(limit).
		Order(order...)
	ec := exp.NewDeleteClauses().
		SetFrom(C("test")).
		CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)).
		WhereAppend(ds.GetClauses().Where()).
		SetLimit(limit).
		SetOrder(order...)
	sds.Equal(ec, ds.Delete().GetClauses())
}

func (sds *selectDatasetSuite) TestTruncate() {
	where := Ex{"a": 1}
	from := From("cte")
	limit := uint(1)
	order := []exp.OrderedExpression{C("a").Asc(), C("b").Desc()}
	ds := From("test").
		With("test-cte", from).
		Where(where).
		Limit(limit).
		Order(order...)
	ec := exp.NewTruncateClauses().
		SetTable(exp.NewColumnListExpression("test"))
	sds.Equal(ec, ds.Truncate().GetClauses())
}

func (sds *selectDatasetSuite) TestWith() {
	from := From("cte")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.With("test-cte", from),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				CommonTablesAppend(exp.NewCommonTableExpression(false, "test-cte", from)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestWithRecursive() {
	from := From("cte")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.WithRecursive("test-cte", from),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				CommonTablesAppend(exp.NewCommonTableExpression(true, "test-cte", from)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestSelect() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Select("a", "b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("a", "b")),
		},
		selectTestCase{
			ds: bd.Select("a").Select("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("b")),
		},
		selectTestCase{
			ds: bd.Select("a").Select(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestSelectDistinct() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.SelectDistinct("a", "b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("a", "b")).
				SetDistinct(exp.NewColumnListExpression()),
		},
		selectTestCase{
			ds: bd.SelectDistinct("a").SelectDistinct("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("b")).
				SetDistinct(exp.NewColumnListExpression()),
		},
		selectTestCase{
			ds: bd.Select("a").SelectDistinct("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("b")).
				SetDistinct(exp.NewColumnListExpression()),
		},
		selectTestCase{
			ds: bd.Select("a").SelectDistinct(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression(Star())).
				SetDistinct(nil),
		},
		selectTestCase{
			ds: bd.SelectDistinct("a").SelectDistinct(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression(Star())).
				SetDistinct(nil),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestClearSelect() {
	bd := From("test").Select("a")
	sds.assertCases(
		selectTestCase{
			ds: bd.ClearSelect(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("a")),
		},
	)
}

func (sds *selectDatasetSuite) TestSelectAppend() {
	bd := From("test").Select("a")
	sds.assertCases(
		selectTestCase{
			ds: bd.SelectAppend("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("a", "b")),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetSelect(exp.NewColumnListExpression("a")),
		},
	)
}

func (sds *selectDatasetSuite) TestDistinct() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Distinct("a", "b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetDistinct(exp.NewColumnListExpression("a", "b")),
		},
		selectTestCase{
			ds: bd.Distinct("a").Distinct("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetDistinct(exp.NewColumnListExpression("b")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestFrom() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.From(T("test2")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(T("test2"))),
		},
		selectTestCase{
			ds: bd.From(From("test")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(From("test").As("t1"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestFromSelf() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.FromSelf(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(bd.As("t1"))),
		},
		selectTestCase{
			ds: bd.As("alias").FromSelf(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(bd.As("alias"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestCompoundFromSelf() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds:      bd.CompoundFromSelf(),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds:      bd.Limit(10).CompoundFromSelf(),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression(bd.Limit(10).As("t1"))),
		},
		selectTestCase{
			ds: bd.Order(C("a").Asc()).CompoundFromSelf(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(bd.Order(C("a").Asc()).As("t1"))),
		},
		selectTestCase{
			ds: bd.As("alias").FromSelf(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression(bd.As("alias"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Join(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.InnerJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestInnerJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.InnerJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.InnerJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestFullOuterJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.FullOuterJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.FullOuterJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestRightOuterJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.RightOuterJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.RightOuterJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestLeftOuterJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.LeftOuterJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.LeftOuterJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestFullJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.FullJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.FullJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestRightJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.RightJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.RightJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestLeftJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.LeftJoin(T("foo"), On(C("a").IsNull())),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewConditionedJoinExpression(exp.LeftJoinType, T("foo"), On(C("a").IsNull())),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestNaturalJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.NaturalJoin(T("foo")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewUnConditionedJoinExpression(exp.NaturalJoinType, T("foo")),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestNaturalLeftJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.NaturalLeftJoin(T("foo")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewUnConditionedJoinExpression(exp.NaturalLeftJoinType, T("foo")),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestNaturalRightJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.NaturalRightJoin(T("foo")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewUnConditionedJoinExpression(exp.NaturalRightJoinType, T("foo")),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestNaturalFullJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.NaturalFullJoin(T("foo")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewUnConditionedJoinExpression(exp.NaturalFullJoinType, T("foo")),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestCrossJoin() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.CrossJoin(T("foo")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				JoinsAppend(
					exp.NewUnConditionedJoinExpression(exp.CrossJoinType, T("foo")),
				),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestWhere() {
	w := Ex{"a": 1}
	w2 := Ex{"b": "c"}
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Where(w),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WhereAppend(w),
		},
		selectTestCase{
			ds: bd.Where(w).Where(w2),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WhereAppend(w).WhereAppend(w2),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestClearWhere() {
	w := Ex{"a": 1}
	bd := From("test").Where(w)
	sds.assertCases(
		selectTestCase{
			ds: bd.ClearWhere(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).WhereAppend(w),
		},
	)
}

func (sds *selectDatasetSuite) TestForUpdate() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.ForUpdate(NoWait),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForUpdate, NoWait)),
		},
		selectTestCase{
			ds: bd.ForUpdate(NoWait, T("table1")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForUpdate, NoWait, T("table1"))),
		},
		selectTestCase{
			ds: bd.ForUpdate(NoWait, T("table1"), T("table2")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForUpdate, NoWait, T("table1"), T("table2"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestForNoKeyUpdate() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.ForNoKeyUpdate(NoWait),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForNoKeyUpdate, NoWait)),
		},
		selectTestCase{
			ds: bd.ForNoKeyUpdate(NoWait, T("table1")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForNoKeyUpdate, NoWait, T("table1"))),
		},
		selectTestCase{
			ds: bd.ForNoKeyUpdate(NoWait, T("table1"), T("table2")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForNoKeyUpdate, NoWait, T("table1"), T("table2"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestForKeyShare() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.ForKeyShare(NoWait),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForKeyShare, NoWait)),
		},
		selectTestCase{
			ds: bd.ForKeyShare(NoWait, T("table1")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForKeyShare, NoWait, T("table1"))),
		},
		selectTestCase{
			ds: bd.ForKeyShare(NoWait, T("table1"), T("table2")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForKeyShare, NoWait, T("table1"), T("table2"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestForShare() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.ForShare(NoWait),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForShare, NoWait)),
		},
		selectTestCase{
			ds: bd.ForShare(NoWait, T("table1")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForShare, NoWait, T("table1"))),
		},
		selectTestCase{
			ds: bd.ForShare(NoWait, T("table1"), T("table2")),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLock(exp.NewLock(exp.ForShare, NoWait, T("table1"), T("table2"))),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestGroupBy() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.GroupBy("a"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetGroupBy(exp.NewColumnListExpression("a")),
		},
		selectTestCase{
			ds: bd.GroupBy("a").GroupBy("b"),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetGroupBy(exp.NewColumnListExpression("b")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestWindow() {
	w1 := W("w1").PartitionBy("a").OrderBy("b")
	w2 := W("w2").PartitionBy("a").OrderBy("b")

	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Window(w1),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w1),
		},
		selectTestCase{
			ds: bd.Window(w1).Window(w2),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w2),
		},
		selectTestCase{
			ds: bd.Window(w1, w2),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w1, w2),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestWindowAppend() {
	w1 := W("w1").PartitionBy("a").OrderBy("b")
	w2 := W("w2").PartitionBy("a").OrderBy("b")

	bd := From("test").Window(w1)
	sds.assertCases(
		selectTestCase{
			ds: bd.WindowAppend(w2),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w1, w2),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w1),
		},
	)
}

func (sds *selectDatasetSuite) TestClearWindow() {
	w1 := W("w1").PartitionBy("a").OrderBy("b")

	bd := From("test").Window(w1)
	sds.assertCases(
		selectTestCase{
			ds:      bd.ClearWindow(),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				WindowsAppend(w1),
		},
	)
}

func (sds *selectDatasetSuite) TestHaving() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Having(C("a").Gt(1)),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				HavingAppend(C("a").Gt(1)),
		},
		selectTestCase{
			ds: bd.Having(C("a").Gt(1)).Having(Ex{"b": "c"}),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				HavingAppend(C("a").Gt(1)).HavingAppend(Ex{"b": "c"}),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestOrder() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Order(C("a").Asc()),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("a").Asc()),
		},
		selectTestCase{
			ds: bd.Order(C("a").Asc()).Order(C("b").Asc()),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("b").Asc()),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestOrderAppend() {
	bd := From("test").Order(C("a").Asc())
	sds.assertCases(
		selectTestCase{
			ds: bd.OrderAppend(C("b").Asc()),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("a").Asc(), C("b").Asc()),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("a").Asc()),
		},
	)
}

func (sds *selectDatasetSuite) TestOrderPrepend() {
	bd := From("test").Order(C("a").Asc())
	sds.assertCases(
		selectTestCase{
			ds: bd.OrderPrepend(C("b").Asc()),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("b").Asc(), C("a").Asc()),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("a").Asc()),
		},
	)
}

func (sds *selectDatasetSuite) TestClearOrder() {
	bd := From("test").Order(C("a").Asc())
	sds.assertCases(
		selectTestCase{
			ds: bd.ClearOrder(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				SetOrder(C("a").Asc()),
		},
	)
}

func (sds *selectDatasetSuite) TestLimit() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Limit(10),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLimit(uint(10)),
		},
		selectTestCase{
			ds: bd.Limit(0),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds: bd.Limit(10).Limit(2),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLimit(uint(2)),
		},
		selectTestCase{
			ds: bd.Limit(10).Limit(0),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestLimitAll() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.LimitAll(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLimit(L("ALL")),
		},
		selectTestCase{
			ds: bd.Limit(10).LimitAll(),
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLimit(L("ALL")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestClearLimit() {
	bd := From("test").Limit(10)
	sds.assertCases(
		selectTestCase{
			ds:      bd.ClearLimit(),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds: bd,
			clauses: exp.NewSelectClauses().
				SetFrom(exp.NewColumnListExpression("test")).
				SetLimit(uint(10)),
		},
	)
}

func (sds *selectDatasetSuite) TestOffset() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds:      bd.Offset(10),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).SetOffset(10),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestClearOffset() {
	bd := From("test").Offset(10)
	sds.assertCases(
		selectTestCase{
			ds:      bd.ClearOffset(),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).SetOffset(10),
		},
	)
}

func (sds *selectDatasetSuite) TestUnion() {
	uds := From("union_test")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Union(uds),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				CompoundsAppend(exp.NewCompoundExpression(exp.UnionCompoundType, uds)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestUnionAll() {
	uds := From("union_test")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.UnionAll(uds),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				CompoundsAppend(exp.NewCompoundExpression(exp.UnionAllCompoundType, uds)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestIntersect() {
	uds := From("union_test")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.Intersect(uds),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				CompoundsAppend(exp.NewCompoundExpression(exp.IntersectCompoundType, uds)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestIntersectAll() {
	uds := From("union_test")
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.IntersectAll(uds),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				CompoundsAppend(exp.NewCompoundExpression(exp.IntersectAllCompoundType, uds)),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestAs() {
	bd := From("test")
	sds.assertCases(
		selectTestCase{
			ds: bd.As("t"),
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")).
				SetAlias(T("t")),
		},
		selectTestCase{
			ds:      bd,
			clauses: exp.NewSelectClauses().SetFrom(exp.NewColumnListExpression("test")),
		},
	)
}

func (sds *selectDatasetSuite) TestBuild() {
	md := new(mocks.SQLDialect)
	ds := From("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToSelectSQL", sqlB, c).Return(nil).Once()
	sql, args, err := ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Nil(err)
	md.AssertExpectations(sds.T())
}

func (sds *selectDatasetSuite) TestBuild_prepared() {
	md := new(mocks.SQLDialect)
	ds := From("test").Prepared(true).SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(true)
	md.On("ToSelectSQL", sqlB, c).Return(nil).Once()
	sql, args, err := ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Nil(err)
	md.AssertExpectations(sds.T())
}

func (sds *selectDatasetSuite) TestBuild_ReturnedError() {
	md := new(mocks.SQLDialect)
	ds := From("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	ee := errors.New("expected error")
	md.On("ToSelectSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(ee)
	}).Once()

	sql, args, err := ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Equal(ee, err)
	md.AssertExpectations(sds.T())
}

func (sds *selectDatasetSuite) TestAppendSQL() {
	md := new(mocks.SQLDialect)
	ds := From("test").SetDialect(md)
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToSelectSQL", sqlB, c).Return(nil).Once()
	ds.AppendSQL(sqlB)
	sds.NoError(sqlB.Error())
	md.AssertExpectations(sds.T())
}

func (sds *selectDatasetSuite) TestScanStructs() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).
			FromCSVString("111 Test Addr,Test1\n211 Test Addr,Test2"))

	sqlMock.ExpectQuery(`SELECT DISTINCT "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).
			FromCSVString("111 Test Addr,Test1\n211 Test Addr,Test2"))
	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))
	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))
	sqlMock.ExpectQuery(`SELECT "test" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))

	db := New("mock", mDB)
	var items []dsTestActionItem
	sds.NoError(db.From("items").ScanStructs(&items))
	sds.Equal([]dsTestActionItem{
		{Address: "111 Test Addr", Name: "Test1"},
		{Address: "211 Test Addr", Name: "Test2"},
	}, items)

	items = items[0:0]
	sds.NoError(db.From("items").Select("name").Distinct().ScanStructs(&items))
	sds.Equal([]dsTestActionItem{
		{Address: "111 Test Addr", Name: "Test1"},
		{Address: "211 Test Addr", Name: "Test2"},
	}, items)

	items = items[0:0]
	sds.EqualError(db.From("items").ScanStructs(items),
		"pp: type must be a pointer to a slice when scanning into structs")
	sds.EqualError(db.From("items").ScanStructs(&dsTestActionItem{}),
		"pp: type must be a pointer to a slice when scanning into structs")
	sds.EqualError(db.From("items").Select("test").ScanStructs(&items),
		`pp: unable to find corresponding field to column "test" returned by query`)

	sds.Equal(ErrQueryFactoryNotFoundError, From("items").ScanStructs(items))
}

func (sds *selectDatasetSuite) TestScanStructs_WithPreparedStatements() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT "address", "name" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\)`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy").
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).
			FromCSVString("111 Test Addr,Test1\n211 Test Addr,Test2"))

	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))
	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))

	sqlMock.ExpectQuery(
		`SELECT "test" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\)`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy").
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))

	db := New("mock", mDB)
	var items []dsTestActionItem
	sds.NoError(db.From("items").Prepared(true).Where(Ex{
		"name":    []string{"Bob", "Sally", "Billy"},
		"address": "111 Test Addr",
	}).ScanStructs(&items))
	sds.Equal(items, []dsTestActionItem{
		{Address: "111 Test Addr", Name: "Test1"},
		{Address: "211 Test Addr", Name: "Test2"},
	})

	items = items[0:0]
	sds.EqualError(db.From("items").ScanStructs(items),
		"pp: type must be a pointer to a slice when scanning into structs")
	sds.EqualError(db.From("items").ScanStructs(&dsTestActionItem{}),
		"pp: type must be a pointer to a slice when scanning into structs")
	sds.EqualError(db.From("items").
		Prepared(true).
		Select("test").
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		ScanStructs(&items), `pp: unable to find corresponding field to column "test" returned by query`)
}

func (sds *selectDatasetSuite) TestScanStruct() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).FromCSVString("111 Test Addr,Test1"))

	sqlMock.ExpectQuery(`SELECT DISTINCT "name" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).FromCSVString("111 Test Addr,Test1"))

	sqlMock.ExpectQuery(`SELECT "test" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))

	db := New("mock", mDB)
	var item dsTestActionItem
	found, err := db.From("items").ScanStruct(&item)
	sds.NoError(err)
	sds.True(found)
	sds.Equal("111 Test Addr", item.Address)
	sds.Equal("Test1", item.Name)

	item = dsTestActionItem{}
	found, err = db.From("items").Select("name").Distinct().ScanStruct(&item)
	sds.NoError(err)
	sds.True(found)
	sds.Equal("111 Test Addr", item.Address)
	sds.Equal("Test1", item.Name)

	_, err = db.From("items").ScanStruct(item)
	sds.EqualError(err, "pp: type must be a pointer to a struct when scanning into a struct")
	_, err = db.From("items").ScanStruct([]dsTestActionItem{})
	sds.EqualError(err, "pp: type must be a pointer to a struct when scanning into a struct")
	_, err = db.From("items").Select("test").ScanStruct(&item)
	sds.EqualError(err, `pp: unable to find corresponding field to column "test" returned by query`)

	_, err = From("items").ScanStruct(item)
	sds.Equal(ErrQueryFactoryNotFoundError, err)
}

func (sds *selectDatasetSuite) TestScanStruct_WithPreparedStatements() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT "address", "name" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\) LIMIT \?`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy", 1).
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).FromCSVString("111 Test Addr,Test1"))

	sqlMock.ExpectQuery(`SELECT "test" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\) LIMIT \?`).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy", 1).
		WillReturnRows(sqlmock.NewRows([]string{"test"}).FromCSVString("test1\ntest2"))

	db := New("mock", mDB)
	var item dsTestActionItem
	found, err := db.From("items").Prepared(true).Where(Ex{
		"name":    []string{"Bob", "Sally", "Billy"},
		"address": "111 Test Addr",
	}).ScanStruct(&item)
	sds.NoError(err)
	sds.True(found)
	sds.Equal("111 Test Addr", item.Address)
	sds.Equal("Test1", item.Name)

	_, err = db.From("items").ScanStruct(item)
	sds.EqualError(err, "pp: type must be a pointer to a struct when scanning into a struct")
	_, err = db.From("items").ScanStruct([]dsTestActionItem{})
	sds.EqualError(err, "pp: type must be a pointer to a struct when scanning into a struct")
	_, err = db.From("items").
		Prepared(true).
		Select("test").
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		ScanStruct(&item)
	sds.EqualError(err, `pp: unable to find corresponding field to column "test" returned by query`)
}

func (sds *selectDatasetSuite) TestScanStructUntagged() {
	defer SetIgnoreUntaggedFields(false)

	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "address", "name", "untagged" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name", "untagged"}).FromCSVString("111 Test Addr,Test1,Test2"))

	sqlMock.ExpectQuery(`SELECT "address", "name" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"address", "name"}).FromCSVString("111 Test Addr,Test1"))

	db := New("mock", mDB)
	var item dsUntaggedTestActionItem

	found, err := db.From("items").ScanStruct(&item)
	sds.NoError(err)
	sds.True(found)
	sds.Equal("111 Test Addr", item.Address)
	sds.Equal("Test1", item.Name)
	sds.Equal("Test2", item.Untagged)

	// Ignore untagged fields, which will suppress the "untagged" column
	SetIgnoreUntaggedFields(true)

	item = dsUntaggedTestActionItem{}
	found, err = db.From("items").ScanStruct(&item)
	sds.NoError(err)
	sds.True(found)
	sds.Equal("111 Test Addr", item.Address)
	sds.Equal("Test1", item.Name)
	sds.Equal("", item.Untagged)
}

func (sds *selectDatasetSuite) TestScanVals() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "id" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))
	sqlMock.ExpectQuery(`SELECT \* FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))
	sqlMock.ExpectQuery(`SELECT \* FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))

	db := New("mock", mDB)
	var ids []uint32
	sds.NoError(db.From("items").Select("id").ScanVals(&ids))
	sds.Equal(ids, []uint32{1, 2, 3, 4, 5})

	sds.EqualError(db.From("items").ScanVals([]uint32{}),
		"pp: type must be a pointer to a slice when scanning into vals")
	sds.EqualError(db.From("items").ScanVals(dsTestActionItem{}),
		"pp: type must be a pointer to a slice when scanning into vals")

	err = From("items").ScanVals(&ids)
	sds.Equal(ErrQueryFactoryNotFoundError, err)
}

func (sds *selectDatasetSuite) TestScanVals_WithPreparedStatment() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT "id" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\)`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))

	sqlMock.ExpectQuery(`SELECT \* FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))
	sqlMock.ExpectQuery(`SELECT \* FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("1\n2\n3\n4\n5"))

	db := New("mock", mDB)
	var ids []uint32
	sds.NoError(db.From("items").
		Prepared(true).
		Select("id").
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		ScanVals(&ids))
	sds.Equal([]uint32{1, 2, 3, 4, 5}, ids)

	sds.EqualError(db.From("items").ScanVals([]uint32{}),
		"pp: type must be a pointer to a slice when scanning into vals")

	sds.EqualError(db.From("items").ScanVals(dsTestActionItem{}),
		"pp: type must be a pointer to a slice when scanning into vals")
}

func (sds *selectDatasetSuite) TestScanVal() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "id" FROM "items" LIMIT 1`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("10"))

	db := New("mock", mDB)
	var id int64
	found, err := db.From("items").Select("id").ScanVal(&id)
	sds.NoError(err)
	sds.Equal(id, int64(10))
	sds.True(found)

	found, err = db.From("items").ScanVal([]int64{})
	sds.False(found)
	sds.EqualError(err, "pp: type must be a pointer when scanning into val")
	found, err = db.From("items").ScanVal(10)
	sds.False(found)
	sds.EqualError(err, "pp: type must be a pointer when scanning into val")

	_, err = From("items").ScanVal(&id)
	sds.Equal(ErrQueryFactoryNotFoundError, err)
}

func (sds *selectDatasetSuite) TestScanVal_WithPreparedStatement() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT "id" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\) LIMIT ?`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString("10"))

	db := New("mock", mDB)
	var id int64
	found, err := db.From("items").
		Prepared(true).
		Select("id").
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		ScanVal(&id)
	sds.NoError(err)
	sds.Equal(int64(10), id)
	sds.True(found)

	found, err = db.From("items").ScanVal([]int64{})
	sds.False(found)
	sds.EqualError(err, "pp: type must be a pointer when scanning into val")
	found, err = db.From("items").ScanVal(10)
	sds.False(found)
	sds.EqualError(err, "pp: type must be a pointer when scanning into val")
}

func (sds *selectDatasetSuite) TestCount() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT COUNT\(\*\) AS "count" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("10"))

	db := New("mock", mDB)
	count, err := db.From("items").Count()
	sds.NoError(err)
	sds.Equal(count, int64(10))
}

func (sds *selectDatasetSuite) TestCount_WithPreparedStatement() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT COUNT\(\*\) AS "count" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\)`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy", 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("10"))

	ds := New("mock", mDB)
	count, err := ds.From("items").
		Prepared(true).
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		Count()
	sds.NoError(err)
	sds.Equal(int64(10), count)
}

func (sds *selectDatasetSuite) TestPluck() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(`SELECT "name" FROM "items"`).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"name"}).FromCSVString("test1\ntest2\ntest3\ntest4\ntest5"))

	db := New("mock", mDB)
	var names []string
	sds.NoError(db.From("items").Pluck(&names, "name"))
	sds.Equal([]string{"test1", "test2", "test3", "test4", "test5"}, names)
}

func (sds *selectDatasetSuite) TestPluck_WithPreparedStatement() {
	mDB, sqlMock, err := sqlmock.New()
	sds.NoError(err)
	sqlMock.ExpectQuery(
		`SELECT "name" FROM "items" WHERE \(\("address" = \?\) AND \("name" IN \(\?, \?, \?\)\)\)`,
	).
		WithArgs("111 Test Addr", "Bob", "Sally", "Billy").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).FromCSVString("Bob\nSally\nBilly"))

	db := New("mock", mDB)
	var names []string
	sds.NoError(db.From("items").
		Prepared(true).
		Where(Ex{"name": []string{"Bob", "Sally", "Billy"}, "address": "111 Test Addr"}).
		Pluck(&names, "name"))
	sds.Equal([]string{"Bob", "Sally", "Billy"}, names)
}

func (sds *selectDatasetSuite) TestSetError() {
	err1 := errors.New("error #1")
	err2 := errors.New("error #2")
	err3 := errors.New("error #3")

	// Verify initial error set/get works properly
	md := new(mocks.SQLDialect)
	ds := From("test").SetDialect(md)
	ds = ds.SetError(err1)
	sds.Equal(err1, ds.Error())
	sql, args, err := ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Equal(err1, err)

	// Repeated SetError calls on Dataset should not overwrite the original error
	ds = ds.SetError(err2)
	sds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Equal(err1, err)

	// Builder functions should not lose the error
	ds = ds.ClearWindow()
	sds.Equal(err1, ds.Error())
	sql, args, err = ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Equal(err1, err)

	// Deeper errors inside SQL generation should still return original error
	c := ds.GetClauses()
	sqlB := builder.NewSQLBuilder(false)
	md.On("ToInsertSQL", sqlB, c).Run(func(args mock.Arguments) {
		args.Get(0).(builder.SQLBuilder).SetError(err3)
	}).Once()

	sql, args, err = ds.Build()
	sds.Empty(sql)
	sds.Empty(args)
	sds.Equal(err1, err)
}

func TestSelectDataset(t *testing.T) {
	suite.Run(t, new(selectDatasetSuite))
}
