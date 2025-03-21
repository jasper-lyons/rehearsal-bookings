package data_access

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func fieldsFor[T any]() []reflect.StructField {
	Type := reflect.TypeFor[T]()
	results := make([]reflect.StructField, Type.NumField())

	for i := 0; i < Type.NumField(); i++ {
		results[i] = Type.Field(i)
	}

	return results
}

func settableSqlColumnNames(fields []reflect.StructField) []string {
	var results []string

	for _, field := range fields {
		_, generated := field.Tag.Lookup("generated")
		columnName := field.Tag.Get("sql")
		if columnName != "" && !generated {
			results = append(results, columnName)
		}
	}

	return results
}

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

	return &SqliteDriver{db: db}
}

func (d *SqliteDriver) Query(query string, params ...any) (*sql.Rows, error) {
	return d.db.Query(query, params...)
}

func (d *SqliteDriver) Insert(statement string, records [][]interface{}) ([]int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	compiled, err := tx.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer compiled.Close()

	ids := make([]int64, len(records))
	for i, record := range records {
		result, err := compiled.Exec(record...)
		if err != nil {
			return nil, err
		}
		id, err := result.LastInsertId()
		if err == nil {
			ids[i] = id
		} else {
			ids[i] = -1
		}
	}

	return ids, tx.Commit()
}

