package data_access

type User struct {
	UserName  string `sql:"user_name" json:"user_name"`
	UserEmail string `sql:"user_email" json:"user_email"`
	UserPhone string `sql:"user_phone" json:"user_phone"`
}

type UsersRepository[D StorageDriver] struct {
	driver D
}

func NewUsersRepository(sd StorageDriver) *UsersRepository[StorageDriver] {
	return &UsersRepository[StorageDriver]{driver: sd}
}

func (ur *UsersRepository[StorageDriver]) All() ([]User, error) {
	rows, err := ur.driver.Query("select user_name, user_email, user_phone from user_view")
	if err != nil {
		return nil, err
	}
	return RowsToType[User](rows)
}
