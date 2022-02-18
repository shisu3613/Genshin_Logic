package csvs

import (
	"fmt"
	"os"
	"server/utils"
)

type ConfigCard struct {
	CardId       int `json:"CardId" gorm:"primary_key"`
	Friendliness int `json:"Friendliness"`
	Check        int `json:"Check"`
}

var (
	ConfigCardMap         map[int]*ConfigCard
	ConfigCardMapByRoleId map[int]*ConfigCard
)

func init() {

	ConfigCardMap = make(map[int]*ConfigCard)
	fmt.Println(os.Getwd())
	utils.GetCsvUtilMgr().LoadCsv("Card", &ConfigCardMap)
	ConfigCardMapByRoleId = make(map[int]*ConfigCard)
	for _, v := range ConfigCardMap {
		//err := db.Debug().Create(v).Error
		//if err != nil {
		//	fmt.Println(err)
		//}
		ConfigCardMapByRoleId[v.Check] = v
	}
	return
}

func GetCardConfig(cardId int) *ConfigCard {
	return ConfigCardMap[cardId]
}

func GetCardConfigByRoleId(roleId int) *ConfigCard {
	return ConfigCardMapByRoleId[roleId]
}
