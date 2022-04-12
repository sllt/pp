### Database

The Database also allows you to execute queries but expects raw SQL to execute. The supported methods are

* [`Exec`](#Database.Exec)
* [`Prepare`](#Database.Prepare)
* [`Query`](#Database.Query)
* [`QueryRow`](#Database.QueryRow)
* [`ScanStructs`](#Database.ScanStructs)
* [`ScanStruct`](#Database.ScanStruct)
* [`ScanVals`](#Database.ScanVals)
* [`ScanVal`](#Database.ScanVal)
* [`Begin`](#Database.Begin)

### Transactions

`pp` has builtin support for transactions to make the use of the Datasets and querying seamless

```go
tx, err := db.Begin()
if err != nil{
   return err
}
//use tx.From to get a dataset that will execute within this transaction
update := tx.From("user").
    Where(pp.Ex{"password": nil}).
    Update(pp.Record{"status": "inactive"})
if _, err = update.Exec(); err != nil{
    if rErr := tx.Rollback(); rErr != nil{
        return rErr
    }
    return err
}
if err = tx.Commit(); err != nil{
    return err
}
return
```

The [`TxDatabase`](#TxDatabase)  also has all methods that the [`Database`](#Database) has along with

* [`Commit`](#TxDatabase.Commit)
* [`Rollback`](#TxDatabase.Rollback)
* [`Wrap`](#TxDatabase.Wrap)

#### Wrap

The [`TxDatabase.Wrap`](#TxDatabase.Wrap) is a convience method for automatically handling `COMMIT` and `ROLLBACK`

```go
tx, err := db.Begin()
if err != nil{
   return err
}
err = tx.Wrap(func() error{
  update := tx.From("user").
      Where(pp.Ex{"password": nil}).
      Update(pp.Record{"status": "inactive"})
  return update.Exec()
})
//err will be the original error from the update statement, unless there was an error executing ROLLBACK
if err != nil{
    return err
}
```

## Logging

To enable trace logging of SQL statements use the [`Database.Logger`](#Database.Logger) method to set your logger.

**NOTE** The logger must implement the [`Logger`](#Logger) interface

**NOTE** If you start a transaction using a database your set a logger on the transaction will inherit that logger automatically

