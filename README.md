# SI
SI (Secret Ingredient), is a database tool to handle query generation and relation handling. It can also be used to save models in the database.
The idea is to simplify the relation handling and to work only with model structs, instead of queries.

Its syntax is heavily based on the laravel eloquent, but very simplified and typed.


## Configuration

### Defining a simple model
```go
type artist struct {
    // Note: `si.Model` must be the first declared thing on the struct.
    si.Model
    
    Name string `si:"DB_COLUMN_NAME"`
}

func (a artist) GetModel() si.Model {
    return a.Model
}

func (a artist) GetTable() string {
    return "artists"
}
```
A model must implement the [`Modeler`](https://github.com/derivatan/si/blob/095f3ca8e974635a8ac20e8b2e327af27556c781/common.go#L20C2-L20C2) interface.

It must also embed the `si.Model` as the first declared field on the model.

All exported fields on the model should have a matching column in the database with the naming as `snake_case(FieldName)`.
This can be overwritten with the si-tag (`DB_COLUMN_NAME` in the example above). The tag can be excluded.


### Define Relationships

```go
type Album struct {
    si.Model

    Name string
    Year int
    ArtistID uuid.UUID

    // Note: must be unexported.
    artist si.RelationData[Artist] `si:FIELD_NAME`
}

func (a Album) GetModel() si.Model {
    return a.Model
}

func (a Album) GetTable() string {
    return "albums"
}

func (a Album) Artist() si.Relation[album, artist] {
    return si.BelongsTo[album, artist](a, "artist", func(a *album) *si.RelationData[artist] {
        return &a.artist
    })
}
```

A relationship is defined by two things.
* An **unexported** field with the type `si.RelationData[T]`
* An exported function that returns a `si.Relation[From, To]`

The field is only for si:s internal use and should not be used or modified in any way. To get a relation you must use the function as a query builder.


### Database setup
In order to be completely agnostic about the database, si uses these [interfaces](https://github.com/derivatan/si/blob/main/db.go) for database communication.
This is based on the `sql.DB`, but can easily be implemented with whatever you want to use. A simple example of such an implementaiton can be found in [`sql_wrap.go`](https://github.com/derivatan/si/blob/main/sql_wrap.go)


## Usage

* `si.Query[T]()` is the main entry point for retrieving data from the database.

  Example, Get all albums with that start with the leter 'a'.
  ```go
  albums, err := si.Query[album]().Where("name", "ILIKE", "a%").OrderBy("name", true).Get(db)
  ```

* `si.Save(model)` is used to create or update a model, with the values upon the model.
  To save relations, you must update the ID column, just as a normal column. This will **not** change what's stored in relation field if it is already loaded. 

* If you need to debug the generated queries, or get some silent errors, you can use `si.SetLogger(...)`.
  This logger will be called with all the queries and their arguments that _si_ generates, and might in some cases give some debugging messages. 

  This example will print all queries.
  ```go
  si.SetLogger(func(a ...any) {
      fmt.Println(a...)
  })
  ```


## Example and tests

There are integration tests for all major functionalities in a [separate repo](http://github.com/derivatan/si_test)
These are a also a good example for how to use the library.


### Comments

 * There is no mapping for a many-to-many relation yet. In order to achieve this anyway, you can make a model for the pivot-table and use the one-to-many relation in both directions to the other models.
