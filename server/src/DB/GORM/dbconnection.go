package DB

import (
	"encoding/json"
	"fmt"
	//"github.com/go-sql-driver/mysql"

	//_ "github.com/jinzhu/gorm/dialects/mysql"
	//"github.com/jinzhu/gorm"
	//"github.com/jinzhu/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", configure.UserName, configure.PWD, configure.DBServer, configure.DBName)), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Access to the database %s success\n", configure.DBName)
	}
	//db.SingularTable(true)
	return db
}

var GormDB *gorm.DB

func init() {
	GormDB = NewDBConnection()

}
