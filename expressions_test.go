package pp_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"manlu.org/pp"
	"manlu.org/pp/exp"
)

type (
	ppExpressionsSuite struct {
		suite.Suite
	}
)

func (ges *ppExpressionsSuite) TestCast() {
	ges.Equal(exp.NewCastExpression(pp.C("test"), "string"), pp.Cast(pp.C("test"), "string"))
}

func (ges *ppExpressionsSuite) TestDoNothing() {
	ges.Equal(exp.NewDoNothingConflictExpression(), pp.DoNothing())
}

func (ges *ppExpressionsSuite) TestDoUpdate() {
	ges.Equal(exp.NewDoUpdateConflictExpression("test", pp.Record{"a": "b"}), pp.DoUpdate("test", pp.Record{"a": "b"}))
}

func (ges *ppExpressionsSuite) TestOr() {
	e1 := pp.C("a").Eq("b")
	e2 := pp.C("b").Eq(2)
	ges.Equal(exp.NewExpressionList(exp.OrType, e1, e2), pp.Or(e1, e2))
}

func (ges *ppExpressionsSuite) TestAnd() {
	e1 := pp.C("a").Eq("b")
	e2 := pp.C("b").Eq(2)
	ges.Equal(exp.NewExpressionList(exp.AndType, e1, e2), pp.And(e1, e2))
}

func (ges *ppExpressionsSuite) TestFunc() {
	ges.Equal(exp.NewSQLFunctionExpression("count", pp.L("*")), pp.Func("count", pp.L("*")))
}

func (ges *ppExpressionsSuite) TestDISTINCT() {
	ges.Equal(exp.NewSQLFunctionExpression("DISTINCT", pp.I("col")), pp.DISTINCT("col"))
}

func (ges *ppExpressionsSuite) TestCOUNT() {
	ges.Equal(exp.NewSQLFunctionExpression("COUNT", pp.I("col")), pp.COUNT("col"))
}

func (ges *ppExpressionsSuite) TestMIN() {
	ges.Equal(exp.NewSQLFunctionExpression("MIN", pp.I("col")), pp.MIN("col"))
}

func (ges *ppExpressionsSuite) TestMAX() {
	ges.Equal(exp.NewSQLFunctionExpression("MAX", pp.I("col")), pp.MAX("col"))
}

func (ges *ppExpressionsSuite) TestAVG() {
	ges.Equal(exp.NewSQLFunctionExpression("AVG", pp.I("col")), pp.AVG("col"))
}

func (ges *ppExpressionsSuite) TestFIRST() {
	ges.Equal(exp.NewSQLFunctionExpression("FIRST", pp.I("col")), pp.FIRST("col"))
}

func (ges *ppExpressionsSuite) TestLAST() {
	ges.Equal(exp.NewSQLFunctionExpression("LAST", pp.I("col")), pp.LAST("col"))
}

func (ges *ppExpressionsSuite) TestSUM() {
	ges.Equal(exp.NewSQLFunctionExpression("SUM", pp.I("col")), pp.SUM("col"))
}

func (ges *ppExpressionsSuite) TestCOALESCE() {
	ges.Equal(exp.NewSQLFunctionExpression("COALESCE", pp.I("col"), nil), pp.COALESCE(pp.I("col"), nil))
}

func (ges *ppExpressionsSuite) TestROW_NUMBER() {
	ges.Equal(exp.NewSQLFunctionExpression("ROW_NUMBER"), pp.ROW_NUMBER())
}

func (ges *ppExpressionsSuite) TestRANK() {
	ges.Equal(exp.NewSQLFunctionExpression("RANK"), pp.RANK())
}

func (ges *ppExpressionsSuite) TestDENSE_RANK() {
	ges.Equal(exp.NewSQLFunctionExpression("DENSE_RANK"), pp.DENSE_RANK())
}

func (ges *ppExpressionsSuite) TestPERCENT_RANK() {
	ges.Equal(exp.NewSQLFunctionExpression("PERCENT_RANK"), pp.PERCENT_RANK())
}

func (ges *ppExpressionsSuite) TestCUME_DIST() {
	ges.Equal(exp.NewSQLFunctionExpression("CUME_DIST"), pp.CUME_DIST())
}

