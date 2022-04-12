// nolint:lll // sql statements are long
package pp

import (
	goSQL "database/sql"
	"fmt"
	"github.com/lib/pq"
	"manlu.org/pp/exp"
	"os"
	"regexp"
	"time"
)

const schema = `
		DROP TABLE IF EXISTS "user_role";
		DROP TABLE IF EXISTS "pp_user";	
		CREATE  TABLE "pp_user" (
			"id" SERIAL PRIMARY KEY NOT NULL,
			"first_name" VARCHAR(45) NOT NULL,
			"last_name" VARCHAR(45) NOT NULL,
			"created" TIMESTAMP NOT NULL DEFAULT now()
		);
		CREATE  TABLE "user_role" (
			"id" SERIAL PRIMARY KEY NOT NULL,
			"user_id" BIGINT NOT NULL REFERENCES pp_user(id) ON DELETE CASCADE,
			"name" VARCHAR(45) NOT NULL,
			"created" TIMESTAMP NOT NULL DEFAULT now()
		); 
    `

const defaultDBURI = "postgres://postgres:@localhost:5435/pppostgres?sslmode=disable"

var ppDB *Database

func getDB() *Database {
	if ppDB == nil {
		dbURI := os.Getenv("PG_URI")
		if dbURI == "" {
			dbURI = defaultDBURI
		}
		uri, err := pq.ParseURL(dbURI)
		if err != nil {
			panic(err)
		}
		pdb, err := goSQL.Open("postgres", uri)
		if err != nil {
			panic(err)
		}
		ppDB = New("postgres", pdb)
	}
	// reset the db
	if _, err := ppDB.Exec(schema); err != nil {
		panic(err)
	}
	type ppUser struct {
		ID        int64     `db:"id" pp:"skipinsert"`
		FirstName string    `db:"first_name"`
		LastName  string    `db:"last_name"`
		Created   time.Time `db:"created" pp:"skipupdate"`
	}

	users := []ppUser{
		{FirstName: "Bob", LastName: "Yukon"},
		{FirstName: "Sally", LastName: "Yukon"},
		{FirstName: "Vinita", LastName: "Yukon"},
		{FirstName: "John", LastName: "Doe"},
	}
	var userIds []int64
	err := ppDB.Insert("pp_user").Rows(users).Returning("id").Executor().ScanVals(&userIds)
	if err != nil {
		panic(err)
	}
	type userRole struct {
		ID      int64     `db:"id" pp:"skipinsert"`
		UserID  int64     `db:"user_id"`
		Name    string    `db:"name"`
		Created time.Time `db:"created" pp:"skipupdate"`
	}

	roles := []userRole{
		{UserID: userIds[0], Name: "Admin"},
		{UserID: userIds[1], Name: "Manager"},
		{UserID: userIds[2], Name: "Manager"},
		{UserID: userIds[3], Name: "User"},
	}
	_, err = ppDB.Insert("user_role").Rows(roles).Executor().Exec()
	if err != nil {
		panic(err)
	}
	return ppDB
}

func ExampleSelectDataset() {
	ds := From("test").
		Select(COUNT("*")).
		InnerJoin(T("test2"), On(I("test.fkey").Eq(I("test2.id")))).
		LeftJoin(T("test3"), On(I("test2.fkey").Eq(I("test3.id")))).
		Where(
			Ex{
				"test.name": Op{
					"like": regexp.MustCompile("^[ab]"),
				},
				"test2.amount": Op{
					"isNot": nil,
				},
			},
			ExOr{
				"test3.id":     nil,
				"test3.status": []string{"passed", "active", "registered"},
			}).
		Order(I("test.created").Desc().NullsLast()).
		GroupBy(I("test.user_id")).
		Having(AVG("test3.age").Gt(10))

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// nolint:lll // SQL statements are long
	// Output:
	// SELECT COUNT(*) FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."id") LEFT JOIN "test3" ON ("test2"."fkey" = "test3"."id") WHERE ((("test"."name" ~ '^[ab]') AND ("test2"."amount" IS NOT NULL)) AND (("test3"."id" IS NULL) OR ("test3"."status" IN ('passed', 'active', 'registered')))) GROUP BY "test"."user_id" HAVING (AVG("test3"."age") > 10) ORDER BY "test"."created" DESC NULLS LAST []
	// SELECT COUNT(*) FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."id") LEFT JOIN "test3" ON ("test2"."fkey" = "test3"."id") WHERE ((("test"."name" ~ ?) AND ("test2"."amount" IS NOT NULL)) AND (("test3"."id" IS NULL) OR ("test3"."status" IN (?, ?, ?)))) GROUP BY "test"."user_id" HAVING (AVG("test3"."age") > ?) ORDER BY "test"."created" DESC NULLS LAST [^[ab] passed active registered 10]
}

