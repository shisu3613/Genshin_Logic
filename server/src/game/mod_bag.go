package game

import (
	"errors"
	"fmt"
	"server/csvs"
)

type ItemInfo struct {
	ItemId  int
	ItemNum int64
}

type ModBag struct {
	BagInfo map[int]*ItemInfo
}

func (mb *ModBag) AddItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	mb.AddItemToBag(itemId, num)
	case csvs.ItemTypeRole:
		player.ModRole.AddItem(itemId, num, player)
	case csvs.ItemTypeIcon:
		player.ModIcon.AddItem(itemId)
	case csvs.ItemTypeCard:
		player.ModCard.AddItem(itemId, 12)
	case csvs.ItemTypeWeapon:
		player.ModWeapon.AddItem(itemId, int(num))
	case csvs.ItemTypeRelic:
		player.ModRelic.AddItem(itemId, int(num))
	case csvs.ItemTypeCook:
		player.ModCook.AddItem(itemId)
	case csvs.ItemTypeCookBook:
		if num > 1 {
			fmt.Println("注意：只能有一份食谱！")
		}
		mb.AddItemToBag(itemId, 1)
	case csvs.ItemTypeFurn: //家园模块，家具识别
		player.ModHome.AddItem(itemId, num, player)

	default: //目前上是可以放进背包里面的物品
		mb.AddItemToBag(itemId, num)
	}
}

func (mb *ModBag) AddItemToBag(itemId int, num int64) {
	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum += num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得物品", config.ItemName, "----数量：", num, "----当前数量：", mb.BagInfo[itemId].ItemNum)
	}

}

func (mb *ModBag) RemoveItem(itemId int, num int64, player *Player) error {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return errors.New("物品不存在")
	}

	switch itemConfig.SortType {
	//case csvs.ItemTypeNormal:
	//	mb.RemoveItemToBagGM(itemId, num)
	case csvs.ItemTypeRole:
		//fmt.Println("无法删除角色")
		return errors.New("无法删除角色")
	case csvs.ItemTypeIcon:
		//fmt.Println("无法删除头像")
		return errors.New("无法删除头像")
	case csvs.ItemTypeCard:
		//fmt.Println("无法删除卡片")
		return errors.New("无法删除卡片")
	case csvs.ItemTypeCook:
		//fmt.Println("无法删除烹饪技能")
		return errors.New("无法删除烹饪技能")
	//case csvs.ItemTypeWeapon: //目前情况下，武器的移除只能一个一个移除
	//	player.ModWeapon.RemoveItem(itemId)
	//case csvs.ItemTypeRelic:
	//	player.ModRelic.RemoveItem(itemId)
	default: //同普通
		return mb.RemoveItemToBag(itemId, num)
	}
}

func (mb *ModBag) RemoveItemToBagGM(itemId int, num int64) { //恶意退款，将钱变为负数
	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum -= num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("扣除物品", config.ItemName, "----数量：", num, "----当前数量：", mb.BagInfo[itemId].ItemNum)
	}
}

func (mb *ModBag) RemoveItemToBag(itemId int, num int64) error {
	itemConfig := csvs.GetItemConfig(itemId)

	if !mb.HasEnoughItem(itemId, num) {
		nowNum := int64(0)
		_, ok := mb.BagInfo[itemId]
		if ok {
			nowNum = mb.BagInfo[itemId].ItemNum
		}
		//fmt.Sprintln(itemConfig.ItemName, "数量不足", "----当前数量：", nowNum)
		return errors.New(fmt.Sprint(itemConfig.ItemName, "数量不足", "----当前数量：", nowNum))
	}

	_, ok := mb.BagInfo[itemId]
	if ok {
		mb.BagInfo[itemId].ItemNum -= num
	} else {
		mb.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	fmt.Println("扣除物品", itemConfig.ItemName, "----数量：", num, "----当前数量：", mb.BagInfo[itemId].ItemNum)
	return nil
}

func (mb *ModBag) HasEnoughItem(itemId int, num int64) bool {
	//注意事件模块设置
	if itemId == 0 {
		return true
	}
	_, ok := mb.BagInfo[itemId]
	if !ok {
		return false
	} else if mb.BagInfo[itemId].ItemNum < num {
		return false
	}
	return true
}

func (mb *ModBag) UseItem(itemId int, num int64, player *Player) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}
	if !mb.HasEnoughItem(itemId, num) {
		nowNum := int64(0)
		_, ok := mb.BagInfo[itemId]
		if ok {
			nowNum = mb.BagInfo[itemId].ItemNum
		}
		fmt.Println(itemConfig.ItemName, "数量不足", "----当前数量：", nowNum)
		return
	}

	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	mb.AddItemToBag(itemId, num)
	//case csvs.ItemTypeRole:
	//	player.ModRole.AddItem(itemId, num, player)
	//case csvs.ItemTypeIcon:
	//	player.ModIcon.AddItem(itemId)
	//case csvs.ItemTypeCard:
	//	player.ModCard.AddItem(itemId, 12)
	//case csvs.ItemTypeWeapon:
	//	player.ModWeapon.AddItem(itemId, int(num))
	//case csvs.ItemTypeRelic:
	//	player.ModRelic.AddItem(itemId, int(num))
	case csvs.ItemTypeCookBook:
		mb.UseCookBook(itemId, num, player)

	default: //目前上是可以放进背包里面的物品
		//mb.AddItemToBag(itemId, num)
		fmt.Println("此物品无法使用")
		return
	}
}

func (mb *ModBag) UseCookBook(itemId int, num int64, player *Player) {
	CookBook := csvs.GetCookBookConfig(itemId)
	if CookBook == nil {
		fmt.Println("食谱不存在")
		return
	}
	if err := mb.RemoveItem(itemId, num, player); err != nil {
		fmt.Println(err)
		return
	}
	mb.AddItem(CookBook.Reward, num, player)
}