func (ges *ppExpressionsSuite) TestNTILE() {
	ges.Equal(exp.NewSQLFunctionExpression("NTILE", 1), pp.NTILE(1))
}

func (ges *ppExpressionsSuite) TestFIRST_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("FIRST_VALUE", pp.I("col")), pp.FIRST_VALUE("col"))
}

func (ges *ppExpressionsSuite) TestLAST_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("LAST_VALUE", pp.I("col")), pp.LAST_VALUE("col"))
}

func (ges *ppExpressionsSuite) TestNTH_VALUE() {
	ges.Equal(exp.NewSQLFunctionExpression("NTH_VALUE", pp.I("col"), 1), pp.NTH_VALUE("col", 1))
	ges.Equal(exp.NewSQLFunctionExpression("NTH_VALUE", pp.I("col"), 1), pp.NTH_VALUE(pp.C("col"), 1))
}

func (ges *ppExpressionsSuite) TestI() {
	ges.Equal(exp.NewIdentifierExpression("s", "t", "c"), pp.I("s.t.c"))
}

func (ges *ppExpressionsSuite) TestC() {
	ges.Equal(exp.NewIdentifierExpression("", "", "c"), pp.C("c"))
}

func (ges *ppExpressionsSuite) TestS() {
	ges.Equal(exp.NewIdentifierExpression("s", "", ""), pp.S("s"))
}

func (ges *ppExpressionsSuite) TestT() {
	ges.Equal(exp.NewIdentifierExpression("", "t", ""), pp.T("t"))
}

func (ges *ppExpressionsSuite) TestW() {
	ges.Equal(exp.NewWindowExpression(nil, nil, nil, nil), pp.W())
	ges.Equal(exp.NewWindowExpression(pp.I("a"), nil, nil, nil), pp.W("a"))
	ges.Equal(exp.NewWindowExpression(pp.I("a"), pp.I("b"), nil, nil), pp.W("a", "b"))
	ges.Equal(exp.NewWindowExpression(pp.I("a"), pp.I("b"), nil, nil), pp.W("a", "b", "c"))
}

func (ges *ppExpressionsSuite) TestOn() {
	ges.Equal(exp.NewJoinOnCondition(pp.Ex{"a": "b"}), pp.On(pp.Ex{"a": "b"}))
}

func (ges *ppExpressionsSuite) TestUsing() {
	ges.Equal(exp.NewJoinUsingCondition("a", "b"), pp.Using("a", "b"))
}

func (ges *ppExpressionsSuite) TestL() {
	ges.Equal(exp.NewLiteralExpression("? + ?", 1, 2), pp.L("? + ?", 1, 2))
}

func (ges *ppExpressionsSuite) TestLiteral() {
	ges.Equal(exp.NewLiteralExpression("? + ?", 1, 2), pp.Literal("? + ?", 1, 2))
}

func (ges *ppExpressionsSuite) TestV() {
	ges.Equal(exp.NewLiteralExpression("?", "a"), pp.V("a"))
}

func (ges *ppExpressionsSuite) TestRange() {
	ges.Equal(exp.NewRangeVal("a", "b"), pp.Range("a", "b"))
}

func (ges *ppExpressionsSuite) TestStar() {
	ges.Equal(exp.NewLiteralExpression("*"), pp.Star())
}

func (ges *ppExpressionsSuite) TestDefault() {
	ges.Equal(exp.Default(), pp.Default())
}

func (ges *ppExpressionsSuite) TestLateral() {
	ds := pp.From("test")
	ges.Equal(exp.NewLateralExpression(ds), pp.Lateral(ds))
}

func (ges *ppExpressionsSuite) TestAny() {
	ds := pp.From("test").Select("id")
	ges.Equal(exp.NewSQLFunctionExpression("ANY ", ds), pp.Any(ds))
}

func (ges *ppExpressionsSuite) TestAll() {
	ds := pp.From("test").Select("id")
	ges.Equal(exp.NewSQLFunctionExpression("ALL ", ds), pp.All(ds))
}

func TestPpExpressions(t *testing.T) {
	suite.Run(t, new(ppExpressionsSuite))
}
