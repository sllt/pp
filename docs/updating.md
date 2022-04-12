# Updating

* [Create a UpdateDataset](#create)
* Examples
  * [Set with `pp.Record`](#set-record)
  * [Set with struct](#set-struct)
  * [Set with map](#set-map)
  * [Multi Table](#from)
  * [Where](#where)
  * [Order](#order)
  * [Limit](#limit)
  * [Returning](#returning)
  * [SetError](#seterror)
  * [Executing](#executing)

<a name="create"></a>
To create a [`UpdateDataset`](#UpdateDataset)  you can use

**[`pp.Update`](#Update)**

When you just want to create some quick SQL, this mostly follows the `Postgres` with the exception of placeholders for prepared statements.

```go
ds := pp.Update("user").Set(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley'
```

**[`SelectDataset.Update`](#SelectDataset.Update)**

If you already have a `SelectDataset` you can invoke `Update()` to get a `UpdateDataset`

**NOTE** This method will also copy over the `WITH`, `WHERE`, `ORDER`, and `LIMIT` clauses from the update

```go
ds := pp.From("user")

updateSQL, _, _ := ds.Update().Set(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
).Build()
fmt.Println(insertSQL, args)

updateSQL, _, _ = ds.Where(pp.C("first_name").Eq("Gregory")).Update().Set(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
).Build()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley'
UPDATE "user" SET "first_name"='Greg', "last_name"='Farley' WHERE "first_name"='Gregory'
```

**[`DialectWrapper.Update`](#DialectWrapper.Update)**

Use this when you want to create SQL for a specific `dialect`

```go
// import _ "manlu.org/pp/dialect/mysql"

dialect := pp.Dialect("mysql")

ds := dialect.Update("user").Set(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE `user` SET `first_name`='Greg', `last_name`='Farley'
```

**[`Database.Update`](#DialectWrapper.Update)**

Use this when you want to execute the SQL or create SQL for the drivers dialect.

```go
// import _ "manlu.org/pp/dialect/mysql"

mysqlDB := //initialize your db
db := pp.New("mysql", mysqlDB)

ds := db.Update("user").Set(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
updateSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
UPDATE `user` SET `first_name`='Greg', `last_name`='Farley'
```

### Examples

For more examples visit the **[Docs](#UpdateDataset)**

<a name="set-record"></a>
**[Set with `pp.Record`](#UpdateDataset.Set)**

```go
sql, args, _ := pp.Update("items").Set(
	pp.Record{"name": "Test", "address": "111 Test Addr"},
).Build()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

<a name="set-struct"></a>
**[Set with Struct](#UpdateDataset.Set)**

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name"`
}
sql, args, _ := pp.Update("items").Set(
	item{Name: "Test", Address: "111 Test Addr"},
).Build()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

With structs you can also skip fields by using the `skipupdate` tag

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name" pp:"skipupdate"`
}
sql, args, _ := pp.Update("items").Set(
	item{Name: "Test", Address: "111 Test Addr"},
).Build()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr' []
```

If you want to use the database `DEFAULT` when the struct field is a zero value you can use the `defaultifempty` tag.

```go
type item struct {
	Address string `db:"address"`
	Name    string `db:"name" pp:"defaultifempty"`
}
sql, args, _ := pp.Update("items").Set(
	item{Address: "111 Test Addr"},
).Build()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"=DEFAULT []
```

`pp` will also use fields in embedded structs when creating an update.

**NOTE** unexported fields will be ignored!

```go
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
```

Output:
```
UPDATE "user" SET "address_state"='NY',"address_street"='111 Street',"firstname"='Greg',"lastname"='Farley' []
```

**NOTE** When working with embedded pointers if the embedded struct is nil then the fields will be ignored.

```go
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
```

Output:
```
UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
```

You can ignore an embedded struct or struct pointer by using `db:"-"`

```go
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
```

Output:
```
UPDATE "user" SET "firstname"='Greg',"lastname"='Farley' []
```


<a name="set-map"></a>
**[Set with Map](#UpdateDataset.Set)**

```go
sql, args, _ := pp.Update("items").Set(
	map[string]interface{}{"name": "Test", "address": "111 Test Addr"},
).Build()
fmt.Println(sql, args)
```

Output:
```
UPDATE "items" SET "address"='111 Test Addr',"name"='Test' []
```

<a name="from"></a>
**[From / Multi Table](#UpdateDataset.From)**

`pp` allows joining multiple tables in a update clause through `From`.

**NOTE** The `sqlite3` adapter does not support a multi table syntax.

`Postgres` Example

```go
dialect := pp.Dialect("postgres")

ds := dialect.Update("table_one").
    Set(pp.Record{"foo": pp.I("table_two.bar")}).
    From("table_two").
    Where(pp.Ex{"table_one.id": pp.I("table_two.id")})

sql, _, _ := ds.Build()
fmt.Println(sql)
```

Output:
```sql
UPDATE "table_one" SET "foo"="table_two"."bar" FROM "table_two" WHERE ("table_one"."id" = "table_two"."id")
```

`MySQL` Example

```go
dialect := pp.Dialect("mysql")

ds := dialect.Update("table_one").
    Set(pp.Record{"foo": pp.I("table_two.bar")}).
    From("table_two").
    Where(pp.Ex{"table_one.id": pp.I("table_two.id")})

sql, _, _ := ds.Build()
fmt.Println(sql)
```
Output:
```sql
UPDATE `table_one`,`table_two` SET `foo`=`table_two`.`bar` WHERE (`table_one`.`id` = `table_two`.`id`)
```

<a name="where"></a>
**[Where](#UpdateDataset.Where)**

```go
sql, _, _ := pp.Update("test").
	Set(pp.Record{"foo": "bar"}).
	Where(pp.Ex{
		"a": pp.Op{"gt": 10},
		"b": pp.Op{"lt": 10},
		"c": nil,
		"d": []string{"a", "b", "c"},
	}).Build()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
```

<a name="order"></a>
**[Order](#UpdateDataset.Order)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "manlu.org/pp/dialect/mysql"

ds := pp.Dialect("mysql").
	Update("test").
	Set(pp.Record{"foo": "bar"}).
	Order(pp.C("a").Asc())
sql, _, _ := ds.Build()
fmt.Println(sql)
```

Output:
```
UPDATE `test` SET `foo`='bar' ORDER BY `a` ASC
```

<a name="limit"></a>
**[Order](#UpdateDataset.Limit)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "manlu.org/pp/dialect/mysql"

ds := pp.Dialect("mysql").
	Update("test").
	Set(pp.Record{"foo": "bar"}).
	Limit(10)
sql, _, _ := ds.Build()
fmt.Println(sql)
```

Output:
```
UPDATE `test` SET `foo`='bar' LIMIT 10
```

<a name="returning"></a>
**[Returning](#UpdateDataset.Returning)**

Returning a single column example.

```go
sql, _, _ := pp.Update("test").
	Set(pp.Record{"foo": "bar"}).
	Returning("id").
	Build()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "id"
```

Returning multiple columns

```go
sql, _, _ := pp.Update("test").
	Set(pp.Record{"foo": "bar"}).
	Returning("a", "b").
	Build()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "a", "b"
```

Returning all columns

```go
sql, _, _ := pp.Update("test").
	Set(pp.Record{"foo": "bar"}).
	Returning(pp.T("test").All()).
	Build()
fmt.Println(sql)
```

Output:
```
UPDATE "test" SET "foo"='bar' RETURNING "test".*
```

<a name="seterror"></a>
**[`SetError`](#UpdateDataset.SetError)**

Sometimes while building up a query with pp you will encounter situations where certain
preconditions are not met or some end-user contraint has been violated. While you could
track this error case separately, pp provides a convenient built-in mechanism to set an
error on a dataset if one has not already been set to simplify query building.

Set an Error on a dataset:

```go
func GetUpdate(name string, value string) *pp.UpdateDataset {

    var ds = pp.Update("test")

    if len(name) == 0 {
        return ds.SetError(fmt.Errorf("name is empty"))
    }

    if len(value) == 0 {
        return ds.SetError(fmt.Errorf("value is empty"))
    }

    return ds.Set(pp.Record{name: value})
}

```

This error is returned on any subsequent call to `Error` or `Build`:

```go
var field, value string
ds = GetUpdate(field, value)
fmt.Println(ds.Error())

sql, args, err = ds.Build()
fmt.Println(err)
```

Output:
```
name is empty
name is empty
```

<a name="executing"></a>
## Executing Updates

To execute Updates use [`pp.Database#Update`](#Database.Update) to create your dataset

### Examples

**Executing an update**
```go
db := getDb()

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
```

Output:

```
Updated 1 users
```

**Executing with Returning**

```go
db := getDb()

update := db.Update("pp_user").
	Set(pp.Record{"last_name": "ucon"}).
	Where(pp.Ex{"last_name": "Yukon"}).
	Returning("id").
	Executor()

var ids []int64
if err := update.ScanVals(&ids); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Printf("Updated users with ids %+v", ids)
}

```

Output:
```
Updated users with ids [1 2 3]
```
