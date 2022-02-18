package csvs

import (
	"server/utils"
)

const (
	ItemTypeNormal   = 1
	ItemTypeRole     = 2
	ItemTypeIcon     = 3
	ItemTypeCard     = 4
	ItemTypeWeapon   = 6
	ItemTypeRelic    = 7
	ItemTypeCookBook = 8
	ItemTypeCook     = 9
	ItemTypeFood     = 10
	ItemTypeFurn     = 11
)

type ConfigItem struct {
	ItemId   int    `json:"ItemId"`
	SortType int    `json:"SortType"`
	ItemName string `json:"ItemName"`
}

var (
	ConfigItemMap map[int]*ConfigItem
)

func init() {
	ConfigItemMap = make(map[int]*ConfigItem)
	utils.GetCsvUtilMgr().LoadCsv("Item", &ConfigItemMap)
	return
}

func GetItemConfig(itemId int) *ConfigItem {
	return ConfigItemMap[itemId]
}

func GetItemName(itemId int) string {
	config := GetItemConfig(itemId)
	if config == nil {
		return ""
	}
	return config.ItemName
}
