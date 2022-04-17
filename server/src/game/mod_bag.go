package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	DB "server/DB/GORM"
	"server/csvs"
)

type ItemInfo struct {
	ItemId  int
	ItemNum int64
}

type ModBag struct {
	BagInfo map[int]*ItemInfo
	player  *Player
}

// AddItem @Modified By WangYuding 2022/4/17 17:43:00
// @Modified description 删除player参数，使用装饰模式
func (mb *ModBag) AddItem(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	//case csvs.ITEMTYPE_NORMAL:
	//	mb.AddItemToBag(itemId, num)
	case csvs.ItemTypeRole:
		mb.player.ModRole.AddItem(itemId, num)
	case csvs.ItemTypeIcon:
		mb.player.GetMod(IconMod).(*ModIcon).AddItem(itemId)
	case csvs.ItemTypeCard:
		mb.player.ModCard.AddItem(itemId, 12)
	case csvs.ItemTypeWeapon:
		mb.player.ModWeapon.AddItem(itemId, int(num))
	case csvs.ItemTypeRelic:
		mb.player.ModRelic.AddItem(itemId, int(num))
	case csvs.ItemTypeCook:
		mb.player.ModCook.AddItem(itemId)
	case csvs.ItemTypeCookBook:
		if num > 1 {
			fmt.Println("注意：只能有一份食谱！")
		}
		mb.AddItemToBag(itemId, 1)
	case csvs.ItemTypeFurn: //家园模块，家具识别
		mb.player.ModHome.AddItem(itemId, num)

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

func (mb *ModBag) RemoveItem(itemId int, num int64) error {
	if itemId == 0 {
		return nil
	}
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
		return errors.New(fmt.Sprint(itemConfig.ItemName, "数量不足", "----当前数量：", nowNum, ",请通过背包系统物品Id:", itemId, "增加物品"))
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

func (mb *ModBag) UseItem(itemId int, num int64) {
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
		mb.UseCookBook(itemId, num)

	default: //目前上是可以放进背包里面的物品
		//mb.AddItemToBag(itemId, num)
		fmt.Println("此物品无法使用")
		return
	}
}

func (mb *ModBag) UseCookBook(itemId int, num int64) {
	CookBook := csvs.GetCookBookConfig(itemId)
	if CookBook == nil {
		fmt.Println("食谱不存在")
		return
	}
	if err := mb.RemoveItem(itemId, num); err != nil {
		fmt.Println(err)
		return
	}
	mb.AddItem(CookBook.Reward, num)
}

// LoadData
// @Description: 当前所有模块中最重要的背包模块的LoadData,
//粗糙处理，暂时使用json格式与数据库交互
// @receiver mb
func (mb *ModBag) LoadData() {
	pid, err := mb.player.Conn.GetProperty("PID")
	uid := pid.(int) + 100000000
	if err != nil {
		mb.player.SendStringMsg(800, "意外错误，请重新输入id")
	}
	var test DBModBag
	if errors.Is(DB.GormDB.First(&test, "user_id", uid).Error, gorm.ErrRecordNotFound) {
		//fmt.Println(DB.GormDB.Find(&test, "user_id", uid).Error)
		fmt.Println("No Icon map, create new record")
		content, _ := json.Marshal(mb)
		tmp := DBModBag{
			UserId:   uid,
			JsonData: content,
		}
		DB.GormDB.Create(&tmp)
	} else {
		configFile := test.JsonData
		err = json.Unmarshal(configFile, &mb)
		if err != nil {
			fmt.Println("Bag json empty")
			return
		}
	}
}

// SaveData
// @Description: 同IconMod，可能有问题
// @receiver mb
func (mb *ModBag) SaveData() {
	uid, _ := mb.player.GetUserID()
	uid += 100000000
	content, _ := json.Marshal(mb)
	var test DBModBag
	DB.GormDB.Find(&test, "user_id", uid)
	test.JsonData = content
	DB.GormDB.Save(test)
}

func (mb *ModBag) init(player *Player) {
	mb.player = player
	mb.BagInfo = make(map[int]*ItemInfo)
}
