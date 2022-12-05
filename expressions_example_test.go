// nolint:lll // sql statements are long
package pp_test

import (
	"fmt"
	"github.com/sllt/pp"
	"regexp"

	"github.com/sllt/pp/exp"
)

func ExampleAVG() {
	ds := pp.From("test").Select(pp.AVG("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT AVG("col") FROM "test" []
	// SELECT AVG("col") FROM "test" []
}

func ExampleAVG_as() {
	sql, _, _ := pp.From("test").Select(pp.AVG("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT AVG("a") AS "a" FROM "test"
}

func ExampleAVG_havingClause() {
	ds := pp.
		From("test").
		Select(pp.AVG("a").As("avg")).
		GroupBy("a").
		Having(pp.AVG("a").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT AVG("a") AS "avg" FROM "test" GROUP BY "a" HAVING (AVG("a") > 10) []
	// SELECT AVG("a") AS "avg" FROM "test" GROUP BY "a" HAVING (AVG("a") > ?) [10]
}

func ExampleAnd() {
	ds := pp.From("test").Where(
		pp.And(
			pp.C("col").Gt(10),
			pp.C("col").Lt(20),
		),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col" > 10) AND ("col" < 20)) []
	// SELECT * FROM "test" WHERE (("col" > ?) AND ("col" < ?)) [10 20]
}

// You can use And with Or to create more complex queries
func ExampleAnd_withOr() {
	ds := pp.From("test").Where(
		pp.And(
			pp.C("col1").IsTrue(),
			pp.Or(
				pp.C("col2").Gt(10),
				pp.C("col2").Lt(20),
			),
		),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// by default expressions are anded together
	ds = pp.From("test").Where(
		pp.C("col1").IsTrue(),
		pp.Or(
			pp.C("col2").Gt(10),
			pp.C("col2").Lt(20),
		),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col2" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col2" < ?))) [10 20]
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col2" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col2" < ?))) [10 20]
}

// You can use ExOr inside of And expression lists.
func ExampleAnd_withExOr() {
	// by default expressions are anded together
	ds := pp.From("test").Where(
		pp.C("col1").IsTrue(),
		pp.ExOr{
			"col2": pp.Op{"gt": 10},
			"col3": pp.Op{"lt": 20},
		},
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > 10) OR ("col3" < 20))) []
	// SELECT * FROM "test" WHERE (("col1" IS TRUE) AND (("col2" > ?) OR ("col3" < ?))) [10 20]
}

func ExampleC() {
	sql, args, _ := pp.From("test").
		Select(pp.C("*")).
		Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").
		Select(pp.C("col1")).
		Build()
	fmt.Println(sql, args)

	ds := pp.From("test").Where(
		pp.C("col1").Eq(10),
		pp.C("col2").In([]int64{1, 2, 3, 4}),
		pp.C("col3").Like(regexp.MustCompile("^[ab]")),
		pp.C("col4").IsNull(),
	)

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
	// SELECT "col1" FROM "test" []
	// SELECT * FROM "test" WHERE (("col1" = 10) AND ("col2" IN (1, 2, 3, 4)) AND ("col3" ~ '^[ab]') AND ("col4" IS NULL)) []
	// SELECT * FROM "test" WHERE (("col1" = ?) AND ("col2" IN (?, ?, ?, ?)) AND ("col3" ~ ?) AND ("col4" IS NULL)) [10 1 2 3 4 ^[ab]]
}

