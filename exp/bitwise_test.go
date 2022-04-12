package exp

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type bitwiseExpressionSuite struct {
	suite.Suite
}

func TestBitwiseExpressionSuite(t *testing.T) {
	suite.Run(t, &bitwiseExpressionSuite{})
}

func (bes *bitwiseExpressionSuite) TestClone() {
	be := NewBitwiseExpression(BitwiseAndOp, NewIdentifierExpression("", "", "col"), 1)
	bes.Equal(be, be.Clone())
}

func (bes *bitwiseExpressionSuite) TestExpression() {
	be := NewBitwiseExpression(BitwiseAndOp, NewIdentifierExpression("", "", "col"), 1)
	bes.Equal(be, be.Expression())
}

func (bes *bitwiseExpressionSuite) TestAs() {
	be := NewBitwiseExpression(BitwiseInversionOp, NewIdentifierExpression("", "", "col"), 1)
	bes.Equal(NewAliasExpression(be, "a"), be.As("a"))
}

func (bes *bitwiseExpressionSuite) TestAsc() {
	be := NewBitwiseExpression(BitwiseAndOp, NewIdentifierExpression("", "", "col"), 1)
	bes.Equal(NewOrderedExpression(be, AscDir, NoNullsSortType), be.Asc())
}

func (bes *bitwiseExpressionSuite) TestDesc() {
	be := NewBitwiseExpression(BitwiseOrOp, NewIdentifierExpression("", "", "col"), 1)
	bes.Equal(NewOrderedExpression(be, DescSortDir, NoNullsSortType), be.Desc())
}

func (bes *bitwiseExpressionSuite) TestAllOthers() {
	be := NewBitwiseExpression(BitwiseRightShiftOp, NewIdentifierExpression("", "", "col"), 1)
	rv := NewRangeVal(1, 2)
	pattern := "bitwiseExp like%"
	inVals := []interface{}{1, 2}
	testCases := []struct {
		Ex       Expression
		Expected Expression
	}{
		{Ex: be.Eq(1), Expected: NewBooleanExpression(EqOp, be, 1)},
		{Ex: be.Neq(1), Expected: NewBooleanExpression(NeqOp, be, 1)},
		{Ex: be.Gt(1), Expected: NewBooleanExpression(GtOp, be, 1)},
		{Ex: be.Gte(1), Expected: NewBooleanExpression(GteOp, be, 1)},
		{Ex: be.Lt(1), Expected: NewBooleanExpression(LtOp, be, 1)},
		{Ex: be.Lte(1), Expected: NewBooleanExpression(LteOp, be, 1)},
		{Ex: be.Between(rv), Expected: NewRangeExpression(BetweenOp, be, rv)},
		{Ex: be.NotBetween(rv), Expected: NewRangeExpression(NotBetweenOp, be, rv)},
		{Ex: be.Like(pattern), Expected: NewBooleanExpression(LikeOp, be, pattern)},
		{Ex: be.NotLike(pattern), Expected: NewBooleanExpression(NotLikeOp, be, pattern)},
		{Ex: be.ILike(pattern), Expected: NewBooleanExpression(ILikeOp, be, pattern)},
		{Ex: be.NotILike(pattern), Expected: NewBooleanExpression(NotILikeOp, be, pattern)},
		{Ex: be.RegexpLike(pattern), Expected: NewBooleanExpression(RegexpLikeOp, be, pattern)},
		{Ex: be.RegexpNotLike(pattern), Expected: NewBooleanExpression(RegexpNotLikeOp, be, pattern)},
		{Ex: be.RegexpILike(pattern), Expected: NewBooleanExpression(RegexpILikeOp, be, pattern)},
		{Ex: be.RegexpNotILike(pattern), Expected: NewBooleanExpression(RegexpNotILikeOp, be, pattern)},
		{Ex: be.In(inVals), Expected: NewBooleanExpression(InOp, be, inVals)},
		{Ex: be.NotIn(inVals), Expected: NewBooleanExpression(NotInOp, be, inVals)},
		{Ex: be.Is(true), Expected: NewBooleanExpression(IsOp, be, true)},
		{Ex: be.IsNot(true), Expected: NewBooleanExpression(IsNotOp, be, true)},
		{Ex: be.IsNull(), Expected: NewBooleanExpression(IsOp, be, nil)},
		{Ex: be.IsNotNull(), Expected: NewBooleanExpression(IsNotOp, be, nil)},
		{Ex: be.IsTrue(), Expected: NewBooleanExpression(IsOp, be, true)},
		{Ex: be.IsNotTrue(), Expected: NewBooleanExpression(IsNotOp, be, true)},
		{Ex: be.IsFalse(), Expected: NewBooleanExpression(IsOp, be, false)},
		{Ex: be.IsNotFalse(), Expected: NewBooleanExpression(IsNotOp, be, false)},
		{Ex: be.Distinct(), Expected: NewSQLFunctionExpression("DISTINCT", be)},
	}

	for _, tc := range testCases {
		bes.Equal(tc.Expected, tc.Ex)
	}
}