func (d *SqliteDriver) Delete(statement string, ids []int64) (int64, error) {
	result, err := d.db.Exec(statement, ToInterfaceSlice(ids)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (d *SqliteDriver) Update(statement string, records [][]interface{}) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	compiled, err := tx.Prepare(statement)
	if err != nil {
		return err
	}
	defer compiled.Close()

	for _, record := range records {
		_, err := compiled.Exec(record...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func RowsToType[T any](rows *sql.Rows) ([]T, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i, _ := range columns {
		pointers[i] = &values[i]
	}
	fields := fieldsFor[T]()

	var results []T

	for rows.Next() {
		rows.Scan(pointers...)
		// create a new object of type
		instance := reflect.New(reflect.TypeFor[T]()).Elem()
		for _, field := range fields {
			// identify the position of the column with a matching sql name
			for j, column := range columns {
				if column == field.Tag.Get("sql") {
					// set the value on the object to the column value from the row in values
					instance.FieldByName(field.Name).Set(reflect.ValueOf(values[j]))
				}
			}
		}
		// add the object to a list of object
		results = append(results, instance.Interface().(T))
	}

	return results, nil
}

func ObjectsToInsertParams(objects []interface{}) [][]interface{} {
	var results [][]interface{}

	for _, object := range objects {
		value := reflect.ValueOf(object)
		Type := value.Type()

		var params []interface{}
		for i := 0; i < value.NumField(); i++ {
			field := Type.Field(i)
			if _, skip := field.Tag.Lookup("generated"); !skip {
				params = append(params, value.Field(i).Interface())
			}
		}

		results = append(results, params)
	}

	return results
}

func ObjectsToUpdateParams(objects []interface{}) [][]interface{} {
	var results [][]interface{}

	for _, object := range objects {
		value := reflect.ValueOf(object)
		Type := value.Type()

		var params []interface{}
		for i := 0; i < value.NumField(); i++ {
			field := Type.Field(i)
			if _, skip := field.Tag.Lookup("generated"); !skip {
				params = append(params, value.Field(i).Interface())
			}
		}

		// id will be the last param
		params = append(params, value.FieldByName("Id").Interface())

		results = append(results, params)
	}

	return results
}

func ToInterfaceSlice[T any](slice []T) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

type StorageDriver interface {
	Query(query string, params ...any) (*sql.Rows, error)
	Insert(statement string, records [][]interface{}) ([]int64, error)
	Delete(statement string, ids []int64) (int64, error)
	Update(statement string, records [][]interface{}) error
}

type BookingsRepository[D StorageDriver] struct {
	driver D
}

func NewBookingsRepository(sd StorageDriver) *BookingsRepository[StorageDriver] {
	return &BookingsRepository[StorageDriver]{driver: sd}
}

func (br *BookingsRepository[StorageDriver]) Find(id int) (*Booking, error) {
	rows, err := br.driver.Query("select * from bookings where id = ? limit 1", id)
	if err != nil {
		return nil, err
	}
	bookings, err := RowsToType[Booking](rows)
	return &bookings[0], err
}

func (br *BookingsRepository[StorageDriver]) All() ([]Booking, error) {
	rows, err := br.driver.Query("select * from bookings")
	if err != nil {
		return nil, err
	}
	return RowsToType[Booking](rows)
}

func (br *BookingsRepository[StorageDriver]) Where(query string, params ...any) ([]Booking, error) {
	rows, err := br.driver.Query("select * from bookings where "+query, params...)
	if err != nil {
		return nil, err
	}
	return RowsToType[Booking](rows)
}

func generateCreateStatement[T any](tableName string) string {
	columns := settableSqlColumnNames(fieldsFor[T]())
	placeholders := strings.Repeat("?,", len(columns)-1) + "?"

	return fmt.Sprintf(`
	insert into %s
		(%s)
	values
		(%s)`,
		tableName,
		strings.Join(columns, ", "),
		placeholders,
	)
}

func (br *BookingsRepository[StorageDriver]) Create(bookings []Booking) ([]Booking, error) {
	query := generateCreateStatement[Booking]("bookings")
	ids, err := br.driver.Insert(
		query,
		ObjectsToInsertParams(ToInterfaceSlice(bookings)),
	)
	if err != nil {
		return nil, err
	}
	for i := range bookings {
		if ids[i] > 0 {
			bookings[i].Id = ids[i]
		}
	}
	return bookings, nil
}

func generateDeleteStatement(tableName string, count int) string {
	placeholders := strings.Repeat("?,", count-1) + "?"
	return fmt.Sprintf(
		`delete from %s
		where id in (%s)`,
		tableName,
		placeholders,
	)
}

func (br *BookingsRepository[StorageDriver]) Delete(bookings []Booking) ([]Booking, error) {
	ids := make([]int64, len(bookings))
	for i, booking := range bookings {
		ids[i] = booking.Id
	}

	_, err := br.driver.Delete(
		generateDeleteStatement("bookings", len(ids)),
		ids,
	)

	return bookings, err
}

func generateUpdateStatement[T any](tableName string) string {
	columns := settableSqlColumnNames(fieldsFor[T]())

	sets := make([]string, len(columns))
	for i, column := range columns {
		sets[i] = fmt.Sprintf("%s = ?", column)
	}

	return fmt.Sprintf(
		"update %s set %s, updated_at = current_timestamp where id = ?",
		tableName,
		strings.Join(sets, ", "),
	)
}

func (br *BookingsRepository[StorageDriver]) Update(bookings []Booking) ([]Booking, error) {
	statement := generateUpdateStatement[Booking]("bookings")
	err := br.driver.Update(
		statement,
		ObjectsToUpdateParams(ToInterfaceSlice(bookings)),
	)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

type Booking struct {
	Id             int64     `sql:"id" generated:"true" json:"id"`
	Type           string    `sql:"type" json:"type"`
	CustomerName   string    `sql:"customer_name" json:"cutomer_name"`
	CustomerEmail  string    `sql:"customer_email" json:"customer_email"`
	CustomerPhone  string    `sql:"customer_phone" json:"customer_phone"`
	RoomName       string    `sql:"room_name" json:"room_name"`
	StartTime      time.Time `sql:"start_time" json:"start_time"`
	EndTime        time.Time `sql:"end_time" json:"end_time"`
	Status         string    `sql:"status" json:"status"`
	Expiration     time.Time `sql:"expiration" json:"expiration"`
	Price          float64   `sql:"price" json:"price"`
	DiscountAmount float64   `sql:"discount_amount" json:"discount_amount"`
	Cymbals        int64     `sql:"cymbals" json:"cymbals"`
	TransactionId  string    `sql:"transaction_id" json:"transaction_id"`
}

func FuckTheError[T any](result T, err error) T {
	if err != nil {
		fmt.Println(err)
	}
	return result
}

const TimeFormat = "2006-01-02 15:04:00"

func Time(t string) (time.Time, error) {
	return time.Parse(TimeFormat, t)
}
