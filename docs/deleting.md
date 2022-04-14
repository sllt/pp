# Deleting

* [Creating A DeleteDataset](#create)
* Examples
  * [Delete All](#delete-all)
  * [Prepared](#prepared)
  * [Where](#where)
  * [Order](#order)
  * [Limit](#limit)
  * [Returning](#returning)
  * [SetError](#seterror)
  * [Executing](#exec)

<a name="create"></a>
To create a [`DeleteDataset`](#DeleteDataset)  you can use

**[`pp.Delete`](#Delete)**

When you just want to create some quick SQL, this mostly follows the `Postgres` with the exception of placeholders for prepared statements.

```go
sql, _, _ := pp.Delete("table").Build()
fmt.Println(sql)
```
Output:
```
DELETE FROM "table"
```

**[`SelectDataset.Delete`](#SelectDataset.Delete)**

If you already have a `SelectDataset` you can invoke `Delete()` to get a `DeleteDataset`

**NOTE** This method will also copy over the `WITH`, `WHERE`, `ORDER`, and `LIMIT` from the `SelectDataset`

```go

ds := pp.From("table")

sql, _, _ := ds.Delete().Build()
fmt.Println(sql)

sql, _, _ = ds.Where(pp.C("foo").Eq("bar")).Delete().Build()
fmt.Println(sql)
```
Output:
```
DELETE FROM "table"
DELETE FROM "table" WHERE "foo"='bar'
```

**[`DialectWrapper.Delete`](#DialectWrapper.Delete)**

Use this when you want to create SQL for a specific `dialect`

```go
// import _ "manlu.org/pp/dialect/mysql"

dialect := pp.Dialect("mysql")

sql, _, _ := dialect.Delete("table").Build()
fmt.Println(sql)
```
Output:
```
DELETE FROM `table`
```

**[`Database.Delete`](#DialectWrapper.Delete)**

Use this when you want to execute the SQL or create SQL for the drivers dialect.

```go
// import _ "manlu.org/pp/dialect/mysql"

mysqlDB := //initialize your db
db := pp.New("mysql", mysqlDB)

sql, _, _ := db.Delete("table").Build()
fmt.Println(sql)
```
Output:
```
DELETE FROM `table`
```

### Examples

For more examples visit the **[Docs](#DeleteDataset)**

<a name="delete-all"></a>
**Delete All Records**

```go
ds := pp.Delete("items")

sql, args, _ := ds.Build()
fmt.Println(sql, args)
```

Output:
```
DELETE FROM "items" []
```

<a name="prepared"></a>
**[`Prepared`](#DeleteDataset.Prepared)**

```go
sql, _, _ := pp.Delete("test").Where(pp.Ex{
	"a": pp.Op{"gt": 10},
	"b": pp.Op{"lt": 10},
	"c": nil,
	"d": []string{"a", "b", "c"},
}).Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM "test" WHERE (("a" > ?) AND ("b" < ?) AND ("c" IS NULL) AND ("d" IN (?, ?, ?))) [10 10 a b c]
```

<a name="where"></a>
**[`Where`](#DeleteDataset.Where)**

```go
sql, _, _ := pp.Delete("test").Where(pp.Ex{
	"a": pp.Op{"gt": 10},
	"b": pp.Op{"lt": 10},
	"c": nil,
	"d": []string{"a", "b", "c"},
}).Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM "test" WHERE (("a" > 10) AND ("b" < 10) AND ("c" IS NULL) AND ("d" IN ('a', 'b', 'c')))
```

<a name="order"></a>
**[`Order`](#DeleteDataset.Order)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "manlu.org/pp/dialect/mysql"

ds := pp.Dialect("mysql").Delete("test").Order(pp.C("a").Asc())
sql, _, _ := ds.Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM `test` ORDER BY `a` ASC
```

<a name="limit"></a>
**[`Limit`](#DeleteDataset.Limit)**

**NOTE** This will only work if your dialect supports it

```go
// import _ "manlu.org/pp/dialect/mysql"

ds := pp.Dialect("mysql").Delete("test").Limit(10)
sql, _, _ := ds.Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM `test` LIMIT 10
```

<a name="returning"></a>
**[`Returning`](#DeleteDataset.Returning)**

Returning a single column example.

```go
ds := pp.Delete("items")
sql, args, _ := ds.Returning("id").Build()
fmt.Println(sql, args)
```

Output:
```
DELETE FROM "items" RETURNING "id" []
```

Returning multiple columns

```go
sql, _, _ := pp.Delete("test").Returning("a", "b").Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM "items" RETURNING "a", "b"
```

Returning all columns

```go
sql, _, _ := pp.Delete("test").Returning(pp.T("test").All()).Build()
fmt.Println(sql)
```

Output:
```
DELETE FROM "test" RETURNING "test".*
```

<a name="seterror"></a>
**[`SetError`](#DeleteDataset.SetError)**

Sometimes while building up a query with pp you will encounter situations where certain
preconditions are not met or some end-user contraint has been violated. While you could
track this error case separately, pp provides a convenient built-in mechanism to set an
error on a dataset if one has not already been set to simplify query building.

Set an Error on a dataset:

```go
func GetDelete(name string, value string) *pp.DeleteDataset {

    var ds = pp.Delete("test")

    if len(name) == 0 {
        return ds.SetError(fmt.Errorf("name is empty"))
    }

    if len(value) == 0 {
        return ds.SetError(fmt.Errorf("value is empty"))
    }

    return ds.Where(pp.C(name).Eq(value))
}

```

This error is returned on any subsequent call to `Error` or `Build`:

```go
var field, value string
ds = GetDelete(field, value)
fmt.Println(ds.Error())

sql, args, err = ds.Build()
fmt.Println(err)
```

Output:
```
name is empty
name is empty
```

## Executing Deletes

To execute DELETES use [`Database.Delete`](#Database.Delete) to create your dataset

### Examples

<a name="exec"></a>
**Executing a Delete**
```go
db := getDb()

de := db.Delete("pp_user").
	Where(pp.Ex{"first_name": "Bob"}).
	Executor()

if r, err := de.Exec(); err != nil {
	fmt.Println(err.Error())
} else {
	c, _ := r.RowsAffected()
	fmt.Printf("Deleted %d users", c)
}
```

Output:

```
Deleted 1 users
```

If you use the RETURNING clause you can scan into structs or values.

```go
db := getDb()

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
```

Output:

```
Deleted users [ids:=[1 2 3]]
```
