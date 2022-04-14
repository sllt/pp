package pp

import (
	"fmt"

	_ "manlu.org/pp/dialect/mysql"
)

func ExampleDelete() {
	ds := Delete("items")

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
}

func ExampleDeleteDataset_Executor() {
	db := getDB()

	de := db.Delete("user").
		Where(Ex{"first_name": "Bob"}).
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

	de := db.Delete("user").
		Where(C("last_name").Eq("Yukon")).
		Returning(C("id")).
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
	sql, _, _ := Delete("test").
		With("check_vals(val)", From().Select(L("123"))).
		Where(C("val").Eq(From("check_vals").Select("val"))).
		Build()
	fmt.Println(sql)

	// Output:
	// WITH check_vals(val) AS (SELECT 123) DELETE FROM "test" WHERE ("val" IN (SELECT "val" FROM "check_vals"))
}

func ExampleDeleteDataset_WithRecursive() {
	sql, _, _ := Delete("nums").
		WithRecursive("nums(x)",
			From().Select(L("1")).
				UnionAll(From("nums").
					Select(L("x+1")).Where(C("x").Lt(5)))).
		Build()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 UNION ALL (SELECT x+1 FROM "nums" WHERE ("x" < 5))) DELETE FROM "nums"
}

func ExampleDeleteDataset_Where() {
	// By default everything is anded together
	sql, _, _ := Delete("test").Where(Ex{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = Delete("test").Where(ExOr{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = Delete("test").Where(
		Or(
			Ex{
				"a": Op{"gt": 10},
				"b": Op{"lt": 10},
			},
			Ex{
				"c": nil,
				"d": []string{"a", "b", "c"},
			},
		),
	).Build()
	fmt.Println(sql)
	// By default everything is anded together
	sql, _, _ = Delete("test").Where(
		C("a").Gt(10),
		C("b").Lt(10),
		C("c").IsNull(),
		C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = Delete("test").Where(
		Or(
			C("a").Gt(10),
			And(
				C("b").Lt(10),
				C("c").IsNull(),
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
	sql, args, _ := Delete("test").Prepared(true).Where(Ex{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = Delete("test").Prepared(true).Where(ExOr{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = Delete("test").Prepared(true).Where(
		Or(
			Ex{
				"a": Op{"gt": 10},
				"b": Op{"lt": 10},
			},
			Ex{
				"c": nil,
				"d": []string{"a", "b", "c"},
			},
		),
	).Build()
	fmt.Println(sql, args)
	// By default everything is anded together
	sql, args, _ = Delete("test").Prepared(true).Where(
		C("a").Gt(10),
		C("b").Lt(10),
		C("c").IsNull(),
		C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = Delete("test").Prepared(true).Where(
		Or(
			C("a").Gt(10),
			And(
				C("b").Lt(10),
				C("c").IsNull(),
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
	ds := Delete("test").Where(
		Or(
			C("a").Gt(10),
			And(
				C("b").Lt(10),
				C("c").IsNull(),
			),
		),
	)
	sql, _, _ := ds.ClearWhere().Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test"
}

func ExampleDeleteDataset_Limit() {
	ds := Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT 10
}

func ExampleDeleteDataset_LimitAll() {
	// Using mysql dialect because it surts limit on delete
	ds := Dialect("mysql").Delete("test").LimitAll()
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` LIMIT ALL
}

func ExampleDeleteDataset_ClearLimit() {
	// Using mysql dialect because it surts limit on delete
	ds := Dialect("mysql").Delete("test").Limit(10)
	sql, _, _ := ds.ClearLimit().Build()
	fmt.Println(sql)
	// Output:
	// DELETE `test` FROM `test`
}

func ExampleDeleteDataset_Order() {
	// use mysql dialect because it surts order by on deletes
	ds := Dialect("mysql").Delete("test").Order(C("a").Asc())
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC
}

func ExampleDeleteDataset_OrderAnd() {
	// use mysql dialect because it surts order by on deletes
	ds := Dialect("mysql").Delete("test").Order(C("a").Asc())
	sql, _, _ := ds.Order(C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `a` ASC, `b` DESC NULLS LAST
}

func ExampleDeleteDataset_OrderPrepend() {
	// use mysql dialect because it surts order by on deletes
	ds := Dialect("mysql").Delete("test").Order(C("a").Asc())
	sql, _, _ := ds.OrderPrepend(C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM `test` ORDER BY `b` DESC NULLS LAST, `a` ASC
}

func ExampleDeleteDataset_ClearOrder() {
	ds := Delete("test").Order(C("a").Asc())
	sql, _, _ := ds.ClearOrder().Build()
	fmt.Println(sql)
	// Output:
	// DELETE FROM "test"
}

func ExampleDeleteDataset_Build() {
	sql, args, _ := Delete("items").Build()
	fmt.Println(sql, args)

	sql, args, _ = Delete("items").
		Where(Ex{"id": Op{"gt": 10}}).
		Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > 10) []
}

func ExampleDeleteDataset_Prepared() {
	sql, args, _ := Delete("items").Prepared(true).Build()
	fmt.Println(sql, args)

	sql, args, _ = Delete("items").
		Prepared(true).
		Where(Ex{"id": Op{"gt": 10}}).
		Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > ?) [10]
}

func ExampleDeleteDataset_Returning() {
	ds := Delete("items")
	sql, args, _ := ds.Returning("id").Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Returning("id").Where(C("id").IsNotNull()).Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" RETURNING "id" []
	// DELETE FROM "items" WHERE ("id" IS NOT NULL) RETURNING "id" []
}
