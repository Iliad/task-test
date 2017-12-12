package model

import (
	"fmt"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"../config"
)

var (
	DBConn *gorm.DB
)

//Подключение к бд
func GormInit() {

	var dc config.DBConfig;

	dc.Read();

	var err error;
	DBConn, err = gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dc.DBHost, dc.DBUser, dc.DBName, dc.DBPass))
	if err != nil {
		log.Fatal("Can't open db connection. Error: ", err)
	}
}

//Закрытие бд
func GormClose() error {
	if DBConn != nil {
		return DBConn.Close()
	}
	return nil
}
