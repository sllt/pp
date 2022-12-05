// nolint:lll // sql statements are long
package pp_test

import (
	"fmt"
	"github.com/sllt/pp"

	_ "github.com/sllt/pp/dialect/mysql"
)

func ExampleUpdate_withStruct() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withGoquRecord() {
	sql, args, _ := pp.Update("items").Set(
		pp.Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withMap() {
	sql, args, _ := pp.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdate_withSkipUpdateTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" pp:"skipupdate"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr' []
}

func ExampleUpdateDataset_Executor() {
	db := getDB()
	update := db.Update("pp_user").
		Where(pp.C("first_name").Eq("Bob")).
		Set(pp.Record{"first_name": "Bobby"}).
		Executor()

	if r, err := update.Exec(); err != nil {
		fmt.Println(err.Error())
	} else {
		c, _ := r.RowsAffected()
		fmt.Printf("Updated %d users", c)
	}

	// Output:
	// Updated 1 users
}

func ExampleUpdateDataset_Executor_returning() {
	db := getDB()
	var ids []int64
	update := db.Update("pp_user").
		Set(pp.Record{"last_name": "ucon"}).
		Where(pp.Ex{"last_name": "Yukon"}).
		Returning("id").
		Executor()
	if err := update.ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Updated users with ids %+v", ids)
	}

	// Output:
	// Updated users with ids [1 2 3]
}

func ExampleUpdateDataset_Returning() {
	sql, _, _ := pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Returning("id").
		Build()
	fmt.Println(sql)
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Returning(pp.T("test").All()).
		Build()
	fmt.Println(sql)
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Returning("a", "b").
		Build()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" SET "foo"='bar' RETURNING "id"
	// UPDATE "test" SET "foo"='bar' RETURNING "test".*
	// UPDATE "test" SET "foo"='bar' RETURNING "a", "b"
}

func ExampleUpdateDataset_With() {
	sql, _, _ := pp.Update("test").
		With("some_vals(val)", pp.From().Select(pp.L("123"))).
		Where(pp.C("val").Eq(pp.From("some_vals").Select("val"))).
		Set(pp.Record{"name": "Test"}).Build()
	fmt.Println(sql)

	// Output:
	// WITH some_vals(val) AS (SELECT 123) UPDATE "test" SET "name"='Test' WHERE ("val" IN (SELECT "val" FROM "some_vals"))
}

func ExampleUpdateDataset_WithRecursive() {
	sql, _, _ := pp.Update("nums").
		WithRecursive("nums(x)", pp.From().Select(pp.L("1").As("num")).
			UnionAll(pp.From("nums").
				Select(pp.L("x+1").As("num")).Where(pp.C("x").Lt(5)))).
		Set(pp.Record{"foo": pp.T("nums").Col("num")}).
		Build()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 AS "num" UNION ALL (SELECT x+1 AS "num" FROM "nums" WHERE ("x" < 5))) UPDATE "nums" SET "foo"="nums"."num"
}

func ExampleUpdateDataset_Limit() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Limit(10)
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' LIMIT 10
}

func ExampleUpdateDataset_LimitAll() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		LimitAll()
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' LIMIT ALL
}

func ExampleUpdateDataset_ClearLimit() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Limit(10)
	sql, _, _ := ds.ClearLimit().Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar'
}

func ExampleUpdateDataset_Order() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Order(pp.C("a").Asc())
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC
}

func ExampleUpdateDataset_OrderAppend() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Order(pp.C("a").Asc())
	sql, _, _ := ds.OrderAppend(pp.C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC, `b` DESC NULLS LAST
}

func ExampleUpdateDataset_OrderPrepend() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Order(pp.C("a").Asc())

	sql, _, _ := ds.OrderPrepend(pp.C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar' ORDER BY `b` DESC NULLS LAST, `a` ASC
}

func ExampleUpdateDataset_ClearOrder() {
	ds := pp.Dialect("mysql").
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Order(pp.C("a").Asc())
	sql, _, _ := ds.ClearOrder().Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `test` SET `foo`='bar'
}

