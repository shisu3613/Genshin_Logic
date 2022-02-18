package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"server/DB"
	"server/csvs"
)

func main() {
	db := DB.NewDBConnection()
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("Close Datebase failure:", err)
		}
	}(db)
	//creatTable(&csvs.ConfigBanWord{}, db)
	var tests []*csvs.ConfigBanWord
	db.Find(&tests)
	for _, x := range tests {
		fmt.Println(*x)
	}
}

func creatTable(val interface{}, db *gorm.DB) {
	if !db.HasTable(val) {
		db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(val)
	}
}
