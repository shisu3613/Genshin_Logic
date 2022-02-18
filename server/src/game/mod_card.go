package game

import (
	"fmt"
	"server/csvs"
)

type Card struct {
	CardId int
}

type ModCard struct {
	CardInfo map[int]*Card
}

func (mc *ModCard) IsHasCard(cardId int) bool {
	_, ok := mc.CardInfo[cardId]
	return ok
}

func (mc *ModCard) AddItem(itemId int, friendliness int) {
	_, ok := mc.CardInfo[itemId]
	if ok {
		fmt.Println("已存在名片：", itemId)
		return
	}
	config := csvs.GetCardConfig(itemId)
	if config == nil {
		fmt.Println("非法名片：", itemId)
		return
	}
	if friendliness < config.Friendliness {
		fmt.Println("好感度不足：", itemId)
		return
	}

	mc.CardInfo[itemId] = &Card{CardId: itemId}
	fmt.Println("获得名片：", itemId)
}

func (mc *ModCard) CheckGetCard(roleId int, friendliness int) {
	config := csvs.GetCardConfigByRoleId(roleId)
	if config == nil {
		return
	}
	mc.AddItem(config.CardId, friendliness)
}
