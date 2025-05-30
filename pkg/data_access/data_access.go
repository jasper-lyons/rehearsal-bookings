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
		err :=	rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}

		// create a new object of type
		instance := reflect.New(reflect.TypeFor[T]()).Elem()
		for _, field := range fields {
			sqlTag := field.Tag.Get("sql")
			if sqlTag == "" {
				continue
			}

			// identify the position of the column with a matching sql name
			for j, column := range columns {
				if column == sqlTag {
					fieldValue := instance.FieldByName(field.Name)
					if !fieldValue.IsValid() || !fieldValue.CanSet() {
						// we can't set this field
						continue
					}

					if values[j] == nil {
						// There is no value to set
						continue
					}

					dbValue := reflect.ValueOf(values[j])

					if dbValue.Type().AssignableTo(fieldValue.Type()) {
						fieldValue.Set(dbValue)
					} else {
						if err:= convertAndSet(fieldValue, values[j]); err != nil {
							return nil, fmt.Errorf("failed to set field %s: %w", field.Name, err)
						}
					}
				}
			}
		}
		// add the object to a list of object
		results = append(results, instance.Interface().(T))
	}

	return results, nil
}

func convertAndSet(fieldValue reflect.Value, value interface{}) error {
    switch fieldValue.Kind() {
    case reflect.Bool:
        // Convert SQLite integer to boolean (0 = false, non-zero = true)
        switch v := value.(type) {
        case int64:
            fieldValue.SetBool(v != 0)
        default:
            return fmt.Errorf("cannot convert %T to bool", value)
        }
    default:
        return fmt.Errorf("unsupported field type: %s", fieldValue.Type())
    }
    
    return nil
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

func generateDeleteStatement(tableName string, count int) string {
	placeholders := strings.Repeat("?,", count-1) + "?"
	return fmt.Sprintf(
		`delete from %s
		where id in (%s)`,
		tableName,
		placeholders,
	)
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

const TimeFormat = "2006-01-02 15:04:00"

func Time(t string) (time.Time, error) {
	return time.Parse(TimeFormat, t)
}
