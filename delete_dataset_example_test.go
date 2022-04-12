package pp_test

import (
	"fmt"

	"manlu.org/pp"
	_ "manlu.org/pp/dialect/mysql"
)

func ExampleDelete() {
	ds := pp.Delete("items")

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
}

func ExampleDeleteDataset_Executor() {
	db := getDB()

	de := db.Delete("pp_user").
		Where(pp.Ex{"first_name": "Bob"}).
		Executor()
	if r, err := de.Exec(); err != nil {
		fmt.Println(err.Error())
	} else {
		c, _ := r.RowsAffected()
		fmt.Printf("Deleted %d users", c)
	}

	// Output:
	// Deleted 1 users
}

func ExampleDeleteDataset_Executor_returning() {
	db := getDB()

	de := db.Delete("pp_user").
		Where(pp.C("last_name").Eq("Yukon")).
		Returning(pp.C("id")).
		Executor()

	var ids []int64
	if err := de.ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Deleted users [ids:=%+v]", ids)
	}

	// Output:
	// Deleted users [ids:=[1 2 3]]
}

func ExampleDeleteDataset_With() {
	sql, _, _ := pp.Delete("test").
		With("check_vals(val)", pp.From().Select(pp.L("123"))).
		Where(pp.C("val").Eq(pp.From("check_vals").Select("val"))).
		Build()
	fmt.Println(sql)

	// Output:
	// WITH check_vals(val) AS (SELECT 123) DELETE FROM "test" WHERE ("val" IN (SELECT "val" FROM "check_vals"))
}

func ExampleDeleteDataset_WithRecursive() {
	sql, _, _ := pp.Delete("nums").
		WithRecursive("nums(x)",
			pp.From().Select(pp.L("1")).
				UnionAll(pp.From("nums").
					Select(pp.L("x+1")).Where(pp.C("x").Lt(5)))).
		Build()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 UNION ALL (SELECT x+1 FROM "nums" WHERE ("x" < 5))) DELETE FROM "nums"
}

func ExampleDeleteDataset_Where() {
	// By default everything is anded together
	sql, _, _ := pp.Delete("test").Where(pp.Ex{
		"a": pp.Op{"gt": 10},
		"b": pp.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = pp.Delete("test").Where(pp.ExOr{
		"a": pp.Op{"gt": 10},
		"b": pp.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = pp.Delete("test").Where(
		pp.Or(
			pp.Ex{
				"a": pp.Op{"gt": 10},
				"b": pp.Op{"lt": 10},
			},
			pp.Ex{
				"c": nil,
				"d": []string{"a", "b", "c"},
			},
		),
	).Build()
	fmt.Println(sql)
	// By default everything is anded together
	sql, _, _ = pp.Delete("test").Where(
		pp.C("a").Gt(10),
		pp.C("b").Lt(10),
		pp.C("c").IsNull(),
		pp.C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = pp.Delete("test").Where(
		pp.Or(
			pp.C("a").Gt(10),
			pp.And(
				pp.C("b").Lt(10),
				pp.C("c").IsNull(),
			),
		),
	).Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// DELETE FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// DELETE FROM "test" WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleDeleteDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := pp.Delete("test").Prepared(true).Where(pp.Ex{
		"a": pp.Op{"gt": 10},
		"b": pp.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = pp.Delete("test").Prepared(true).Where(pp.ExOr{
		"a": pp.Op{"gt": 10},
		"b": pp.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = pp.Delete("test").Prepared(true).Where(
		pp.Or(
			pp.Ex{
				"a": pp.Op{"gt": 10},
				"b": pp.Op{"lt": 10},
			},
			pp.Ex{
				"c": nil,
				"d": []string{"a", "b", "c"},
			},
		),
	).Build()
	fmt.Println(sql, args)
	// By default everything is anded together
	sql, args, _ = pp.Delete("test").Prepared(true).Where(
		pp.C("a").Gt(10),
		pp.C("b").Lt(10),
		pp.C("c").IsNull(),
		pp.C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = pp.Delete("test").Prepared(true).Where(
		pp.Or(
			pp.C("a").Gt(10),
			pp.And(
				pp.C("b").Lt(10),
				pp.C("c").IsNull(),
			),
		),
	).Build()
	fmt.Println(sql, args)
	// Output:
	// DELETE FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// DELETE FROM "test" WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [10 10]
}

func ExampleDeleteDataset_ClearWhere() {
	ds := pp.Delete("test").Where(
		pp.Or(
			pp.C("a").Gt(10),
			pp.And(
				pp.C("b").Lt(10),
				pp.C("c").IsNull(),
			),
		),
	)
	sql, _, _ := ds.ClearWhere().Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test"
}

func ExampleDeleteDataset_Limit() {
	ds := pp.Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT 10
}

func ExampleDeleteDataset_LimitAll() {
	// Using mysql dialect because it supports limit on delete
	ds := pp.Dialect("mysql").Delete("test").LimitAll()
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT ALL
}

func ExampleDeleteDataset_ClearLimit() {
	// Using mysql dialect because it supports limit on delete
	ds := pp.Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.ClearLimit().Build()
	fmt.Println(sql)
	// Output:
	// DELETE `test` FROM `test`
}

func ExampleDeleteDataset_Order() {
	// use mysql dialect because it supports order by on deletes
	ds := pp.Dialect("mysql").Delete("test").Order(pp.C("a").Asc())
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC
}

func ExampleDeleteDataset_OrderAppend() {
	// use mysql dialect because it supports order by on deletes
	ds := pp.Dialect("mysql").Delete("test").Order(pp.C("a").Asc())
	sql, _, _ := ds.OrderAppend(pp.C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC, `b` DESC NULLS LAST
}

func ExampleDeleteDataset_OrderPrepend() {
	// use mysql dialect because it supports order by on deletes
	ds := pp.Dialect("mysql").Delete("test").Order(pp.C("a").Asc())
	sql, _, _ := ds.OrderPrepend(pp.C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `b` DESC NULLS LAST, `a` ASC
}

func ExampleDeleteDataset_ClearOrder() {
	ds := pp.Delete("test").Order(pp.C("a").Asc())
	sql, _, _ := ds.ClearOrder().Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test"
}

func ExampleDeleteDataset_Build() {
	sql, args, _ := pp.Delete("items").Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.Delete("items").
		Where(pp.Ex{"id": pp.Op{"gt": 10}}).
		Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > 10) []
}

func ExampleDeleteDataset_Prepared() {
	sql, args, _ := pp.Delete("items").Prepared(true).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.Delete("items").
		Prepared(true).
		Where(pp.Ex{"id": pp.Op{"gt": 10}}).
		Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > ?) [10]
}

func ExampleDeleteDataset_Returning() {
	ds := pp.Delete("items")
	sql, args, _ := ds.Returning("id").Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Returning("id").Where(pp.C("id").IsNotNull()).Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" RETURNING "id" []
	// DELETE FROM "items" WHERE ("id" IS NOT NULL) RETURNING "id" []
}
