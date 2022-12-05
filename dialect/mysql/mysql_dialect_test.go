package mysql_test

import (
	"regexp"
	"testing"

	"github.com/sllt/pp"
	"github.com/sllt/pp/exp"
	"github.com/stretchr/testify/suite"
)

type (
	mysqlDialectSuite struct {
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

func (mds *mysqlDialectSuite) GetDs(table string) *pp.SelectDataset {
	return pp.Dialect("mysql").From(table)
}

func (mds *mysqlDialectSuite) assertSQL(cases ...sqlTestCase) {
	for i, c := range cases {
		actualSQL, actualArgs, err := c.ds.Build()
		if c.err == "" {
			mds.NoError(err, "test case %d failed", i)
		} else {
			mds.EqualError(err, c.err, "test case %d failed", i)
		}
		mds.Equal(c.sql, actualSQL, "test case %d failed", i)
		if c.isPrepared && c.args != nil || len(c.args) > 0 {
			mds.Equal(c.args, actualArgs, "test case %d failed", i)
		} else {
			mds.Empty(actualArgs, "test case %d failed", i)
		}
	}
}

func (mds *mysqlDialectSuite) TestIdentifiers() {
	ds := mds.GetDs("test")
	mds.assertSQL(
		sqlTestCase{ds: ds.Select(
			"a",
			pp.I("a.b.c"),
			pp.I("c.d"),
			pp.C("test").As("test"),
		), sql: "SELECT `a`, `a`.`b`.`c`, `c`.`d`, `test` AS `test` FROM `test`"},
	)
}

func (mds *mysqlDialectSuite) TestLiteralString() {
	ds := mds.GetDs("test")
	col := pp.C("a")
	mds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq("test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(col.Eq("test'test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\'test')"},
		sqlTestCase{ds: ds.Where(col.Eq(`test"test`)), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\"test')"},
		sqlTestCase{ds: ds.Where(col.Eq(`test\test`)), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\\test')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\ntest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ntest')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\rtest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\rtest')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\x00test")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\x00test')"},
		sqlTestCase{ds: ds.Where(col.Eq("test\x1atest")), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\x1atest')"},
	)
}

func (mds *mysqlDialectSuite) TestLiteralBytes() {
	col := pp.C("a")
	ds := mds.GetDs("test")
	mds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test'test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\'test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte(`test"test`))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\"test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte(`test\test`))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\\\test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\ntest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\ntest')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\rtest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\rtest')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\x00test"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\x00test')"},
		sqlTestCase{ds: ds.Where(col.Eq([]byte("test\x1atest"))), sql: "SELECT * FROM `test` WHERE (`a` = 'test\\x1atest')"},
	)
}

func (mds *mysqlDialectSuite) TestBooleanOperations() {
	col := pp.C("a")
	ds := mds.GetDs("test")
	mds.assertSQL(
		sqlTestCase{ds: ds.Where(col.Eq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.Eq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.Is(true)), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.Is(false)), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS FALSE)"},
		sqlTestCase{ds: ds.Where(col.Neq(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.Neq(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsNot(true)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsNot(false)), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.IsNotTrue()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT TRUE)"},
		sqlTestCase{ds: ds.Where(col.IsNotFalse()), sql: "SELECT * FROM `test` WHERE (`a` IS NOT FALSE)"},
		sqlTestCase{ds: ds.Where(col.Like("a%")), sql: "SELECT * FROM `test` WHERE (`a` LIKE BINARY 'a%')"},
		sqlTestCase{ds: ds.Where(col.NotLike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT LIKE BINARY 'a%')"},
		sqlTestCase{ds: ds.Where(col.ILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.NotILike("a%")), sql: "SELECT * FROM `test` WHERE (`a` NOT LIKE 'a%')"},
		sqlTestCase{ds: ds.Where(col.Like(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` REGEXP BINARY '[ab]')"},
		sqlTestCase{ds: ds.Where(col.NotLike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` NOT REGEXP BINARY '[ab]')"},
		sqlTestCase{ds: ds.Where(col.ILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` REGEXP '[ab]')"},
		sqlTestCase{ds: ds.Where(col.NotILike(regexp.MustCompile("[ab]"))), sql: "SELECT * FROM `test` WHERE (`a` NOT REGEXP '[ab]')"},
	)
}

func (mds *mysqlDialectSuite) TestBitwiseOperations() {
	col := pp.C("a")
	ds := mds.GetDs("test")
	mds.assertSQL(
		sqlTestCase{ds: ds.Where(col.BitwiseInversion()), sql: "SELECT * FROM `test` WHERE (~ `a`)"},
		sqlTestCase{ds: ds.Where(col.BitwiseAnd(1)), sql: "SELECT * FROM `test` WHERE (`a` & 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseOr(1)), sql: "SELECT * FROM `test` WHERE (`a` | 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseXor(1)), sql: "SELECT * FROM `test` WHERE (`a` ^ 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseLeftShift(1)), sql: "SELECT * FROM `test` WHERE (`a` << 1)"},
		sqlTestCase{ds: ds.Where(col.BitwiseRightShift(1)), sql: "SELECT * FROM `test` WHERE (`a` >> 1)"},
	)
}

func (mds *mysqlDialectSuite) TestUpdateSQL() {
	ds := mds.GetDs("test").Update()
	mds.assertSQL(
		sqlTestCase{
			ds: ds.
				Set(pp.Record{"foo": "bar"}).
				From("test_2").
				Where(pp.I("test.id").Eq(pp.I("test_2.test_id"))),
			sql: "UPDATE `test`,`test_2` SET `foo`='bar' WHERE (`test`.`id` = `test_2`.`test_id`)",
		},
	)
}

func TestDatasetAdapterSuite(t *testing.T) {
	suite.Run(t, new(mysqlDialectSuite))
}
