package game

import (
	"fmt"
	"server/csvs"
)

type Relic struct {
	RelicId int
	KeyId   int
}

type ModRelic struct {
	RelicInfo map[int]*Relic
	MaxKey    int //永远递增的编号
}

func (mr *ModRelic) AddItem(itemId int, num int) {
	//武器表的验证环节
	config := csvs.GetRelicsConfig(itemId)
	if config == nil {
		fmt.Println("非法圣遗物")
		return
	}
	if len(mr.RelicInfo)+num > csvs.MaxRelicSize {
		fmt.Println("圣遗物背包已满！！！")
		return
	}
	for i := 0; i < num; i++ {
		relic := new(Relic)
		relic.RelicId = itemId
		mr.MaxKey++
		relic.KeyId = mr.MaxKey
		mr.RelicInfo[relic.KeyId] = relic
		fmt.Println("获得圣遗物：", csvs.GetItemName(itemId), "------圣遗物星级：", config.Star, "-----圣遗物编号：", relic.KeyId)
	}

}

func (mr *ModRelic) RemoveItem(keyId int) {
	if _, ok := mr.RelicInfo[keyId]; !ok {
		fmt.Println("当前编号圣遗物不存在")
		return
	}
	fmt.Println("移除圣遗物id为：", keyId, "移除圣遗物名称为", csvs.GetItemName(mr.RelicInfo[keyId].RelicId))
	delete(mr.RelicInfo, keyId)
}
