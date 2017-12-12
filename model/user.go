package model

type User struct {
	ID         int64  `gorm:primary key;not_nil`
	Login      string `gorm:not_nil`
	Pass      string `gorm:not_nil`
}

//Проверка логина в бд
func (u User) Get(login string, password string) error {
	return DBConn.Where("login = ? and pass = ?", login, password).First(&u).Error
}

//Смена пароля
func (u User) ChangePassword(login string, newpassword string) error {
	return DBConn.Model(&u).Where("login = ?", login).Update("pass", newpassword).Error
}