func ExampleUpdateDataset_From() {
	ds := pp.Update("table_one").
		Set(pp.Record{"foo": pp.I("table_two.bar")}).
		From("table_two").
		Where(pp.Ex{"table_one.id": pp.I("table_two.id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
}

func ExampleUpdateDataset_From_postgres() {
	dialect := pp.Dialect("postgres")

	ds := dialect.Update("table_one").
		Set(pp.Record{"foo": pp.I("table_two.bar")}).
		From("table_two").
		Where(pp.Ex{"table_one.id": pp.I("table_two.id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
}

func ExampleUpdateDataset_From_mysql() {
	dialect := pp.Dialect("mysql")

	ds := dialect.Update("table_one").
		Set(pp.Record{"foo": pp.I("table_two.bar")}).
		From("table_two").
		Where(pp.Ex{"table_one.id": pp.I("table_two.id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// UPDATE `table_one`,`table_two` SET `foo`=`table_two`.`bar` WHERE (`table_one`.`id` = `table_two`.`id`)
}

func ExampleUpdateDataset_Where() {
	// By default everything is anded together
	sql, _, _ := pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(pp.Ex{
			"a": pp.Op{"gt": 10},
			"b": pp.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).Build()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(pp.ExOr{
			"a": pp.Op{"gt": 10},
			"b": pp.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).Build()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(
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
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(
			pp.C("a").Gt(10),
			pp.C("b").Lt(10),
			pp.C("c").IsNull(),
			pp.C("d").In("a", "b", "c"),
		).Build()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = pp.Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(
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
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleUpdateDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := pp.Update("test").
		Prepared(true).
		Set(pp.Record{"foo": "bar"}).
		Where(pp.Ex{
			"a": pp.Op{"gt": 10},
			"b": pp.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).Build()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = pp.Update("test").Prepared(true).
		Set(pp.Record{"foo": "bar"}).
		Where(pp.ExOr{
			"a": pp.Op{"gt": 10},
			"b": pp.Op{"lt": 10},
			"c": nil,
			"d": []string{"a", "b", "c"},
		}).Build()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = pp.Update("test").Prepared(true).
		Set(pp.Record{"foo": "bar"}).
		Where(
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
	sql, args, _ = pp.Update("test").Prepared(true).
		Set(pp.Record{"foo": "bar"}).
		Where(
			pp.C("a").Gt(10),
			pp.C("b").Lt(10),
			pp.C("c").IsNull(),
			pp.C("d").In("a", "b", "c"),
		).Build()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = pp.Update("test").Prepared(true).
		Set(pp.Record{"foo": "bar"}).
		Where(
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
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [bar 10 10 a b c]
	// UPDATE "test" SET "foo"=? WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [bar 10 10]
}

func ExampleUpdateDataset_ClearWhere() {
	ds := pp.
		Update("test").
		Set(pp.Record{"foo": "bar"}).
		Where(
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
	// UPDATE "test" SET "foo"='bar'
}

func ExampleUpdateDataset_Table() {
	ds := pp.Update("test")
	sql, _, _ := ds.Table("test2").Set(pp.Record{"foo": "bar"}).Build()
	fmt.Println(sql)
	// Output:
	// UPDATE "test2" SET "foo"='bar'
}

func ExampleUpdateDataset_Table_aliased() {
	ds := pp.Update("test")
	sql, _, _ := ds.Table(pp.T("test").As("t")).Set(pp.Record{"foo": "bar"}).Build()
	fmt.Println(sql)
	// Output:
	// UPDATE "test" AS "t" SET "foo"='bar'
}

func ExampleUpdateDataset_Set() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.Update("items").Set(
		pp.Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_struct() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_ppRecord() {
	sql, args, _ := pp.Update("items").Set(
		pp.Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_map() {
	sql, args, _ := pp.Update("items").Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_withSkipUpdateTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" pp:"skipupdate"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr' []
}

func ExampleUpdateDataset_Set_withDefaultIfEmptyTag() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name" pp:"defaultifempty"`
	}
	sql, args, _ := pp.Update("items").Set(
		item{Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.Update("items").Set(
		item{Name: "Bob Yukon", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"=DEFAULT []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Bob Yukon' []
}

func ExampleUpdateDataset_Set_withNoTags() {
	type item struct {
		Address string
		Name    string
	}
	sql, args, _ := pp.Update("items").Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleUpdateDataset_Set_withEmbeddedStruct() {
	type Address struct {
		Street string `db:"address_street"`
		State  string `db:"address_state"`
	}
	type User struct {
		Address
		FirstName string
		LastName  string
	}
	ds := pp.Update("user").Set(
		User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
	)
	updateSQL, args, _ := ds.Build()
	fmt.Println(updateSQL, args)

	// Output:
	// UPDATE "user" SET "address_state"='NY',"address_street"='111 Street',"firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_Set_withIgnoredEmbedded() {
	type Address struct {
		Street string
		State  string
	}
	type User struct {
		Address   `db:"-"`
		FirstName string
		LastName  string
	}
	ds := pp.Update("user").Set(
		User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
	)
	updateSQL, args, _ := ds.Build()
	fmt.Println(updateSQL, args)

	// Output:
	// UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_Set_withNilEmbeddedPointer() {
	type Address struct {
		Street string
		State  string
	}
	type User struct {
		*Address
		FirstName string
		LastName  string
	}
	ds := pp.Update("user").Set(
		User{FirstName: "Greg", LastName: "Farley"},
	)
	updateSQL, args, _ := ds.Build()
	fmt.Println(updateSQL, args)

	// Output:
	// UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
}

func ExampleUpdateDataset_Build_prepared() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}

	sql, args, _ := pp.From("items").Prepared(true).Update().Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("items").Prepared(true).Update().Set(
		pp.Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = pp.From("items").Prepared(true).Update().Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)
	// Output:
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
}

func ExampleUpdateDataset_Prepared() {
	sql, args, _ := pp.Update("items").Prepared(true).Set(
		pp.Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"=?,"name"=? [111 Test Addr Test]
}
