package csvs

import "server/utils"

type ConfigDropItem struct {
	DropId     int `json:"DropId"`
	DropType   int `json:"DropType"`
	Weight     int `json:"Weight"`
	ItemId     int `json:"ItemId"`
	ItemNumMin int `json:"ItemNumMin"`
	ItemNumMax int `json:"ItemNumMax"`
	WorldAdd   int `json:"WorldAdd"`
}

var (
	ConfigDropItemSlice []*ConfigDropItem
)

func init() {

	ConfigDropItemSlice = make([]*ConfigDropItem, 0)
	utils.GetCsvUtilMgr().LoadCsv("DropItem", &ConfigDropItemSlice)
	return
}
