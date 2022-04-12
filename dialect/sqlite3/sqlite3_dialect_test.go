package sqlite3_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/suite"
	"manlu.org/pp"
	"manlu.org/pp/exp"
)

type (
	sqlite3DialectSuite struct {
		suite.Suite
	}
	sqlTestCase struct {
		ds         exp.SQLExpression
		sql        string
		err        string
		isPrepared bool
		args       []interface{}
	}
)

func (sds *sqlite3DialectSuite) GetDs(table string) *pp.SelectDataset {
	return pp.Dialect("sqlite3").From(table)
}

func (sds *sqlite3DialectSuite) assertSQL(cases ...sqlTestCase) {
	for i, c := range cases {
		actualSQL, actualArgs, err := c.ds.Build()
		if c.err == "" {
			sds.NoError(err, "test case %d failed", i)
		} else {
			sds.EqualError(err, c.err, "test case %d failed", i)
		}
		sds.Equal(c.sql, actualSQL, "test case %d failed", i)
		if c.isPrepared && c.args != nil || len(c.args) > 0 {
			sds.Equal(c.args, actualArgs, "test case %d failed", i)
		} else {
			sds.Empty(actualArgs, "test case %d failed", i)
		}
	}
}

func (sds *sqlite3DialectSuite) TestIdentifiers() {
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Select(
			"a",
			pp.I("a.b.c"),
			pp.I("c.d"),
			pp.C("test").As("test"),
		), sql: "SELECT `a`, `a`.`b`.`c`, `c`.`d`, `test` AS `test` FROM `test`"},
	)
}

func (sds *sqlite3DialectSuite) TestUpdateSQL_multipleTables() {
	ds := sds.GetDs("test").Update()
	sds.assertSQL(
		sqlTestCase{
			ds: ds.
				Set(pp.Record{"foo": "bar"}).
				From("test_2").
				Where(pp.I("test.id").Eq(pp.I("test_2.test_id"))),
			err: "pp: sqlite3 dialect does not support multiple tables in UPDATE",
		},
	)
}

func (sds *sqlite3DialectSuite) TestCompoundExpressions() {
	ds1 := sds.GetDs("test").Select("a")
	ds2 := sds.GetDs("test2").Select("b")
	sds.assertSQL(
		sqlTestCase{ds: ds1.Union(ds2), sql: "SELECT `a` FROM `test` UNION SELECT `b` FROM `test2`"},
		sqlTestCase{ds: ds1.UnionAll(ds2), sql: "SELECT `a` FROM `test` UNION ALL SELECT `b` FROM `test2`"},
		sqlTestCase{ds: ds1.Intersect(ds2), sql: "SELECT `a` FROM `test` INTERSECT SELECT `b` FROM `test2`"},
	)
}

func (sds *sqlite3DialectSuite) TestLiteralString() {
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test'test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test''test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(`test"test`)), sql: "SELECT * FROM `test` WHERE (`a` = 'test\"test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(`test\test`)), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test\ntest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\ntest')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test\rtest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\rtest')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test\x00test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\x00test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq("test\x1atest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\x1atest')"},
	)
}

func (sds *sqlite3DialectSuite) TestLiteralBytes() {
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test'test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test''test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte(`test"test`))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\"test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte(`test\test`))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test\ntest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\ntest')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test\rtest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\rtest')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test\x00test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\x00test')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq([]byte("test\x1atest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\x1atest')"},
	)
}

func (sds *sqlite3DialectSuite) TestBooleanOperations() {
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Is(true)), sql: "SELECT * FROM `test` WHERE (`a` IS 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Is(false)), sql: "SELECT * FROM `test` WHERE (`a` IS 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Neq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Neq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsNot(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsNot(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsNotTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").IsNotFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT 0)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Like("a%")), sql: "SELECT * FROM `test` WHERE (`a` LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(pp.C("a").NotLike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(pp.C("a").ILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(pp.C("a").NotILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(pp.C("a").Like(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` REGEXP '[ab]')"},
		sqlTestCase{ds: ds.Where(pp.C("a").NotLike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` NOT REGEXP '[ab]')"},
		sqlTestCase{ds: ds.Where(pp.C("a").ILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` REGEXP '[ab]')"},
		sqlTestCase{ds: ds.Where(pp.C("a").NotILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` NOT REGEXP '[ab]')"},
	)
}

func (sds *sqlite3DialectSuite) TestBitwiseOperations() {
	col := pp.C("a")
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Where(col.BitwiseInversion()), err: "pp: bitwise operator 'Inversion' not supported"},
		sqlTestCase{ds: ds.Where(col.BitwiseAnd(1)), sql: "SELECT * FROM `test` WHERE (`a` & 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseOr(1)), sql: "SELECT * FROM `test` WHERE (`a` | 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseXor(1)), err: "pp: bitwise operator 'XOR' not supported"},
		sqlTestCase{ds: ds.Where(col.BitwiseLeftShift(1)), sql: "SELECT * FROM `test` WHERE (`a` << 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseRightShift(1)), sql: "SELECT * FROM `test` WHERE (`a` >> 1)"},
	)
}

func (sds *sqlite3DialectSuite) TestForUpdate() {
	ds := sds.GetDs("test")
	sds.assertSQL(
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(1)).ForUpdate(pp.Wait), sql: "SELECT * FROM `test` WHERE (`a` = 1)"},
		sqlTestCase{ds: ds.Where(pp.C("a").Eq(1)).ForUpdate(pp.NoWait), sql: "SELECT * FROM `test` WHERE (`a` = 1)"},
	)
}

func TestDatasetAdapterSuite(t *testing.T) {
	suite.Run(t, new(sqlite3DialectSuite))
}
