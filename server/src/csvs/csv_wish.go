package csvs

import "server/utils"

type ConfigWishes struct {
	DropId int `json:"DropId"` //用于判断当前在第几次色子
	Weight int `json:"Weight"`
	Result int `json:"Result"`
	IsEnd  int `json:"IsEnd"`
}

var (
	ConfigWishesSlice []*ConfigWishes
)

func init() {
	ConfigWishesSlice = make([]*ConfigWishes, 0)
	utils.GetCsvUtilMgr().LoadCsv("WishDrop", &ConfigWishesSlice)

	return
}
