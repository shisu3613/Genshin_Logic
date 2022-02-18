package csvs

type ConfigBanWord struct {
	Id  int `gorm:"primary_key"`
	Txt string
}

var (
	ConfigBanWordSlice []*ConfigBanWord
)

func init() {
	//db := DB.NewDBConnection()
	//defer func(db *gorm.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		fmt.Println("Close Datebase failure:", err)
	//	}
	//}(db)
	ConfigBanWordSlice = append(ConfigBanWordSlice,
		&ConfigBanWord{Id: 1, Txt: "外挂"},
		&ConfigBanWord{Id: 2, Txt: "辅助"},
		&ConfigBanWord{Id: 3, Txt: "微信"},
		&ConfigBanWord{Id: 4, Txt: "代练"},
		&ConfigBanWord{Id: 5, Txt: "赚钱"},
	)

	//for _, v := range ConfigBanWordSlice {
	//	err := db.Debug().Create(v).Error
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

}

func GetBanWordBase() []string {
	relString := make([]string, 0)
	for _, v := range ConfigBanWordSlice {
		relString = append(relString, v.Txt)
	}
	return relString
}
