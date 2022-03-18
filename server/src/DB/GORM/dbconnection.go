package GORM

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
)

type DBConfig struct {
	UserName string
	DBServer string
	PWD      string
	DBName   string
}

func NewDBConnection() *gorm.DB {
	var configure *DBConfig
	data, err := ioutil.ReadFile("json/DBconfig.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &configure)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configure.UserName, configure.PWD, configure.DBServer, configure.DBName)
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configure.UserName, configure.PWD, configure.DBServer, configure.DBName))
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Access to the database %s success\n", configure.DBName)
	}
	db.SingularTable(true)
	return db
}
