package data_access

import (
	"time"
)

type Booking struct {
	Id             int64     `sql:"id" generated:"true" json:"id"`
	Type           string    `sql:"type" json:"type"`
	CustomerName   string    `sql:"customer_name" json:"customer_name"`
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
	BookingNotes   string    `sql:"booking_notes" json:"booking_notes"`
	TransactionId  string    `sql:"transaction_id" json:"transaction_id"`
	PaymentMethod  string    `sql:"payment_method" json:"payment_method"`
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

func (br *BookingsRepository[StorageDriver]) Where(query string, params ...any) ([]Booking, error) {
	rows, err := br.driver.Query("select * from bookings where "+query, params...)
	if err != nil {
		return nil, err
	}
	return RowsToType[Booking](rows)
}
