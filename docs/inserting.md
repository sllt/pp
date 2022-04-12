# Inserting

* [Creating An InsertDataset](#create)
* Examples
  * [Insert Cols and Vals](#insert-cols-vals)
  * [Insert `pp.Record`](#insert-record)
  * [Insert Structs](#insert-structs)
  * [Insert Map](#insert-map)
  * [Insert From Query](#insert-from-query)
  * [Returning](#returning)
  * [SetError](#seterror)
  * [Executing](#executing)

<a name="create"></a>
To create a [`InsertDataset`](#InsertDataset)  you can use

**[`pp.Insert`](#Insert)**

When you just want to create some quick SQL, this mostly follows the `Postgres` with the exception of placeholders for prepared statements.

```go
ds := pp.Insert("user").Rows(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
insertSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley')
```

**[`SelectDataset.Insert`](#SelectDataset.Insert)**

If you already have a `SelectDataset` you can invoke `Insert()` to get a `InsertDataset`

**NOTE** This method will also copy over the `WITH` clause as well as the `FROM`

```go
ds := pp.From("user")

ds := ds.Insert().Rows(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
insertSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley')
```

**[`DialectWrapper.Insert`](#DialectWrapper.Insert)**

Use this when you want to create SQL for a specific `dialect`

```go
// import _ "github.com/doug-martin/pp/v9/dialect/mysql"

dialect := pp.Dialect("mysql")

ds := dialect.Insert().Rows(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
insertSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
INSERT INTO `user` (`first_name`, `last_name`) VALUES ('Greg', 'Farley')
```

**[`Database.Insert`](#DialectWrapper.Insert)**

Use this when you want to execute the SQL or create SQL for the drivers dialect.

```go
// import _ "github.com/doug-martin/pp/v9/dialect/mysql"

mysqlDB := //initialize your db
db := pp.New("mysql", mysqlDB)

ds := db.Insert().Rows(
    pp.Record{"first_name": "Greg", "last_name": "Farley"},
)
insertSQL, _, _ := ds.Build()
fmt.Println(insertSQL, args)
```
Output:
```
INSERT INTO `user` (`first_name`, `last_name`) VALUES ('Greg', 'Farley')
```

### Examples

For more examples visit the **[Docs](#InsertDataset)**

<a name="insert-cols-vals"></a>
**Insert with Cols and Vals**

```go
ds := pp.Insert("user").
	Cols("first_name", "last_name").
	Vals(
		pp.Vals{"Greg", "Farley"},
		pp.Vals{"Jimmy", "Stewart"},
		pp.Vals{"Jeff", "Jeffers"},
	)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```sql
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

<a name="insert-record"></a>
**Insert `pp.Record`**

```go
ds := pp.Insert("user").Rows(
	pp.Record{"first_name": "Greg", "last_name": "Farley"},
	pp.Record{"first_name": "Jimmy", "last_name": "Stewart"},
	pp.Record{"first_name": "Jeff", "last_name": "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

<a name="insert-structs"></a>
**Insert Structs**

```go
type User struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}
ds := pp.Insert("user").Rows(
	User{FirstName: "Greg", LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

You can skip fields in a struct by using the `skipinsert` tag

```go
type User struct {
	FirstName string `db:"first_name" pp:"skipinsert"`
	LastName  string `db:"last_name"`
}
ds := pp.Insert("user").Rows(
	User{FirstName: "Greg", LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("last_name") VALUES ('Farley'), ('Stewart'), ('Jeffers') []
```

If you want to use the database `DEFAULT` when the struct field is a zero value you can use the `defaultifempty` tag.

```go
type User struct {
	FirstName string `db:"first_name" pp:"defaultifempty"`
	LastName  string `db:"last_name"`
}
ds := pp.Insert("user").Rows(
	User{LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES (DEFAULT, 'Farley'), ('Jimmy', 'Stewart'), (DEFAULT, 'Jeffers') []
```

`pp` will also use fields in embedded structs when creating an insert.

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
ds := pp.Insert("user").Rows(
	User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
	User{Address: Address{Street: "211 Street", State: "NY"}, FirstName: "Jimmy", LastName: "Stewart"},
	User{Address: Address{Street: "311 Street", State: "NY"}, FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("address_state", "address_street", "firstname", "lastname") VALUES ('NY', '111 Street', 'Greg', 'Farley'), ('NY', '211 Street', 'Jimmy', 'Stewart'), ('NY', '311 Street', 'Jeff', 'Jeffers') []
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
ds := pp.Insert("user").Rows(
	User{FirstName: "Greg", LastName: "Farley"},
	User{FirstName: "Jimmy", LastName: "Stewart"},
	User{FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("firstname", "lastname") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
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

ds := pp.Insert("user").Rows(
	User{Address: Address{Street: "111 Street", State: "NY"}, FirstName: "Greg", LastName: "Farley"},
	User{Address: Address{Street: "211 Street", State: "NY"}, FirstName: "Jimmy", LastName: "Stewart"},
	User{Address: Address{Street: "311 Street", State: "NY"}, FirstName: "Jeff", LastName: "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("firstname", "lastname") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

<a name="insert-map"></a>
**Insert `map[string]interface{}`**

```go
ds := pp.Insert("user").Rows(
	map[string]interface{}{"first_name": "Greg", "last_name": "Farley"},
	map[string]interface{}{"first_name": "Jimmy", "last_name": "Stewart"},
	map[string]interface{}{"first_name": "Jeff", "last_name": "Jeffers"},
)
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") VALUES ('Greg', 'Farley'), ('Jimmy', 'Stewart'), ('Jeff', 'Jeffers') []
```

<a name="insert-from-query"></a>
**Insert from query**

```go
ds := pp.Insert("user").Prepared(true).
	FromQuery(pp.From("other_table"))
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" SELECT * FROM "other_table" []
```

You can also specify the columns

```go
ds := pp.Insert("user").Prepared(true).
	Cols("first_name", "last_name").
	FromQuery(pp.From("other_table").Select("fn", "ln"))
insertSQL, args, _ := ds.Build()
fmt.Println(insertSQL, args)
```

Output:
```
INSERT INTO "user" ("first_name", "last_name") SELECT "fn", "ln" FROM "other_table" []
```

<a name="returning"></a>
**Returning Clause**

Returning a single column example.

```go
sql, _, _ := pp.Insert("test").
	Rows(pp.Record{"a": "a", "b": "b"}).
	Returning("id").
	Build()
fmt.Println(sql)
```

Output:
```
INSERT INTO "test" ("a", "b") VALUES ('a', 'b') RETURNING "id"
```

Returning multiple columns

```go
sql, _, _ = pp.Insert("test").
	Rows(pp.Record{"a": "a", "b": "b"}).
	Returning("a", "b").
	Build()
fmt.Println(sql)
```

Output:
```
INSERT INTO "test" ("a", "b") VALUES ('a', 'b') RETURNING "a", "b"
```

Returning all columns

```go
sql, _, _ = pp.Insert("test").
	Rows(pp.Record{"a": "a", "b": "b"}).
	Returning(pp.T("test").All()).
	Build()
fmt.Println(sql)
```

Output:
```
INSERT INTO "test" ("a", "b") VALUES ('a', 'b') RETURNING "test".*
```

<a name="seterror"></a>
**[`SetError`](#InsertDataset.SetError)**

Sometimes while building up a query with pp you will encounter situations where certain
preconditions are not met or some end-user contraint has been violated. While you could
track this error case separately, pp provides a convenient built-in mechanism to set an
error on a dataset if one has not already been set to simplify query building.

Set an Error on a dataset:

```go
func GetInsert(name string, value string) *pp.InsertDataset {

    var ds = pp.Insert("test")

    if len(field) == 0 {
        return ds.SetError(fmt.Errorf("name is empty"))
    }

    if len(value) == 0 {
        return ds.SetError(fmt.Errorf("value is empty"))
    }

    return ds.Rows(pp.Record{name: value})
}

```

This error is returned on any subsequent call to `Error` or `Build`:

```go
var field, value string
ds = GetInsert(field, value)
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
## Executing Inserts

To execute INSERTS use [`Database.Insert`](#Database.Insert) to create your dataset

### Examples

**Executing an single Insert**
```go
db := getDb()

insert := db.Insert("pp_user").Rows(
	pp.Record{"first_name": "Jed", "last_name": "Riley", "created": time.Now()},
).Executor()

if _, err := insert.Exec(); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Println("Inserted 1 user")
}
```

Output:

```
Inserted 1 user
```

**Executing multiple inserts**

```go
db := getDb()

users := []pp.Record{
	{"first_name": "Greg", "last_name": "Farley", "created": time.Now()},
	{"first_name": "Jimmy", "last_name": "Stewart", "created": time.Now()},
	{"first_name": "Jeff", "last_name": "Jeffers", "created": time.Now()},
}

insert := db.Insert("pp_user").Rows(users).Executor()
if _, err := insert.Exec(); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Printf("Inserted %d users", len(users))
}

```

Output:
```
Inserted 3 users
```

If you use the RETURNING clause you can scan into structs or values.

```go
db := getDb()

insert := db.Insert("pp_user").Returning(pp.C("id")).Rows(
		pp.Record{"first_name": "Jed", "last_name": "Riley", "created": time.Now()},
).Executor()

var id int64
if _, err := insert.ScanVal(&id); err != nil {
	fmt.Println(err.Error())
} else {
	fmt.Printf("Inserted 1 user id:=%d\n", id)
}
```

Output:

```
Inserted 1 user id:=5
```
