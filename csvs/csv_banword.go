package csvs

import "fmt"

type ConfigBanWord struct {
	Id  int
	Txt string
}

var ConfigBanWordSlice []*ConfigBanWord

func init() {
	ConfigBanWordSlice = append(ConfigBanWordSlice,
		&ConfigBanWord{
			Id:  1,
			Txt: "外挂",
		},
		&ConfigBanWord{
			Id:  2,
			Txt: "微信",
		}, &ConfigBanWord{
			Id:  3,
			Txt: "原神",
		}, &ConfigBanWord{
			Id:  4,
			Txt: "代练",
		}, &ConfigBanWord{
			Id:  5,
			Txt: "辅助",
		}, &ConfigBanWord{
			Id:  6,
			Txt: "赚钱",
		})
	fmt.Println("csv_banword初始化")
}

func GetBanWordBase() (relString []string) {
	for _, v := range ConfigBanWordSlice {
		relString = append(relString, v.Txt)
	}
	return
}