func ExampleC_as() {
	sql, _, _ := pp.From("test").Select(pp.C("a").As("as_a")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Select(pp.C("a").As(pp.C("as_a"))).Build()
	fmt.Println(sql)

	// Output:
	// SELECT "a" AS "as_a" FROM "test"
	// SELECT "a" AS "as_a" FROM "test"
}

func ExampleC_ordering() {
	sql, args, _ := pp.From("test").Order(pp.C("a").Asc()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Order(pp.C("a").Asc().NullsFirst()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Order(pp.C("a").Asc().NullsLast()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Order(pp.C("a").Desc()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Order(pp.C("a").Desc().NullsFirst()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Order(pp.C("a").Desc().NullsLast()).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC []
	// SELECT * FROM "test" ORDER BY "a" ASC NULLS FIRST []
	// SELECT * FROM "test" ORDER BY "a" ASC NULLS LAST []
	// SELECT * FROM "test" ORDER BY "a" DESC []
	// SELECT * FROM "test" ORDER BY "a" DESC NULLS FIRST []
	// SELECT * FROM "test" ORDER BY "a" DESC NULLS LAST []
}

func ExampleC_cast() {
	sql, _, _ := pp.From("test").
		Select(pp.C("json1").Cast("TEXT").As("json_text")).
		Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.C("json1").Cast("TEXT").Neq(
			pp.C("json2").Cast("TEXT"),
		),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT CAST("json1" AS TEXT) AS "json_text" FROM "test"
	// SELECT * FROM "test" WHERE (CAST("json1" AS TEXT) != CAST("json2" AS TEXT))
}

func ExampleC_comparisons() {
	// used from an identifier
	sql, _, _ := pp.From("test").Where(pp.C("a").Eq(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Neq(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Gt(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Gte(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Lt(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Lte(10)).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" = 10)
	// SELECT * FROM "test" WHERE ("a" != 10)
	// SELECT * FROM "test" WHERE ("a" > 10)
	// SELECT * FROM "test" WHERE ("a" >= 10)
	// SELECT * FROM "test" WHERE ("a" < 10)
	// SELECT * FROM "test" WHERE ("a" <= 10)
}

func ExampleC_inOperators() {
	// using identifiers
	sql, _, _ := pp.From("test").Where(pp.C("a").In("a", "b", "c")).Build()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = pp.From("test").Where(pp.C("a").In([]string{"a", "b", "c"})).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").NotIn("a", "b", "c")).Build()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = pp.From("test").Where(pp.C("a").NotIn([]string{"a", "b", "c"})).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c'))
}

func ExampleC_likeComparisons() {
	// using identifiers
	sql, _, _ := pp.From("test").Where(pp.C("a").Like("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").Like(regexp.MustCompile(`[ab]`))).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").ILike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").ILike(regexp.MustCompile("[ab]"))).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").NotLike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").NotLike(regexp.MustCompile("[ab]"))).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").NotILike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.C("a").NotILike(regexp.MustCompile(`[ab]`))).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" LIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" ~ '[ab]')
	// SELECT * FROM "test" WHERE ("a" ILIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" ~* '[ab]')
	// SELECT * FROM "test" WHERE ("a" NOT LIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" !~ '[ab]')
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE '%a%')
	// SELECT * FROM "test" WHERE ("a" !~* '[ab]')
}

func ExampleC_isComparisons() {
	sql, args, _ := pp.From("test").Where(pp.C("a").Is(nil)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").Is(true)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").Is(false)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNull()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsTrue()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsFalse()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNot(nil)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNot(true)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNot(false)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNotNull()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNotTrue()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.C("a").IsNotFalse()).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
}

func ExampleC_betweenComparisons() {
	ds := pp.From("test").Where(
		pp.C("a").Between(pp.Range(1, 10)),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(
		pp.C("a").NotBetween(pp.Range(1, 10)),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN ? AND ?) [1 10]
}

func ExampleCOALESCE() {
	ds := pp.From("test").Select(
		pp.COALESCE(pp.C("a"), "a"),
		pp.COALESCE(pp.C("a"), pp.C("b"), nil),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT COALESCE("a", 'a'), COALESCE("a", "b", NULL) FROM "test" []
	// SELECT COALESCE("a", ?), COALESCE("a", "b", ?) FROM "test" [a <nil>]
}

func ExampleCOALESCE_as() {
	sql, _, _ := pp.From("test").Select(pp.COALESCE(pp.C("a"), "a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT COALESCE("a", 'a') AS "a" FROM "test"
}

func ExampleCOUNT() {
	ds := pp.From("test").Select(pp.COUNT("*"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT COUNT(*) FROM "test" []
	// SELECT COUNT(*) FROM "test" []
}

func ExampleCOUNT_as() {
	sql, _, _ := pp.From("test").Select(pp.COUNT("*").As("count")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT COUNT(*) AS "count" FROM "test"
}

func ExampleCOUNT_havingClause() {
	ds := pp.
		From("test").
		Select(pp.COUNT("a").As("COUNT")).
		GroupBy("a").
		Having(pp.COUNT("a").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT COUNT("a") AS "COUNT" FROM "test" GROUP BY "a" HAVING (COUNT("a") > 10) []
	// SELECT COUNT("a") AS "COUNT" FROM "test" GROUP BY "a" HAVING (COUNT("a") > ?) [10]
}

func ExampleCast() {
	sql, _, _ := pp.From("test").
		Select(pp.Cast(pp.C("json1"), "TEXT").As("json_text")).
		Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.Cast(pp.C("json1"), "TEXT").Neq(
			pp.Cast(pp.C("json2"), "TEXT"),
		),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT CAST("json1" AS TEXT) AS "json_text" FROM "test"
	// SELECT * FROM "test" WHERE (CAST("json1" AS TEXT) != CAST("json2" AS TEXT))
}

func ExampleDISTINCT() {
	ds := pp.From("test").Select(pp.DISTINCT("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT DISTINCT("col") FROM "test" []
	// SELECT DISTINCT("col") FROM "test" []
}

func ExampleDISTINCT_as() {
	sql, _, _ := pp.From("test").Select(pp.DISTINCT("a").As("distinct_a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT DISTINCT("a") AS "distinct_a" FROM "test"
}

func ExampleDefault() {
	ds := pp.Insert("items")

	sql, args, _ := ds.Rows(pp.Record{
		"name":    pp.Default(),
		"address": pp.Default(),
	}).Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(pp.Record{
		"name":    pp.Default(),
		"address": pp.Default(),
	}).Build()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address", "name") VALUES (DEFAULT, DEFAULT) []
	// INSERT INTO "items" ("address", "name") VALUES (DEFAULT, DEFAULT) []
}

func ExampleDoNothing() {
	ds := pp.Insert("items")

	sql, args, _ := ds.Rows(pp.Record{
		"address": "111 Address",
		"name":    "bob",
	}).OnConflict(pp.DoNothing()).Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(pp.Record{
		"address": "111 Address",
		"name":    "bob",
	}).OnConflict(pp.DoNothing()).Build()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address", "name") VALUES ('111 Address', 'bob') ON CONFLICT DO NOTHING []
	// INSERT INTO "items" ("address", "name") VALUES (?, ?) ON CONFLICT DO NOTHING [111 Address bob]
}

func ExampleDoUpdate() {
	ds := pp.Insert("items")

	sql, args, _ := ds.
		Rows(pp.Record{"address": "111 Address"}).
		OnConflict(pp.DoUpdate("address", pp.C("address").Set(pp.I("excluded.address")))).
		Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).
		Rows(pp.Record{"address": "111 Address"}).
		OnConflict(pp.DoUpdate("address", pp.C("address").Set(pp.I("excluded.address")))).
		Build()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address") VALUES ('111 Address') ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" []
	// INSERT INTO "items" ("address") VALUES (?) ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" [111 Address]
}

func ExampleDoUpdate_where() {
	ds := pp.Insert("items")

	sql, args, _ := ds.
		Rows(pp.Record{"address": "111 Address"}).
		OnConflict(pp.DoUpdate(
			"address",
			pp.C("address").Set(pp.I("excluded.address"))).Where(pp.I("items.updated").IsNull()),
		).
		Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).
		Rows(pp.Record{"address": "111 Address"}).
		OnConflict(pp.DoUpdate(
			"address",
			pp.C("address").Set(pp.I("excluded.address"))).Where(pp.I("items.updated").IsNull()),
		).
		Build()
	fmt.Println(sql, args)

	// Output:
	// INSERT INTO "items" ("address") VALUES ('111 Address') ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" WHERE ("items"."updated" IS NULL) []
	// INSERT INTO "items" ("address") VALUES (?) ON CONFLICT (address) DO UPDATE SET "address"="excluded"."address" WHERE ("items"."updated" IS NULL) [111 Address]
}

func ExampleFIRST() {
	ds := pp.From("test").Select(pp.FIRST("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT FIRST("col") FROM "test" []
	// SELECT FIRST("col") FROM "test" []
}

func ExampleFIRST_as() {
	sql, _, _ := pp.From("test").Select(pp.FIRST("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT FIRST("a") AS "a" FROM "test"
}

// This example shows how to create custom SQL Functions
func ExampleFunc() {
	stragg := func(expression exp.Expression, delimiter string) exp.SQLFunctionExpression {
		return pp.Func("str_agg", expression, pp.L(delimiter))
	}
	sql, _, _ := pp.From("test").Select(stragg(pp.C("col"), "|")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT str_agg("col", |) FROM "test"
}

func ExampleI() {
	ds := pp.From("test").
		Select(
			pp.I("my_schema.table.col1"),
			pp.I("table.col2"),
			pp.I("col3"),
		)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Select(pp.I("test.*"))

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT "my_schema"."table"."col1", "table"."col2", "col3" FROM "test" []
	// SELECT "my_schema"."table"."col1", "table"."col2", "col3" FROM "test" []
	// SELECT "test".* FROM "test" []
	// SELECT "test".* FROM "test" []
}

func ExampleL() {
	ds := pp.From("test").Where(
		// literal with no args
		pp.L(`"col"::TEXT = ""other_col"::text`),
		// literal with args they will be interpolated into the sql by default
		pp.L("col IN (?, ?, ?)", "a", "b", "c"),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("col"::TEXT = ""other_col"::text AND col IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("col"::TEXT = ""other_col"::text AND col IN (?, ?, ?)) [a b c]
}

func ExampleL_withArgs() {
	ds := pp.From("test").Where(
		pp.L(
			"(? AND ?) OR ?",
			pp.C("a").Eq(1),
			pp.C("b").Eq("b"),
			pp.C("c").In([]string{"a", "b", "c"}),
		),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE (("a" = 1) AND ("b" = 'b')) OR ("c" IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE (("a" = ?) AND ("b" = ?)) OR ("c" IN (?, ?, ?)) [1 b a b c]
}

func ExampleL_as() {
	sql, _, _ := pp.From("test").Select(pp.L("json_col->>'totalAmount'").As("total_amount")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT json_col->>'totalAmount' AS "total_amount" FROM "test"
}

func ExampleL_comparisons() {
	// used from a literal expression
	sql, _, _ := pp.From("test").Where(pp.L("(a + b)").Eq(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a + b)").Neq(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a + b)").Gt(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a + b)").Gte(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a + b)").Lt(10)).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a + b)").Lte(10)).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ((a + b) = 10)
	// SELECT * FROM "test" WHERE ((a + b) != 10)
	// SELECT * FROM "test" WHERE ((a + b) > 10)
	// SELECT * FROM "test" WHERE ((a + b) >= 10)
	// SELECT * FROM "test" WHERE ((a + b) < 10)
	// SELECT * FROM "test" WHERE ((a + b) <= 10)
}

func ExampleL_inOperators() {
	// using identifiers
	sql, _, _ := pp.From("test").Where(pp.L("json_col->>'val'").In("a", "b", "c")).Build()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = pp.From("test").Where(pp.L("json_col->>'val'").In([]string{"a", "b", "c"})).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("json_col->>'val'").NotIn("a", "b", "c")).Build()
	fmt.Println(sql)
	// with a slice
	sql, _, _ = pp.From("test").Where(pp.L("json_col->>'val'").NotIn([]string{"a", "b", "c"})).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE (json_col->>'val' IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' NOT IN ('a', 'b', 'c'))
	// SELECT * FROM "test" WHERE (json_col->>'val' NOT IN ('a', 'b', 'c'))
}

func ExampleL_likeComparisons() {
	// using identifiers
	sql, _, _ := pp.From("test").Where(pp.L("(a::text || 'bar')").Like("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.L("(a::text || 'bar')").Like(regexp.MustCompile("[ab]")),
	).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a::text || 'bar')").ILike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.L("(a::text || 'bar')").ILike(regexp.MustCompile("[ab]")),
	).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a::text || 'bar')").NotLike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.L("(a::text || 'bar')").NotLike(regexp.MustCompile("[ab]")),
	).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(pp.L("(a::text || 'bar')").NotILike("%a%")).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("test").Where(
		pp.L("(a::text || 'bar')").NotILike(regexp.MustCompile("[ab]")),
	).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ((a::text || 'bar') LIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ~ '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ILIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') ~* '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') NOT LIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') !~ '[ab]')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') NOT ILIKE '%a%')
	// SELECT * FROM "test" WHERE ((a::text || 'bar') !~* '[ab]')
}

func ExampleL_isComparisons() {
	sql, args, _ := pp.From("test").Where(pp.L("a").Is(nil)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").Is(true)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").Is(false)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNull()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsTrue()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsFalse()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNot(nil)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNot(true)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNot(false)).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNotNull()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNotTrue()).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("test").Where(pp.L("a").IsNotFalse()).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (a IS NULL) []
	// SELECT * FROM "test" WHERE (a IS TRUE) []
	// SELECT * FROM "test" WHERE (a IS FALSE) []
	// SELECT * FROM "test" WHERE (a IS NULL) []
	// SELECT * FROM "test" WHERE (a IS TRUE) []
	// SELECT * FROM "test" WHERE (a IS FALSE) []
	// SELECT * FROM "test" WHERE (a IS NOT NULL) []
	// SELECT * FROM "test" WHERE (a IS NOT TRUE) []
	// SELECT * FROM "test" WHERE (a IS NOT FALSE) []
	// SELECT * FROM "test" WHERE (a IS NOT NULL) []
	// SELECT * FROM "test" WHERE (a IS NOT TRUE) []
	// SELECT * FROM "test" WHERE (a IS NOT FALSE) []
}

func ExampleL_betweenComparisons() {
	ds := pp.From("test").Where(
		pp.L("(a + b)").Between(pp.Range(1, 10)),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(
		pp.L("(a + b)").NotBetween(pp.Range(1, 10)),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ((a + b) BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ((a + b) BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ((a + b) NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ((a + b) NOT BETWEEN ? AND ?) [1 10]
}

func ExampleLAST() {
	ds := pp.From("test").Select(pp.LAST("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT LAST("col") FROM "test" []
	// SELECT LAST("col") FROM "test" []
}

func ExampleLAST_as() {
	sql, _, _ := pp.From("test").Select(pp.LAST("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT LAST("a") AS "a" FROM "test"
}

func ExampleMAX() {
	ds := pp.From("test").Select(pp.MAX("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT MAX("col") FROM "test" []
	// SELECT MAX("col") FROM "test" []
}

func ExampleMAX_as() {
	sql, _, _ := pp.From("test").Select(pp.MAX("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT MAX("a") AS "a" FROM "test"
}

func ExampleMAX_havingClause() {
	ds := pp.
		From("test").
		Select(pp.MAX("a").As("MAX")).
		GroupBy("a").
		Having(pp.MAX("a").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT MAX("a") AS "MAX" FROM "test" GROUP BY "a" HAVING (MAX("a") > 10) []
	// SELECT MAX("a") AS "MAX" FROM "test" GROUP BY "a" HAVING (MAX("a") > ?) [10]
}

func ExampleMIN() {
	ds := pp.From("test").Select(pp.MIN("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT MIN("col") FROM "test" []
	// SELECT MIN("col") FROM "test" []
}

func ExampleMIN_as() {
	sql, _, _ := pp.From("test").Select(pp.MIN("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT MIN("a") AS "a" FROM "test"
}

func ExampleMIN_havingClause() {
	ds := pp.
		From("test").
		Select(pp.MIN("a").As("MIN")).
		GroupBy("a").
		Having(pp.MIN("a").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT MIN("a") AS "MIN" FROM "test" GROUP BY "a" HAVING (MIN("a") > 10) []
	// SELECT MIN("a") AS "MIN" FROM "test" GROUP BY "a" HAVING (MIN("a") > ?) [10]
}

func ExampleOn() {
	ds := pp.From("test").Join(
		pp.T("my_table"),
		pp.On(pp.I("my_table.fkey").Eq(pp.I("other_table.id"))),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
}

func ExampleOn_withEx() {
	ds := pp.From("test").Join(
		pp.T("my_table"),
		pp.On(pp.Ex{"my_table.fkey": pp.I("other_table.id")}),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
	// SELECT * FROM "test" INNER JOIN "my_table" ON ("my_table"."fkey" = "other_table"."id") []
}

func ExampleOr() {
	ds := pp.From("test").Where(
		pp.Or(
			pp.C("col").Eq(10),
			pp.C("col").Eq(20),
		),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("col" = 10) OR ("col" = 20)) []
	// SELECT * FROM "test" WHERE (("col" = ?) OR ("col" = ?)) [10 20]
}

func ExampleOr_withAnd() {
	ds := pp.From("items").Where(
		pp.Or(
			pp.C("a").Gt(10),
			pp.And(
				pp.C("b").Eq(100),
				pp.C("c").Neq("test"),
			),
		),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE (("a" > 10) OR (("b" = 100) AND ("c" != 'test'))) []
	// SELECT * FROM "items" WHERE (("a" > ?) OR (("b" = ?) AND ("c" != ?))) [10 100 test]
}

func ExampleOr_withExMap() {
	ds := pp.From("test").Where(
		pp.Or(
			// Ex will be anded together
			pp.Ex{
				"col1": 1,
				"col2": true,
			},
			pp.Ex{
				"col3": nil,
				"col4": "foo",
			},
		),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ((("col1" = 1) AND ("col2" IS TRUE)) OR (("col3" IS NULL) AND ("col4" = 'foo'))) []
	// SELECT * FROM "test" WHERE ((("col1" = ?) AND ("col2" IS TRUE)) OR (("col3" IS NULL) AND ("col4" = ?))) [1 foo]
}

func ExampleRange_numbers() {
	ds := pp.From("test").Where(
		pp.C("col").Between(pp.Range(1, 10)),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(
		pp.C("col").NotBetween(pp.Range(1, 10)),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("col" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN ? AND ?) [1 10]
}

func ExampleRange_strings() {
	ds := pp.From("test").Where(
		pp.C("col").Between(pp.Range("a", "z")),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(
		pp.C("col").NotBetween(pp.Range("a", "z")),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col" BETWEEN 'a' AND 'z') []
	// SELECT * FROM "test" WHERE ("col" BETWEEN ? AND ?) [a z]
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN 'a' AND 'z') []
	// SELECT * FROM "test" WHERE ("col" NOT BETWEEN ? AND ?) [a z]
}

func ExampleRange_identifiers() {
	ds := pp.From("test").Where(
		pp.C("col1").Between(pp.Range(pp.C("col2"), pp.C("col3"))),
	)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(
		pp.C("col1").NotBetween(pp.Range(pp.C("col2"), pp.C("col3"))),
	)
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("col1" BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" NOT BETWEEN "col2" AND "col3") []
	// SELECT * FROM "test" WHERE ("col1" NOT BETWEEN "col2" AND "col3") []
}

func ExampleS() {
	s := pp.S("test_schema")
	t := s.Table("test")
	sql, args, _ := pp.
		From(t).
		Select(
			t.Col("col1"),
			t.Col("col2"),
			t.Col("col3"),
		).
		Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT "test_schema"."test"."col1", "test_schema"."test"."col2", "test_schema"."test"."col3" FROM "test_schema"."test" []
}

func ExampleSUM() {
	ds := pp.From("test").Select(pp.SUM("col"))
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT SUM("col") FROM "test" []
	// SELECT SUM("col") FROM "test" []
}

func ExampleSUM_as() {
	sql, _, _ := pp.From("test").Select(pp.SUM("a").As("a")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT SUM("a") AS "a" FROM "test"
}

func ExampleSUM_havingClause() {
	ds := pp.
		From("test").
		Select(pp.SUM("a").As("SUM")).
		GroupBy("a").
		Having(pp.SUM("a").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT SUM("a") AS "SUM" FROM "test" GROUP BY "a" HAVING (SUM("a") > 10) []
	// SELECT SUM("a") AS "SUM" FROM "test" GROUP BY "a" HAVING (SUM("a") > ?) [10]
}

func ExampleStar() {
	ds := pp.From("test").Select(pp.Star())

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
	// SELECT * FROM "test" []
}

func ExampleT() {
	t := pp.T("test")
	sql, args, _ := pp.
		From(t).
		Select(
			t.Col("col1"),
			t.Col("col2"),
			t.Col("col3"),
		).
		Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT "test"."col1", "test"."col2", "test"."col3" FROM "test" []
}

func ExampleUsing() {
	ds := pp.From("test").Join(
		pp.T("my_table"),
		pp.Using("fkey"),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
}

func ExampleUsing_withIdentifier() {
	ds := pp.From("test").Join(
		pp.T("my_table"),
		pp.Using(pp.C("fkey")),
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
	// SELECT * FROM "test" INNER JOIN "my_table" USING ("fkey") []
}

func ExampleEx() {
	ds := pp.From("items").Where(
		pp.Ex{
			"col1": "a",
			"col2": 1,
			"col3": true,
			"col4": false,
			"col5": nil,
			"col6": []string{"a", "b", "c"},
		},
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 'a') AND ("col2" = 1) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IS NULL) AND ("col6" IN ('a', 'b', 'c'))) []
	// SELECT * FROM "items" WHERE (("col1" = ?) AND ("col2" = ?) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IS NULL) AND ("col6" IN (?, ?, ?))) [a 1 a b c]
}

func ExampleEx_withOp() {
	sql, args, _ := pp.From("items").Where(
		pp.Ex{
			"col1": pp.Op{"neq": "a"},
			"col3": pp.Op{"isNot": true},
			"col6": pp.Op{"notIn": []string{"a", "b", "c"}},
		},
	).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE (("col1" != 'a') AND ("col3" IS NOT TRUE) AND ("col6" NOT IN ('a', 'b', 'c'))) []
}

func ExampleEx_in() {
	// using an Ex expression map
	sql, _, _ := pp.From("test").Where(pp.Ex{
		"a": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c'))
}

func ExampleExOr() {
	sql, args, _ := pp.From("items").Where(
		pp.ExOr{
			"col1": "a",
			"col2": 1,
			"col3": true,
			"col4": false,
			"col5": nil,
			"col6": []string{"a", "b", "c"},
		},
	).Build()
	fmt.Println(sql, args)

	// nolint:lll // sql statements are long
	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 'a') OR ("col2" = 1) OR ("col3" IS TRUE) OR ("col4" IS FALSE) OR ("col5" IS NULL) OR ("col6" IN ('a', 'b', 'c'))) []
}

func ExampleExOr_withOp() {
	sql, _, _ := pp.From("items").Where(pp.ExOr{
		"col1": pp.Op{"neq": "a"},
		"col3": pp.Op{"isNot": true},
		"col6": pp.Op{"notIn": []string{"a", "b", "c"}},
	}).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("items").Where(pp.ExOr{
		"col1": pp.Op{"gt": 1},
		"col2": pp.Op{"gte": 1},
		"col3": pp.Op{"lt": 1},
		"col4": pp.Op{"lte": 1},
	}).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("items").Where(pp.ExOr{
		"col1": pp.Op{"like": "a%"},
		"col2": pp.Op{"notLike": "a%"},
		"col3": pp.Op{"iLike": "a%"},
		"col4": pp.Op{"notILike": "a%"},
	}).Build()
	fmt.Println(sql)

	sql, _, _ = pp.From("items").Where(pp.ExOr{
		"col1": pp.Op{"like": regexp.MustCompile("^[ab]")},
		"col2": pp.Op{"notLike": regexp.MustCompile("^[ab]")},
		"col3": pp.Op{"iLike": regexp.MustCompile("^[ab]")},
		"col4": pp.Op{"notILike": regexp.MustCompile("^[ab]")},
	}).Build()
	fmt.Println(sql)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" != 'a') OR ("col3" IS NOT TRUE) OR ("col6" NOT IN ('a', 'b', 'c')))
	// SELECT * FROM "items" WHERE (("col1" > 1) OR ("col2" >= 1) OR ("col3" < 1) OR ("col4" <= 1))
	// SELECT * FROM "items" WHERE (("col1" LIKE 'a%') OR ("col2" NOT LIKE 'a%') OR ("col3" ILIKE 'a%') OR ("col4" NOT ILIKE 'a%'))
	// SELECT * FROM "items" WHERE (("col1" ~ '^[ab]') OR ("col2" !~ '^[ab]') OR ("col3" ~* '^[ab]') OR ("col4" !~* '^[ab]'))
}

func ExampleOp_comparisons() {
	ds := pp.From("test").Where(pp.Ex{
		"a": 10,
		"b": pp.Op{"neq": 10},
		"c": pp.Op{"gte": 10},
		"d": pp.Op{"lt": 10},
		"e": pp.Op{"lte": 10},
	})

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("a" = 10) AND ("b" != 10) AND ("c" >= 10) AND ("d" < 10) AND ("e" <= 10)) []
	// SELECT * FROM "test" WHERE (("a" = ?) AND ("b" != ?) AND ("c" >= ?) AND ("d" < ?) AND ("e" <= ?)) [10 10 10 10 10]
}

func ExampleOp_inComparisons() {
	// using an Ex expression map
	ds := pp.From("test").Where(pp.Ex{
		"a": pp.Op{"in": []string{"a", "b", "c"}},
	})

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notIn": []string{"a", "b", "c"}},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("a" IN (?, ?, ?)) [a b c]
	// SELECT * FROM "test" WHERE ("a" NOT IN ('a', 'b', 'c')) []
	// SELECT * FROM "test" WHERE ("a" NOT IN (?, ?, ?)) [a b c]
}

func ExampleOp_likeComparisons() {
	// using an Ex expression map
	ds := pp.From("test").Where(pp.Ex{
		"a": pp.Op{"like": "%a%"},
	})
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"like": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"iLike": "%a%"},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"iLike": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notLike": "%a%"},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notLike": regexp.MustCompile("[ab]")},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notILike": "%a%"},
	})

	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notILike": regexp.MustCompile("[ab]")},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" LIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" LIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" ~ '[ab]') []
	// SELECT * FROM "test" WHERE ("a" ~ ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" ILIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" ILIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" ~* '[ab]') []
	// SELECT * FROM "test" WHERE ("a" ~* ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" NOT LIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" NOT LIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" !~ '[ab]') []
	// SELECT * FROM "test" WHERE ("a" !~ ?) [[ab]]
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE '%a%') []
	// SELECT * FROM "test" WHERE ("a" NOT ILIKE ?) [%a%]
	// SELECT * FROM "test" WHERE ("a" !~* '[ab]') []
	// SELECT * FROM "test" WHERE ("a" !~* ?) [[ab]]
}

func ExampleOp_isComparisons() {
	// using an Ex expression map
	ds := pp.From("test").Where(pp.Ex{
		"a": true,
	})
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"is": true},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": false,
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"is": false},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": nil,
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"is": nil},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"isNot": true},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"isNot": false},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"isNot": nil},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)
	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT TRUE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT FALSE) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
	// SELECT * FROM "test" WHERE ("a" IS NOT NULL) []
}

func ExampleOp_betweenComparisons() {
	ds := pp.From("test").Where(pp.Ex{
		"a": pp.Op{"between": pp.Range(1, 10)},
	})
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("test").Where(pp.Ex{
		"a": pp.Op{"notBetween": pp.Range(1, 10)},
	})
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE ("a" BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" BETWEEN ? AND ?) [1 10]
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN 1 AND 10) []
	// SELECT * FROM "test" WHERE ("a" NOT BETWEEN ? AND ?) [1 10]
}

// When using a single op with multiple keys they are ORed together
func ExampleOp_withMultipleKeys() {
	ds := pp.From("items").Where(pp.Ex{
		"col1": pp.Op{"is": nil, "eq": 10},
	})

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "items" WHERE (("col1" = 10) OR ("col1" IS NULL)) []
	// SELECT * FROM "items" WHERE (("col1" = ?) OR ("col1" IS NULL)) [10]
}

func ExampleRecord_insert() {
	ds := pp.Insert("test")

	records := []pp.Record{
		{"col1": 1, "col2": "foo"},
		{"col1": 2, "col2": "bar"},
	}

	sql, args, _ := ds.Rows(records).Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Rows(records).Build()
	fmt.Println(sql, args)
	// Output:
	// INSERT INTO "test" ("col1", "col2") VALUES (1, 'foo'), (2, 'bar') []
	// INSERT INTO "test" ("col1", "col2") VALUES (?, ?), (?, ?) [1 foo 2 bar]
}

func ExampleRecord_update() {
	ds := pp.Update("test")
	update := pp.Record{"col1": 1, "col2": "foo"}

	sql, args, _ := ds.Set(update).Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Set(update).Build()
	fmt.Println(sql, args)
	// Output:
	// UPDATE "test" SET "col1"=1,"col2"='foo' []
	// UPDATE "test" SET "col1"=?,"col2"=? [1 foo]
}

func ExampleV() {
	ds := pp.From("user").Select(
		pp.V(true).As("is_verified"),
		pp.V(1.2).As("version"),
		"first_name",
		"last_name",
	)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	ds = pp.From("user").Where(pp.V(1).Neq(1))
	sql, args, _ = ds.Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT TRUE AS "is_verified", 1.2 AS "version", "first_name", "last_name" FROM "user" []
	// SELECT * FROM "user" WHERE (1 != 1) []
}

func ExampleV_prepared() {
	ds := pp.From("user").Select(
		pp.V(true).As("is_verified"),
		pp.V(1.2).As("version"),
		"first_name",
		"last_name",
	)

	sql, args, _ := ds.Prepared(true).Build()
	fmt.Println(sql, args)

	ds = pp.From("user").Where(pp.V(1).Neq(1))

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT ? AS "is_verified", ? AS "version", "first_name", "last_name" FROM "user" [true 1.2]
	// SELECT * FROM "user" WHERE (? != ?) [1 1]
}

func ExampleVals() {
	ds := pp.Insert("user").
		Cols("first_name", "last_name", "is_verified").
		Vals(
			pp.Vals{"Greg", "Farley", true},
			pp.Vals{"Jimmy", "Stewart", true},
			pp.Vals{"Jeff", "Jeffers", false},
		)
	insertSQL, args, _ := ds.Build()
	fmt.Println(insertSQL, args)

	// Output:
	// INSERT INTO "user" ("first_name", "last_name", "is_verified") VALUES ('Greg', 'Farley', TRUE), ('Jimmy', 'Stewart', TRUE), ('Jeff', 'Jeffers', FALSE) []
}

func ExampleW() {
	ds := pp.From("test").
		Select(pp.ROW_NUMBER().Over(pp.W().PartitionBy("a").OrderBy(pp.I("b").Asc())))
	query, args, _ := ds.Build()
	fmt.Println(query, args)

	ds = pp.From("test").
		Select(pp.ROW_NUMBER().OverName(pp.I("w"))).
		Window(pp.W("w").PartitionBy("a").OrderBy(pp.I("b").Asc()))
	query, args, _ = ds.Build()
	fmt.Println(query, args)

	ds = pp.From("test").
		Select(pp.ROW_NUMBER().OverName(pp.I("w1"))).
		Window(
			pp.W("w1").PartitionBy("a"),
			pp.W("w").Inherit("w1").OrderBy(pp.I("b").Asc()),
		)
	query, args, _ = ds.Build()
	fmt.Println(query, args)

	ds = pp.From("test").
		Select(pp.ROW_NUMBER().Over(pp.W().Inherit("w").OrderBy("b"))).
		Window(pp.W("w").PartitionBy("a"))
	query, args, _ = ds.Build()
	fmt.Println(query, args)
	// Output:
	// SELECT ROW_NUMBER() OVER (PARTITION BY "a" ORDER BY "b" ASC) FROM "test" []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w" AS (PARTITION BY "a" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER "w1" FROM "test" WINDOW "w1" AS (PARTITION BY "a"), "w" AS ("w1" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER ("w" ORDER BY "b") FROM "test" WINDOW "w" AS (PARTITION BY "a") []
}

func ExampleLateral() {
	maxEntry := pp.From("entry").
		Select(pp.MAX("int").As("max_int")).
		Where(pp.Ex{"time": pp.Op{"lt": pp.I("e.time")}}).
		As("max_entry")

	maxID := pp.From("entry").
		Select("id").
		Where(pp.Ex{"int": pp.I("max_entry.max_int")}).
		As("max_id")

	ds := pp.
		Select("e.id", "max_entry.max_int", "max_id.id").
		From(
			pp.T("entry").As("e"),
			pp.Lateral(maxEntry),
			pp.Lateral(maxID),
		)
	query, args, _ := ds.Build()
	fmt.Println(query, args)

	query, args, _ = ds.Prepared(true).Build()
	fmt.Println(query, args)

	// Output:
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e", LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry", LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" []
}

func ExampleLateral_join() {
	maxEntry := pp.From("entry").
		Select(pp.MAX("int").As("max_int")).
		Where(pp.Ex{"time": pp.Op{"lt": pp.I("e.time")}}).
		As("max_entry")

	maxID := pp.From("entry").
		Select("id").
		Where(pp.Ex{"int": pp.I("max_entry.max_int")}).
		As("max_id")

	ds := pp.
		Select("e.id", "max_entry.max_int", "max_id.id").
		From(pp.T("entry").As("e")).
		Join(pp.Lateral(maxEntry), pp.On(pp.V(true))).
		Join(pp.Lateral(maxID), pp.On(pp.V(true)))
	query, args, _ := ds.Build()
	fmt.Println(query, args)

	query, args, _ = ds.Prepared(true).Build()
	fmt.Println(query, args)

	// Output:
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON TRUE INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON TRUE []
	// SELECT "e"."id", "max_entry"."max_int", "max_id"."id" FROM "entry" AS "e" INNER JOIN LATERAL (SELECT MAX("int") AS "max_int" FROM "entry" WHERE ("time" < "e"."time")) AS "max_entry" ON ? INNER JOIN LATERAL (SELECT "id" FROM "entry" WHERE ("int" = "max_entry"."max_int")) AS "max_id" ON ? [true true]
}

func ExampleAny() {
	ds := pp.From("test").Where(pp.Ex{
		"id": pp.Any(pp.From("other").Select("test_id")),
	})
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("id" = ANY ((SELECT "test_id" FROM "other"))) []
	// SELECT * FROM "test" WHERE ("id" = ANY ((SELECT "test_id" FROM "other"))) []
}

func ExampleAll() {
	ds := pp.From("test").Where(pp.Ex{
		"id": pp.All(pp.From("other").Select("test_id")),
	})
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "test" WHERE ("id" = ALL ((SELECT "test_id" FROM "other"))) []
	// SELECT * FROM "test" WHERE ("id" = ALL ((SELECT "test_id" FROM "other"))) []
}

func ExampleCase_search() {
	ds := pp.From("test").
		Select(
			pp.C("col"),
			pp.Case().
				When(pp.C("col").Gt(0), true).
				When(pp.C("col").Lte(0), false).
				As("is_gt_zero"),
		)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE  WHEN ("col" > 0) THEN TRUE WHEN ("col" <= 0) THEN FALSE END AS "is_gt_zero" FROM "test" []
	// SELECT "col", CASE  WHEN ("col" > ?) THEN ? WHEN ("col" <= ?) THEN ? END AS "is_gt_zero" FROM "test" [0 true 0 false]
}

func ExampleCase_searchElse() {
	ds := pp.From("test").
		Select(
			pp.C("col"),
			pp.Case().
				When(pp.C("col").Gt(10), "Gt 10").
				When(pp.C("col").Gt(20), "Gt 20").
				Else("Bad Val").
				As("str_val"),
		)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE  WHEN ("col" > 10) THEN 'Gt 10' WHEN ("col" > 20) THEN 'Gt 20' ELSE 'Bad Val' END AS "str_val" FROM "test" []
	// SELECT "col", CASE  WHEN ("col" > ?) THEN ? WHEN ("col" > ?) THEN ? ELSE ? END AS "str_val" FROM "test" [10 Gt 10 20 Gt 20 Bad Val]
}

func ExampleCase_value() {
	ds := pp.From("test").
		Select(
			pp.C("col"),
			pp.Case().
				Value(pp.C("str")).
				When("foo", "FOO").
				When("bar", "BAR").
				As("foo_bar_upper"),
		)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE "str" WHEN 'foo' THEN 'FOO' WHEN 'bar' THEN 'BAR' END AS "foo_bar_upper" FROM "test" []
	// SELECT "col", CASE "str" WHEN ? THEN ? WHEN ? THEN ? END AS "foo_bar_upper" FROM "test" [foo FOO bar BAR]
}

func ExampleCase_valueElse() {
	ds := pp.From("test").
		Select(
			pp.C("col"),
			pp.Case().
				Value(pp.C("str")).
				When("foo", "FOO").
				When("bar", "BAR").
				Else("Baz").
				As("foo_bar_upper"),
		)
	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT "col", CASE "str" WHEN 'foo' THEN 'FOO' WHEN 'bar' THEN 'BAR' ELSE 'Baz' END AS "foo_bar_upper" FROM "test" []
	// SELECT "col", CASE "str" WHEN ? THEN ? WHEN ? THEN ? ELSE ? END AS "foo_bar_upper" FROM "test" [foo FOO bar BAR Baz]
}