func ExampleSelect() {
	sql, _, _ := Select(L("NOW()")).Build()
	fmt.Println(sql)

	// Output:
	// SELECT NOW()
}

func ExampleFrom() {
	sql, args, _ := From("test").Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" []
}

func ExampleSelectDataset_As() {
	ds := From("test").As("t")
	sql, _, _ := From(ds).Build()
	fmt.Println(sql)
	// Output: SELECT * FROM (SELECT * FROM "test") AS "t"
}

func ExampleSelectDataset_Union() {
	sql, _, _ := From("test").
		Union(From("test2")).
		Build()
	fmt.Println(sql)

	sql, _, _ = From("test").
		Limit(1).
		Union(From("test2")).
		Build()
	fmt.Println(sql)

	sql, _, _ = From("test").
		Limit(1).
		Union(From("test2").
			Order(C("id").Desc())).
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" UNION (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_UnionAll() {
	sql, _, _ := From("test").
		UnionAll(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		UnionAll(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		UnionAll(From("test2").
			Order(C("id").Desc())).
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" UNION ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" UNION ALL (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_With() {
	sql, _, _ := From("one").
		With("one", From().Select(L("1"))).
		Select(Star()).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("derived").
		With("intermed", From("test").Select(Star()).Where(C("x").Gte(5))).
		With("derived", From("intermed").Select(Star()).Where(C("x").Lt(10))).
		Select(Star()).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("multi").
		With("multi(x,y)", From().Select(L("1"), L("2"))).
		Select(C("x"), C("y")).
		Build()
	fmt.Println(sql)

	// Output:
	// WITH one AS (SELECT 1) SELECT * FROM "one"
	// WITH intermed AS (SELECT * FROM "test" WHERE ("x" >= 5)), derived AS (SELECT * FROM "intermed" WHERE ("x" < 10)) SELECT * FROM "derived"
	// WITH multi(x,y) AS (SELECT 1, 2) SELECT "x", "y" FROM "multi"
}

func ExampleSelectDataset_With_insertDataset() {
	insertDs := Insert("foo").Rows(Record{"user_id": 10}).Returning("id")

	ds := From("bar").
		With("ins", insertDs).
		Select("bar_name").
		Where(Ex{"bar.user_id": I("ins.user_id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (10) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id")
	// WITH ins AS (INSERT INTO "foo" ("user_id") VALUES (?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "ins"."user_id") [10]
}

func ExampleSelectDataset_With_updateDataset() {
	updateDs := Update("foo").Set(Record{"bar": "baz"}).Returning("id")

	ds := From("bar").
		With("upd", updateDs).
		Select("bar_name").
		Where(Ex{"bar.user_id": I("upd.user_id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// WITH upd AS (UPDATE "foo" SET "bar"='baz' RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id")
	// WITH upd AS (UPDATE "foo" SET "bar"=? RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "upd"."user_id") [baz]
}

func ExampleSelectDataset_With_deleteDataset() {
	deleteDs := Delete("foo").Where(Ex{"bar": "baz"}).Returning("id")

	ds := From("bar").
		With("del", deleteDs).
		Select("bar_name").
		Where(Ex{"bar.user_id": I("del.user_id")})

	sql, _, _ := ds.Build()
	fmt.Println(sql)

	sql, args, _ := ds.Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// WITH del AS (DELETE FROM "foo" WHERE ("bar" = 'baz') RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id")
	// WITH del AS (DELETE FROM "foo" WHERE ("bar" = ?) RETURNING "id") SELECT "bar_name" FROM "bar" WHERE ("bar"."user_id" = "del"."user_id") [baz]
}

func ExampleSelectDataset_WithRecursive() {
	sql, _, _ := From("nums").
		WithRecursive("nums(x)",
			From().Select(L("1")).
				UnionAll(From("nums").
					Select(L("x+1")).Where(C("x").Lt(5)))).
		Build()
	fmt.Println(sql)
	// Output:
	// WITH RECURSIVE nums(x) AS (SELECT 1 UNION ALL (SELECT x+1 FROM "nums" WHERE ("x" < 5))) SELECT * FROM "nums"
}

func ExampleSelectDataset_Intersect() {
	sql, _, _ := From("test").
		Intersect(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		Intersect(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		Intersect(From("test2").
			Order(C("id").Desc())).
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INTERSECT (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_IntersectAll() {
	sql, _, _ := From("test").
		IntersectAll(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		IntersectAll(From("test2")).
		Build()
	fmt.Println(sql)
	sql, _, _ = From("test").
		Limit(1).
		IntersectAll(From("test2").
			Order(C("id").Desc())).
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INTERSECT ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT ALL (SELECT * FROM "test2")
	// SELECT * FROM (SELECT * FROM "test" LIMIT 1) AS "t1" INTERSECT ALL (SELECT * FROM (SELECT * FROM "test2" ORDER BY "id" DESC) AS "t1")
}

func ExampleSelectDataset_ClearOffset() {
	ds := From("test").
		Offset(2)
	sql, _, _ := ds.
		ClearOffset().
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Offset() {
	ds := From("test").Offset(2)
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" OFFSET 2
}

func ExampleSelectDataset_Limit() {
	ds := From("test").Limit(10)
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LIMIT 10
}

func ExampleSelectDataset_LimitAll() {
	ds := From("test").LimitAll()
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LIMIT ALL
}

func ExampleSelectDataset_ClearLimit() {
	ds := From("test").Limit(10)
	sql, _, _ := ds.ClearLimit().Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Order() {
	ds := From("test").Order(C("a").Asc())
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC
}

func ExampleSelectDataset_Order_caseExpression() {
	ds := From("test").Order(Case().When(C("num").Gt(10), 0).Else(1).Asc())
	sql, _, _ := ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY CASE  WHEN ("num" > 10) THEN 0 ELSE 1 END ASC
}

func ExampleSelectDataset_OrderAppend() {
	ds := From("test").Order(C("a").Asc())
	sql, _, _ := ds.OrderAppend(C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "a" ASC, "b" DESC NULLS LAST
}

func ExampleSelectDataset_OrderPrepend() {
	ds := From("test").Order(C("a").Asc())
	sql, _, _ := ds.OrderPrepend(C("b").Desc().NullsLast()).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" ORDER BY "b" DESC NULLS LAST, "a" ASC
}

func ExampleSelectDataset_ClearOrder() {
	ds := From("test").Order(C("a").Asc())
	sql, _, _ := ds.ClearOrder().Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
}

func ExampleSelectDataset_GroupBy() {
	sql, _, _ := From("test").
		Select(SUM("income").As("income_sum")).
		GroupBy("age").
		Build()
	fmt.Println(sql)
	// Output:
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age"
}

func ExampleSelectDataset_GroupByAppend() {
	ds := From("test").
		Select(SUM("income").As("income_sum")).
		GroupBy("age")
	sql, _, _ := ds.
		GroupByAppend("job").
		Build()
	fmt.Println(sql)
	// the original dataset group by does not change
	sql, _, _ = ds.Build()
	fmt.Println(sql)
	// Output:
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age", "job"
	// SELECT SUM("income") AS "income_sum" FROM "test" GROUP BY "age"
}

func ExampleSelectDataset_Having() {
	sql, _, _ := From("test").Having(SUM("income").Gt(1000)).Build()
	fmt.Println(sql)
	sql, _, _ = From("test").GroupBy("age").Having(SUM("income").Gt(1000)).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" HAVING (SUM("income") > 1000)
	// SELECT * FROM "test" GROUP BY "age" HAVING (SUM("income") > 1000)
}

func ExampleSelectDataset_Window() {
	ds := From("test").
		Select(ROW_NUMBER().Over(W().PartitionBy("a").OrderBy(I("b").Asc())))
	query, args, _ := ds.Build()
	fmt.Println(query, args)

	ds = From("test").
		Select(ROW_NUMBER().OverName(I("w"))).
		Window(W("w").PartitionBy("a").OrderBy(I("b").Asc()))
	query, args, _ = ds.Build()
	fmt.Println(query, args)

	ds = From("test").
		Select(ROW_NUMBER().OverName(I("w1"))).
		Window(
			W("w1").PartitionBy("a"),
			W("w").Inherit("w1").OrderBy(I("b").Asc()),
		)
	query, args, _ = ds.Build()
	fmt.Println(query, args)

	ds = From("test").
		Select(ROW_NUMBER().Over(W().Inherit("w").OrderBy("b"))).
		Window(W("w").PartitionBy("a"))
	query, args, _ = ds.Build()
	fmt.Println(query, args)
	// Output
	// SELECT ROW_NUMBER() OVER (PARTITION BY "a" ORDER BY "b" ASC) FROM "test" []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w" AS (PARTITION BY "a" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER "w" FROM "test" WINDOW "w1" AS (PARTITION BY "a"), "w" AS ("w1" ORDER BY "b" ASC) []
	// SELECT ROW_NUMBER() OVER ("w" ORDER BY "b") FROM "test" WINDOW "w" AS (PARTITION BY "a") []
}

func ExampleSelectDataset_Where() {
	// By default everything is anded together
	sql, _, _ := From("test").Where(Ex{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use ExOr to get ORed expressions together
	sql, _, _ = From("test").Where(ExOr{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, _, _ = From("test").Where(
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
	sql, _, _ = From("test").Where(
		C("a").Gt(10),
		C("b").Lt(10),
		C("c").IsNull(),
		C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql)
	// You can use a combination of Ors and Ands
	sql, _, _ = From("test").Where(
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
	// SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE (("a" > 10) OR ("b" < 10) OR ("c" IS NULL) OR ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE ((("a" > 10) AND ("b" < 10)) OR (("c" IS NULL) AND ("d" IN ('a', 'b', 'c'))))
	// SELECT * FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
	// SELECT * FROM "test" WHERE (("a" > 10) OR (("b" < 10) AND ("c" IS NULL)))
}

func ExampleSelectDataset_Where_prepared() {
	// By default everything is anded together
	sql, args, _ := From("test").Prepared(true).Where(Ex{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use ExOr to get ORed expressions together
	sql, args, _ = From("test").Prepared(true).Where(ExOr{
		"a": Op{"gt": 10},
		"b": Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// You can use Or with Ex to Or multiple Ex maps together
	sql, args, _ = From("test").Prepared(true).Where(
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
	sql, args, _ = From("test").Prepared(true).Where(
		C("a").Gt(10),
		C("b").Lt(10),
		C("c").IsNull(),
		C("d").In("a", "b", "c"),
	).Build()
	fmt.Println(sql, args)
	// You can use a combination of Ors and Ands
	sql, args, _ = From("test").Prepared(true).Where(
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
	// SELECT * FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) OR ("b" < ?) OR ("c" IS NULL) OR ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE ((("a" > ?) AND ("b" < ?)) OR (("c" IS NULL) AND ("d" IN (?, ?, ?)))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
	// SELECT * FROM "test" WHERE (("a" > ?) OR (("b" < ?) AND ("c" IS NULL))) [10 10]
}

func ExampleSelectDataset_ClearWhere() {
	ds := From("test").Where(
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
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Join() {
	sql, _, _ := From("test").Join(
		T("test2"),
		On(Ex{"test.fkey": I("test2.Id")}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").Join(T("test2"), Using("common_column")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").Join(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(T("test2").Col("Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").Join(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(T("test").Col("fkey").Eq(T("t").Col("Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_InnerJoin() {
	sql, _, _ := From("test").InnerJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").InnerJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").InnerJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").InnerJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" INNER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" INNER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_FullOuterJoin() {
	sql, _, _ := From("test").FullOuterJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullOuterJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullOuterJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullOuterJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" FULL OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" FULL OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_RightOuterJoin() {
	sql, _, _ := From("test").RightOuterJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightOuterJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightOuterJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightOuterJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" RIGHT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" RIGHT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_LeftOuterJoin() {
	sql, _, _ := From("test").LeftOuterJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftOuterJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftOuterJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftOuterJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LEFT OUTER JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT OUTER JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" LEFT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT OUTER JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_FullJoin() {
	sql, _, _ := From("test").FullJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").FullJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" FULL JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_RightJoin() {
	sql, _, _ := From("test").RightJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").RightJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" RIGHT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_LeftJoin() {
	sql, _, _ := From("test").LeftJoin(
		T("test2"),
		On(Ex{
			"test.fkey": I("test2.Id"),
		}),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftJoin(
		T("test2"),
		Using("common_column"),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftJoin(
		From("test2").Where(C("amount").Gt(0)),
		On(I("test.fkey").Eq(I("test2.Id"))),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").LeftJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
		On(I("test.fkey").Eq(I("t.Id"))),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" LEFT JOIN "test2" ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT JOIN "test2" USING ("common_column")
	// SELECT * FROM "test" LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) ON ("test"."fkey" = "test2"."Id")
	// SELECT * FROM "test" LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t" ON ("test"."fkey" = "t"."Id")
}

func ExampleSelectDataset_NaturalJoin() {
	sql, _, _ := From("test").NaturalJoin(T("test2")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalJoin(
		From("test2").Where(C("amount").Gt(0)),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL JOIN "test2"
	// SELECT * FROM "test" NATURAL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalLeftJoin() {
	sql, _, _ := From("test").NaturalLeftJoin(T("test2")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalLeftJoin(
		From("test2").Where(C("amount").Gt(0)),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalLeftJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL LEFT JOIN "test2"
	// SELECT * FROM "test" NATURAL LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL LEFT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalRightJoin() {
	sql, _, _ := From("test").NaturalRightJoin(T("test2")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalRightJoin(
		From("test2").Where(C("amount").Gt(0)),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalRightJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL RIGHT JOIN "test2"
	// SELECT * FROM "test" NATURAL RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL RIGHT JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_NaturalFullJoin() {
	sql, _, _ := From("test").NaturalFullJoin(T("test2")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalFullJoin(
		From("test2").Where(C("amount").Gt(0)),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").NaturalFullJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" NATURAL FULL JOIN "test2"
	// SELECT * FROM "test" NATURAL FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" NATURAL FULL JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_CrossJoin() {
	sql, _, _ := From("test").CrossJoin(T("test2")).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").CrossJoin(
		From("test2").Where(C("amount").Gt(0)),
	).Build()
	fmt.Println(sql)

	sql, _, _ = From("test").CrossJoin(
		From("test2").Where(C("amount").Gt(0)).As("t"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test" CROSS JOIN "test2"
	// SELECT * FROM "test" CROSS JOIN (SELECT * FROM "test2" WHERE ("amount" > 0))
	// SELECT * FROM "test" CROSS JOIN (SELECT * FROM "test2" WHERE ("amount" > 0)) AS "t"
}

func ExampleSelectDataset_FromSelf() {
	sql, _, _ := From("test").FromSelf().Build()
	fmt.Println(sql)
	sql, _, _ = From("test").As("my_test_table").FromSelf().Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test") AS "t1"
	// SELECT * FROM (SELECT * FROM "test") AS "my_test_table"
}

func ExampleSelectDataset_From() {
	ds := From("test")
	sql, _, _ := ds.From("test2").Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test2"
}

func ExampleSelectDataset_From_withDataset() {
	ds := From("test")
	fromDs := ds.Where(C("age").Gt(10))
	sql, _, _ := ds.From(fromDs).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "t1"
}

func ExampleSelectDataset_From_withAliasedDataset() {
	ds := From("test")
	fromDs := ds.Where(C("age").Gt(10))
	sql, _, _ := ds.From(fromDs.As("test2")).Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM (SELECT * FROM "test" WHERE ("age" > 10)) AS "test2"
}

func ExampleSelectDataset_Select() {
	sql, _, _ := From("test").Select("a", "b", "c").Build()
	fmt.Println(sql)
	// Output:
	// SELECT "a", "b", "c" FROM "test"
}

func ExampleSelectDataset_Select_withDataset() {
	ds := From("test")
	fromDs := ds.Select("age").Where(C("age").Gt(10))
	sql, _, _ := ds.From().Select(fromDs).Build()
	fmt.Println(sql)
	// Output:
	// SELECT (SELECT "age" FROM "test" WHERE ("age" > 10))
}

func ExampleSelectDataset_Select_withAliasedDataset() {
	ds := From("test")
	fromDs := ds.Select("age").Where(C("age").Gt(10))
	sql, _, _ := ds.From().Select(fromDs.As("ages")).Build()
	fmt.Println(sql)
	// Output:
	// SELECT (SELECT "age" FROM "test" WHERE ("age" > 10)) AS "ages"
}

func ExampleSelectDataset_Select_withLiteral() {
	sql, _, _ := From("test").Select(L("a + b").As("sum")).Build()
	fmt.Println(sql)
	// Output:
	// SELECT a + b AS "sum" FROM "test"
}

func ExampleSelectDataset_Select_withSQLFunctionExpression() {
	sql, _, _ := From("test").Select(
		COUNT("*").As("age_count"),
		MAX("age").As("max_age"),
		AVG("age").As("avg_age"),
	).Build()
	fmt.Println(sql)
	// Output:
	// SELECT COUNT(*) AS "age_count", MAX("age") AS "max_age", AVG("age") AS "avg_age" FROM "test"
}

func ExampleSelectDataset_Select_withStruct() {
	ds := From("test")

	type myStruct struct {
		Name         string
		Address      string `db:"address"`
		EmailAddress string `db:"email_address"`
	}

	// Pass with pointer
	sql, _, _ := ds.Select(&myStruct{}).Build()
	fmt.Println(sql)

	// Pass instance of
	sql, _, _ = ds.Select(myStruct{}).Build()
	fmt.Println(sql)

	type myStruct2 struct {
		myStruct
		Zipcode string `db:"zipcode"`
	}

	// Pass pointer to struct with embedded struct
	sql, _, _ = ds.Select(&myStruct2{}).Build()
	fmt.Println(sql)

	// Pass instance of struct with embedded struct
	sql, _, _ = ds.Select(myStruct2{}).Build()
	fmt.Println(sql)

	var myStructs []myStruct

	// Pass slice of structs, will only select columns from underlying type
	sql, _, _ = ds.Select(myStructs).Build()
	fmt.Println(sql)

	// Output:
	// SELECT "address", "email_address", "name" FROM "test"
	// SELECT "address", "email_address", "name" FROM "test"
	// SELECT "address", "email_address", "name", "zipcode" FROM "test"
	// SELECT "address", "email_address", "name", "zipcode" FROM "test"
	// SELECT "address", "email_address", "name" FROM "test"
}

func ExampleSelectDataset_Distinct() {
	sql, _, _ := From("test").Select("a", "b").Distinct().Build()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT "a", "b" FROM "test"
}

func ExampleSelectDataset_Distinct_on() {
	sql, _, _ := From("test").Distinct("a").Build()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON ("a") * FROM "test"
}

func ExampleSelectDataset_Distinct_onWithLiteral() {
	sql, _, _ := From("test").Distinct(L("COALESCE(?, ?)", C("a"), "empty")).Build()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
}

func ExampleSelectDataset_Distinct_onCoalesce() {
	sql, _, _ := From("test").Distinct(COALESCE(C("a"), "empty")).Build()
	fmt.Println(sql)
	// Output:
	// SELECT DISTINCT ON (COALESCE("a", 'empty')) * FROM "test"
}

func ExampleSelectDataset_SelectAppend() {
	ds := From("test").Select("a", "b")
	sql, _, _ := ds.SelectAppend("c").Build()
	fmt.Println(sql)
	ds = From("test").Select("a", "b").Distinct()
	sql, _, _ = ds.SelectAppend("c").Build()
	fmt.Println(sql)
	// Output:
	// SELECT "a", "b", "c" FROM "test"
	// SELECT DISTINCT "a", "b", "c" FROM "test"
}

func ExampleSelectDataset_ClearSelect() {
	ds := From("test").Select("a", "b")
	sql, _, _ := ds.ClearSelect().Build()
	fmt.Println(sql)
	ds = From("test").Select("a", "b").Distinct()
	sql, _, _ = ds.ClearSelect().Build()
	fmt.Println(sql)
	// Output:
	// SELECT * FROM "test"
	// SELECT * FROM "test"
}

func ExampleSelectDataset_Build() {
	sql, args, _ := From("items").Where(Ex{"a": 1}).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE ("a" = 1) []
}

func ExampleSelectDataset_Build_prepared() {
	sql, args, _ := From("items").Where(Ex{"a": 1}).Prepared(true).Build()
	fmt.Println(sql, args)
	// Output:
	// SELECT * FROM "items" WHERE ("a" = ?) [1]
}

func ExampleSelectDataset_Update() {
	type item struct {
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := From("items").Update().Set(
		item{Name: "Test", Address: "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").Update().Set(
		Record{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").Update().Set(
		map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	// Output:
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
	// UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
}

func ExampleSelectDataset_Insert() {
	type item struct {
		ID      uint32 `db:"id" pp:"skipinsert"`
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	sql, args, _ := From("items").Insert().Rows(
		item{Name: "Test1", Address: "111 Test Addr"},
		item{Name: "Test2", Address: "112 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").Insert().Rows(
		Record{"name": "Test1", "address": "111 Test Addr"},
		Record{"name": "Test2", "address": "112 Test Addr"},
	).Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").Insert().Rows(
		[]item{
			{Name: "Test1", Address: "111 Test Addr"},
			{Name: "Test2", Address: "112 Test Addr"},
		}).Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").Insert().Rows(
		[]Record{
			{"name": "Test1", "address": "111 Test Addr"},
			{"name": "Test2", "address": "112 Test Addr"},
		}).Build()
	fmt.Println(sql, args)
	// Output:
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
	// INSERT INTO "items" ("address", "name") VALUES ('111 Test Addr', 'Test1'), ('112 Test Addr', 'Test2') []
}

func ExampleSelectDataset_Delete() {
	sql, args, _ := From("items").Delete().Build()
	fmt.Println(sql, args)

	sql, args, _ = From("items").
		Where(Ex{"id": Op{"gt": 10}}).
		Delete().
		Build()
	fmt.Println(sql, args)

	// Output:
	// DELETE FROM "items" []
	// DELETE FROM "items" WHERE ("id" > 10) []
}

func ExampleSelectDataset_Truncate() {
	sql, args, _ := From("items").Truncate().Build()
	fmt.Println(sql, args)
	// Output:
	// TRUNCATE "items" []
}

func ExampleSelectDataset_Prepared() {
	sql, args, _ := From("items").Prepared(true).Where(Ex{
		"col1": "a",
		"col2": 1,
		"col3": true,
		"col4": false,
		"col5": []string{"a", "b", "c"},
	}).Build()
	fmt.Println(sql, args)
	// nolint:lll // sql statements are long
	// Output:
	// SELECT * FROM "items" WHERE (("col1" = ?) AND ("col2" = ?) AND ("col3" IS TRUE) AND ("col4" IS FALSE) AND ("col5" IN (?, ?, ?))) [a 1 a b c]
}

func ExampleSelectDataset_ScanStructs() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()
	var users []User
	if err := db.From("pp_user").ScanStructs(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	users = users[0:0]
	if err := db.From("pp_user").Select("first_name").ScanStructs(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	// Output:
	// [{FirstName:Bob LastName:Yukon} {FirstName:Sally LastName:Yukon} {FirstName:Vinita LastName:Yukon} {FirstName:John LastName:Doe}]
	// [{FirstName:Bob LastName:} {FirstName:Sally LastName:} {FirstName:Vinita LastName:} {FirstName:John LastName:}]
}

func ExampleSelectDataset_ScanStructs_prepared() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()

	ds := db.From("pp_user").
		Prepared(true).
		Where(Ex{
			"last_name": "Yukon",
		})

	var users []User
	if err := ds.ScanStructs(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("\n%+v", users)

	// Output:
	// [{FirstName:Bob LastName:Yukon} {FirstName:Sally LastName:Yukon} {FirstName:Vinita LastName:Yukon}]
}

// In this example we create a new struct that has two structs that represent two table
// the User and Role fields are tagged with the table name
func ExampleSelectDataset_ScanStructs_withJoinAutoSelect() {
	type Role struct {
		UserID uint64 `db:"user_id"`
		Name   string `db:"name"`
	}
	type User struct {
		ID        uint64 `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	type UserAndRole struct {
		User User `db:"pp_user"`   // tag as the "pp_user" table
		Role Role `db:"user_role"` // tag as "user_role" table
	}
	db := getDB()

	ds := db.
		From("pp_user").
		Join(T("user_role"), On(I("pp_user.id").Eq(I("user_role.user_id"))))
	var users []UserAndRole
	// Scan structs will auto build the
	if err := ds.ScanStructs(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, u := range users {
		fmt.Printf("\n%+v", u)
	}
	// Output:
	// {User:{ID:1 FirstName:Bob LastName:Yukon} Role:{UserID:1 Name:Admin}}
	// {User:{ID:2 FirstName:Sally LastName:Yukon} Role:{UserID:2 Name:Manager}}
	// {User:{ID:3 FirstName:Vinita LastName:Yukon} Role:{UserID:3 Name:Manager}}
	// {User:{ID:4 FirstName:John LastName:Doe} Role:{UserID:4 Name:User}}
}

// In this example we create a new struct that has the user properties as well as a nested
// Role struct from the join table
func ExampleSelectDataset_ScanStructs_withJoinManualSelect() {
	type Role struct {
		UserID uint64 `db:"user_id"`
		Name   string `db:"name"`
	}
	type User struct {
		ID        uint64 `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Role      Role   `db:"user_role"` // tag as "user_role" table
	}
	db := getDB()

	ds := db.
		Select(
			"pp_user.id",
			"pp_user.first_name",
			"pp_user.last_name",
			// alias the fully qualified identifier `C` is important here so it doesnt parse it
			I("user_role.user_id").As(C("user_role.user_id")),
			I("user_role.name").As(C("user_role.name")),
		).
		From("pp_user").
		Join(T("user_role"), On(I("pp_user.id").Eq(I("user_role.user_id"))))
	var users []User
	if err := ds.ScanStructs(&users); err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, u := range users {
		fmt.Printf("\n%+v", u)
	}

	// Output:
	// {ID:1 FirstName:Bob LastName:Yukon Role:{UserID:1 Name:Admin}}
	// {ID:2 FirstName:Sally LastName:Yukon Role:{UserID:2 Name:Manager}}
	// {ID:3 FirstName:Vinita LastName:Yukon Role:{UserID:3 Name:Manager}}
	// {ID:4 FirstName:John LastName:Doe Role:{UserID:4 Name:User}}
}

func ExampleSelectDataset_ScanStruct() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()
	findUserByName := func(name string) {
		var user User
		ds := db.From("pp_user").Where(C("first_name").Eq(name))
		found, err := ds.ScanStruct(&user)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user: %+v\n", user)
		}
	}

	findUserByName("Bob")
	findUserByName("Zeb")

	// Output:
	// Found user: {FirstName:Bob LastName:Yukon}
	// No user found for first_name Zeb
}

// In this example we create a new struct that has two structs that represent two table
// the User and Role fields are tagged with the table name
func ExampleSelectDataset_ScanStruct_withJoinAutoSelect() {
	type Role struct {
		UserID uint64 `db:"user_id"`
		Name   string `db:"name"`
	}
	type User struct {
		ID        uint64 `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	type UserAndRole struct {
		User User `db:"pp_user"`   // tag as the "pp_user" table
		Role Role `db:"user_role"` // tag as "user_role" table
	}
	db := getDB()
	findUserAndRoleByName := func(name string) {
		var userAndRole UserAndRole
		ds := db.
			From("pp_user").
			Join(
				T("user_role"),
				On(I("pp_user.id").Eq(I("user_role.user_id"))),
			).
			Where(C("first_name").Eq(name))
		found, err := ds.ScanStruct(&userAndRole)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user and role: %+v\n", userAndRole)
		}
	}

	findUserAndRoleByName("Bob")
	findUserAndRoleByName("Zeb")
	// Output:
	// Found user and role: {User:{ID:1 FirstName:Bob LastName:Yukon} Role:{UserID:1 Name:Admin}}
	// No user found for first_name Zeb
}

// In this example we create a new struct that has the user properties as well as a nested
// Role struct from the join table
func ExampleSelectDataset_ScanStruct_withJoinManualSelect() {
	type Role struct {
		UserID uint64 `db:"user_id"`
		Name   string `db:"name"`
	}
	type User struct {
		ID        uint64 `db:"id"`
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
		Role      Role   `db:"user_role"` // tag as "user_role" table
	}
	db := getDB()
	findUserByName := func(name string) {
		var userAndRole User
		ds := db.
			Select(
				"pp_user.id",
				"pp_user.first_name",
				"pp_user.last_name",
				// alias the fully qualified identifier `C` is important here so it doesnt parse it
				I("user_role.user_id").As(C("user_role.user_id")),
				I("user_role.name").As(C("user_role.name")),
			).
			From("pp_user").
			Join(
				T("user_role"),
				On(I("pp_user.id").Eq(I("user_role.user_id"))),
			).
			Where(C("first_name").Eq(name))
		found, err := ds.ScanStruct(&userAndRole)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No user found for first_name %s\n", name)
		default:
			fmt.Printf("Found user and role: %+v\n", userAndRole)
		}
	}

	findUserByName("Bob")
	findUserByName("Zeb")

	// Output:
	// Found user and role: {ID:1 FirstName:Bob LastName:Yukon Role:{UserID:1 Name:Admin}}
	// No user found for first_name Zeb
}

func ExampleSelectDataset_ScanVals() {
	var ids []int64
	if err := getDB().From("pp_user").Select("id").ScanVals(&ids); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("UserIds = %+v", ids)

	// Output:
	// UserIds = [1 2 3 4]
}

func ExampleSelectDataset_ScanVal() {
	db := getDB()
	findUserIDByName := func(name string) {
		var id int64
		ds := db.From("pp_user").
			Select("id").
			Where(C("first_name").Eq(name))

		found, err := ds.ScanVal(&id)
		switch {
		case err != nil:
			fmt.Println(err.Error())
		case !found:
			fmt.Printf("No id found for user %s", name)
		default:
			fmt.Printf("\nFound userId: %+v\n", id)
		}
	}

	findUserIDByName("Bob")
	findUserIDByName("Zeb")
	// Output:
	// Found userId: 1
	// No id found for user Zeb
}

func ExampleSelectDataset_Count() {
	count, err := getDB().From("pp_user").Count()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Count is %d", count)

	// Output:
	// Count is 4
}

func ExampleSelectDataset_Pluck() {
	var lastNames []string
	if err := getDB().From("pp_user").Pluck(&lastNames, "last_name"); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("LastNames = %+v", lastNames)

	// Output:
	// LastNames = [Yukon Yukon Yukon Doe]
}

func ExampleSelectDataset_Executor_scannerScanStruct() {
	type User struct {
		FirstName string `db:"first_name"`
		LastName  string `db:"last_name"`
	}
	db := getDB()

	scanner, err := db.
		From("pp_user").
		Select("first_name", "last_name").
		Where(Ex{
			"last_name": "Yukon",
		}).
		Executor().
		Scanner()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer scanner.Close()

	for scanner.Next() {
		u := User{}

		err = scanner.ScanStruct(&u)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\n%+v", u)
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}

	// Output:
	// {FirstName:Bob LastName:Yukon}
	// {FirstName:Sally LastName:Yukon}
	// {FirstName:Vinita LastName:Yukon}
}

func ExampleSelectDataset_Executor_scannerScanVal() {
	db := getDB()

	scanner, err := db.
		From("pp_user").
		Select("first_name").
		Where(Ex{
			"last_name": "Yukon",
		}).
		Executor().
		Scanner()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer scanner.Close()

	for scanner.Next() {
		name := ""

		err = scanner.ScanVal(&name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println(name)
	}

	if scanner.Err() != nil {
		fmt.Println(scanner.Err().Error())
	}

	// Output:
	// Bob
	// Sally
	// Vinita
}

func ExampleForUpdate() {
	sql, args, _ := From("test").ForUpdate(exp.Wait).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" FOR UPDATE  []
}

func ExampleForUpdate_of() {
	sql, args, _ := From("test").ForUpdate(exp.Wait, T("test")).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" FOR UPDATE OF "test"  []
}

func ExampleForUpdate_ofMultiple() {
	sql, args, _ := From("table1").Join(
		T("table2"),
		On(I("table2.id").Eq(I("table1.id"))),
	).ForUpdate(
		exp.Wait,
		T("table1"),
		T("table2"),
	).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "table1" INNER JOIN "table2" ON ("table2"."id" = "table1"."id") FOR UPDATE OF "table1", "table2"  []
}
