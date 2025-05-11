package data_access

import (
	"time"
	"fmt"
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
	if err != nil {
		return nil, err
	}

	if len(codes) == 0 || codes == nil {
		return nil, fmt.Errorf("No code found with name %s", code_name)
	}

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

type Codes struct {
	cr *CodesRepository[StorageDriver]
}

func NewCodes(cr *CodesRepository[StorageDriver]) Codes {
	return Codes {
		cr: cr,
	}
}

func (codes *Codes) GetCode(key string) (string, error) {
	code, err := codes.cr.Find(key)
	if err != nil {
		return "", err
	}

	return code.CodeValue, nil
}

func (codes *Codes) FrontDoorCodeFor(weekday time.Weekday) (string, error) {
	switch weekday {
	case time.Monday:
		return codes.GetCode("Monday Front Door")
	case time.Tuesday:
		return codes.GetCode("Tuesday Front Door")
	case time.Wednesday:
		return codes.GetCode("Wednesday Front Door")
	case time.Thursday:
		return codes.GetCode("Thursday Front Door")
	case time.Friday:
		return codes.GetCode("Friday Front Door")
	case time.Saturday:
		return codes.GetCode("Saturday Front Door")
	case time.Sunday:
		return codes.GetCode("Sunday Front Door")
	default:
		return "", fmt.Errorf("invalid weekday: %v", weekday)
	}
}

func (codes *Codes) RoomCodeFor(room string) (string, error) {
	switch room {
	case "Room 1":
		return codes.GetCode("Room 1")
	case "Room 2":
		return codes.GetCode("Room 1")
	case "Rec Room":
		return codes.GetCode("Rec Room Store")
	}
}
