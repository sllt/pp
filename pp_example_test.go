package pp_test

import (
	"fmt"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"manlu.org/pp"
	_ "manlu.org/pp/dialect/mysql"
	_ "manlu.org/pp/dialect/postgres"
	_ "manlu.org/pp/dialect/sqlite3"
)

// Creating a mysql dataset. Be sure to import the mysql adapter.
func ExampleDialect_datasetMysql() {
	// import _ "manlu.org/pp/dialect/mysql"

	d := pp.Dialect("mysql")
	ds := d.From("test").Where(pp.Ex{
		"foo": "bar",
		"baz": []int64{1, 2, 3},
	}).Limit(10)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM `test` WHERE ((`baz` IN (1, 2, 3)) AND (`foo` = 'bar')) LIMIT 10 []
	// SELECT * FROM `test` WHERE ((`baz` IN (?, ?, ?)) AND (`foo` = ?)) LIMIT ? [1 2 3 bar 10]
}

// Creating a mysql database. Be sure to import the mysql adapter.
func ExampleDialect_dbMysql() {
	// import _ "manlu.org/pp/dialect/mysql"

	type item struct {
		ID      int64  `db:"id"`
		Address string `db:"address"`
		Name    string `db:"name"`
	}

	// set up a mock db this would normally be
	// db, err := sql.Open("mysql", dbURI)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	mDB, mock, _ := sqlmock.New()

	d := pp.Dialect("mysql")

	db := d.DB(mDB)

	// use the db.From to get a dataset to execute queries
	ds := db.From("items").Where(pp.C("id").Eq(1))

	// set up mock for example purposes
	mock.ExpectQuery("SELECT `address`, `id`, `name` FROM `items` WHERE \\(`id` = 1\\) LIMIT 1").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)
	var it item
	found, err := ds.ScanStruct(&it)
	fmt.Println(it, found, err)

	// set up mock for example purposes
	mock.ExpectQuery("SELECT `address`, `id`, `name` FROM `items` WHERE \\(`id` = \\?\\) LIMIT \\?").
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)

	found, err = ds.Prepared(true).ScanStruct(&it)
	fmt.Println(it, found, err)

	// Output:
	// {1 111 Test Addr Test1} true <nil>
	// {1 111 Test Addr Test1} true <nil>
}

// Creating a mysql dataset. Be sure to import the postgres adapter
func ExampleDialect_datasetPostgres() {
	// import _ "manlu.org/pp/dialect/postgres"

	d := pp.Dialect("postgres")
	ds := d.From("test").Where(pp.Ex{
		"foo": "bar",
		"baz": []int64{1, 2, 3},
	}).Limit(10)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM "test" WHERE (("baz" IN (1, 2, 3)) AND ("foo" = 'bar')) LIMIT 10 []
	// SELECT * FROM "test" WHERE (("baz" IN ($1, $2, $3)) AND ("foo" = $4)) LIMIT $5 [1 2 3 bar 10]
}

// Creating a postgres dataset. Be sure to import the postgres adapter
func ExampleDialect_dbPostgres() {
	// import _ "manlu.org/pp/dialect/postgres"

	type item struct {
		ID      int64  `db:"id"`
		Address string `db:"address"`
		Name    string `db:"name"`
	}

	// set up a mock db this would normally be
	// db, err := sql.Open("postgres", dbURI)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	mDB, mock, _ := sqlmock.New()

	d := pp.Dialect("postgres")

	db := d.DB(mDB)

	// use the db.From to get a dataset to execute queries
	ds := db.From("items").Where(pp.C("id").Eq(1))

	// set up mock for example purposes
	mock.ExpectQuery(`SELECT "address", "id", "name" FROM "items" WHERE \("id" = 1\) LIMIT 1`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)
	var it item
	found, err := ds.ScanStruct(&it)
	fmt.Println(it, found, err)

	// set up mock for example purposes
	mock.ExpectQuery(`SELECT "address", "id", "name" FROM "items" WHERE \("id" = \$1\) LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)

	found, err = ds.Prepared(true).ScanStruct(&it)
	fmt.Println(it, found, err)

	// Output:
	// {1 111 Test Addr Test1} true <nil>
	// {1 111 Test Addr Test1} true <nil>
}

// Creating a mysql dataset. Be sure to import the sqlite3 adapter
func ExampleDialect_datasetSqlite3() {
	// import _ "manlu.org/pp/dialect/sqlite3"

	d := pp.Dialect("sqlite3")
	ds := d.From("test").Where(pp.Ex{
		"foo": "bar",
		"baz": []int64{1, 2, 3},
	}).Limit(10)

	sql, args, _ := ds.Build()
	fmt.Println(sql, args)

	sql, args, _ = ds.Prepared(true).Build()
	fmt.Println(sql, args)

	// Output:
	// SELECT * FROM `test` WHERE ((`baz` IN (1, 2, 3)) AND (`foo` = 'bar')) LIMIT 10 []
	// SELECT * FROM `test` WHERE ((`baz` IN (?, ?, ?)) AND (`foo` = ?)) LIMIT ? [1 2 3 bar 10]
}

// Creating a sqlite3 database. Be sure to import the sqlite3 adapter
func ExampleDialect_dbSqlite3() {
	// import _ "manlu.org/pp/dialect/sqlite3"
	type item struct {
		ID      int64  `db:"id"`
		Address string `db:"address"`
		Name    string `db:"name"`
	}

	// set up a mock db this would normally be
	// db, err := sql.Open("sqlite3", dbURI)
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	mDB, mock, _ := sqlmock.New()

	d := pp.Dialect("sqlite3")

	db := d.DB(mDB)

	// use the db.From to get a dataset to execute queries
	ds := db.From("items").Where(pp.C("id").Eq(1))

	// set up mock for example purposes
	mock.ExpectQuery("SELECT `address`, `id`, `name` FROM `items` WHERE \\(`id` = 1\\) LIMIT 1").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)
	var it item
	found, err := ds.ScanStruct(&it)
	fmt.Println(it, found, err)

	// set up mock for example purposes
	mock.ExpectQuery("SELECT `address`, `id`, `name` FROM `items` WHERE \\(`id` = \\?\\) LIMIT \\?").
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "address", "name"}).
				FromCSVString("1, 111 Test Addr,Test1"),
		)

	found, err = ds.Prepared(true).ScanStruct(&it)
	fmt.Println(it, found, err)

	// Output:
	// {1 111 Test Addr Test1} true <nil>
	// {1 111 Test Addr Test1} true <nil>
}

func ExampleSetTimeLocation() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	created, err := time.Parse(time.RFC3339, "2019-10-01T15:01:00Z")
	if err != nil {
		panic(err)
	}

	// use original time with tz info
	pp.SetTimeLocation(loc)
	ds := pp.Insert("test").Rows(pp.Record{
		"address": "111 Address",
		"name":    "Bob Yukon",
		"created": created,
	})
	sql, _, _ := ds.Build()
	fmt.Println(sql)

	// convert time to UTC
	pp.SetTimeLocation(time.UTC)
	sql, _, _ = ds.Build()
	fmt.Println(sql)

	// Output:
	// INSERT INTO "test" ("address", "created", "name") VALUES ('111 Address', '2019-10-01T23:01:00+08:00', 'Bob Yukon')
	// INSERT INTO "test" ("address", "created", "name") VALUES ('111 Address', '2019-10-01T15:01:00Z', 'Bob Yukon')
}