package config

import (
	"github.com/astaxie/beego/config"
	"log"
)

type DBConfig struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
}

func (dc *DBConfig) Read() error {
	fullConfigIni, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		log.Fatal("Can't open config file. Error: ", err)
	}

	configIni, err := fullConfigIni.GetSection("default")

	//Если невозможно открыть конфиг, то приложение завершается с ошибкой
	if err != nil {
		log.Fatal("Can't open config file. Error: ", err)
	}

	dc.DBUser = configIni["user"]
	dc.DBPass = configIni["pass"]
	dc.DBName = configIni["name"]
	dc.DBHost = configIni["host"]

	log.Println("Config loaded")
	return nil
}