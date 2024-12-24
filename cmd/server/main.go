package main

import (
	"fmt"
	"log"
	"database/sql"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"reflect"
	// GO Embed doesn't support embedding files from outside the module boundary
	// (in this case, anything outside of this directory) but we want to store
	// templates and static files at the route directory so we need to treat them
	// as their own go modules (you'll see the main.go files in the web/static and
	// web/templates directories) so that we can import them here and access the 
	// embedded static files!
	// "rehearsal-bookings/web/static"
	// "rehearsal-bookings/web/templates"
)

type SqliteDriver struct {
	db *sql.DB
}

func NewSqliteDriver(url string) *SqliteDriver {
	db, err := sql.Open("sqlite3", url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to open sqlite3 db. %s: %s", url, err))
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &SqliteDriver { db: db }
}

func (d* SqliteDriver) Query(query string, params ...any) (*sql.Rows, error) {
	return d.db.Query(query, params...)
}

func RowsToType[T any](rows* sql.Rows) ([]T, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i, _ := range columns {
		pointers[i] = &values[i]
	}

	var results []T

	for rows.Next() {
		rows.Scan(pointers...)
		// create a new object of type
		Type := reflect.TypeFor[T]()
		instance := reflect.New(Type).Elem()
		// go through each public struct
		for i := 0; i < instance.NumField(); i++ {
			// extract it's sql name
			fieldType := Type.Field(i)
			sqlName := fieldType.Tag.Get("sql")
			field := instance.Field(i)
			// identify the position of the column with a matching sql name
			for i, column := range columns {
				if column == sqlName {
					// set the value on the object to the column value from the row in values
					field.Set(reflect.ValueOf(values[i]))
				}
			}
		}
		// add the object to a list of object
		results = append(results, instance.Interface().(T))
	}

	return results, nil
}

type StorageDriver interface {
	Query(query string, params ...any) (*sql.Rows, error);
}

type BookingsRepository[D StorageDriver] struct {
	driver D
}

func NewBookingsRepository(sd StorageDriver) (*BookingsRepository[StorageDriver]) {
	return &BookingsRepository[StorageDriver] { driver: sd }
}

func (br* BookingsRepository[StorageDriver]) Find(id int) ([]Booking, error) {
	rows, err := br.driver.Query("select * from bookings where id = ?", id)
	if err != nil {
		return nil, err
	}
	return RowsToType[Booking](rows)
}

type Booking struct {
	Id int64 `sql:"id"`
	CustomerName string `sql:"customer_name"`
}

func main() {
	driver := NewSqliteDriver("db/development.db")
	br := NewBookingsRepository(driver)
	bookings, err := br.Find(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(bookings)
}
