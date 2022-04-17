package game

import (
	"fmt"
	"server/csvs"
)

type HomeItem struct {
	ItemId  int
	ItemNum int64
	KeyId   int
}

type ModHome struct {
	HomeItemInfo     map[int]*HomeItem
	UsedHomeItemInfo map[int]*HomeItem
	player           *Player
}

func (mh *ModHome) AddItem(itemId int, num int64) {
	_, ok := mh.HomeItemInfo[itemId]
	if ok {
		mh.HomeItemInfo[itemId].ItemNum += num
	} else {
		mh.HomeItemInfo[itemId] = &HomeItem{ItemId: itemId, ItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得家具", config.ItemName, "----数量：", num, "----当前数量：", mh.HomeItemInfo[itemId].ItemNum)
	}

}
