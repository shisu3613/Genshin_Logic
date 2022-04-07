package DB

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func DBtest(db *gorm.DB) {
	//db := NewDBConnection()
	//defer func(db *gorm.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		fmt.Println("Close Datebase failure:", err)
	//	}
	//}(db)

	type Like struct {
		ID        int    `gorm:"primary_key"`
		Ip        string `gorm:"type:varchar(20);not null;index:ip_idx"`
		Ua        string `gorm:"type:varchar(256);not null;"`
		Title     string `gorm:"type:varchar(128);not null;index:title_idx"`
		Hash      uint64 `gorm:"unique_index:hash_idx;"`
		CreatedAt time.Time
	}

	//if !db.Table(&Like{}) {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Create(&Like{}).Error; err != nil {
		panic(err)
	}
	//}

	like := &Like{
		Ip:        "1.1.1.1",
		Ua:        "List",
		Title:     "test",
		Hash:      12345678910,
		CreatedAt: time.Now(),
	}

	if err := db.Create(like).Error; err != nil {
		fmt.Println(err)
	}
}
