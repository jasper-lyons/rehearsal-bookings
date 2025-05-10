package data_access

import (
	"time"
)

type Code struct {
	Id        int64     `sql:"id" generated:"true" json:"id"`
	CodeName  string    `sql:"code_name" json:"code_name"`
	CodeValue string    `sql:"code_value" json:"code_value"`
	UpdatedAt time.Time `sql:"updated_at" generated:"true" json:"updated_at"`
	CreatedAt time.Time `sql:"created_at" generated:"true" json:"created_at"`
}

type CodesRepository[D StorageDriver] struct {
	driver D
}

func NewCodesRepository(sd StorageDriver) *CodesRepository[StorageDriver] {
	return &CodesRepository[StorageDriver]{driver: sd}
}

func (cr *CodesRepository[StorageDriver]) Find(code_name string) (*Code, error) {
	rows, err := cr.driver.Query("select * from codes where code_name = ? limit 1", code_name)
	if err != nil {
		return nil, err
	}
	codes, err := RowsToType[Code](rows)
	return &codes[0], err
}

func (cr *CodesRepository[StorageDriver]) All() ([]Code, error) {
	rows, err := cr.driver.Query("select * from codes")
	if err != nil {
		return nil, err
	}
	return RowsToType[Code](rows)
}

func (cr *CodesRepository[StorageDriver]) Update(codes []Code) ([]Code, error) {
	statement := generateUpdateStatement[Code]("codes")
	err := cr.driver.Update(
		statement,
		ObjectsToUpdateParams(ToInterfaceSlice(codes)),
	)
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (cr *CodesRepository[StorageDriver]) Where(query string, params ...any) ([]Code, error) {
	rows, err := cr.driver.Query("select * from codes where "+query, params...)
	if err != nil {
		return nil, err
	}
	return RowsToType[Code](rows)
}
