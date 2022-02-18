package csvs

import "server/utils"

type ConfigRelics struct {
	RelicsId int `json:"RelicsId"`
	Type     int `json:"Type"`
	Pos      int `json:"Pos"`
	Star     int `json:"Star"`
}

var (
	ConfigRelicsMap map[int]*ConfigRelics
)

func init() {
	ConfigRelicsMap = make(map[int]*ConfigRelics)
	utils.GetCsvUtilMgr().LoadCsv("Relics", &ConfigRelicsMap)
	return
}

func GetRelicsConfig(relicsId int) *ConfigRelics {
	return ConfigRelicsMap[relicsId]
}
