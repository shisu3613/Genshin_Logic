package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	DB "server/DB/GORM"
	"server/csvs"
)

type Icon struct {
	IconId int
}

type ModIcon struct {
	IconInfo map[int]*Icon
	player   *Player
}

func (self *ModIcon) IsHasIcon(iconId int) bool {
	_, ok := self.IconInfo[iconId]
	return ok
}

func (self *ModIcon) AddItem(itemId int) {
	_, ok := self.IconInfo[itemId]
	if ok {
		fmt.Println("已存在头像：", itemId)
		return
	}
	config := csvs.GetIconConfig(itemId)
	if config == nil {
		fmt.Println("非法头像：", itemId)
		return
	}
	self.IconInfo[itemId] = &Icon{IconId: itemId}
	fmt.Println("获得头像：", csvs.GetItemName(itemId))
}

func (self *ModIcon) CheckGetIcon(roleId int) {
	config := csvs.GetIconConfigByRoleId(roleId)
	if config == nil {
		return
	}
	self.AddItem(config.IconId)
}

// SaveData
// @Description 选择换成json模式存到数据库里面（初步设想：背包模块统一管理 ）
// @Author WangYuding 2022-04-09 21:42:10
func (self *ModIcon) SaveData() {
	//pid, err := self.player.Conn.GetProperty("PID")
	//uid := pid.(int) + 100000000
	//if err != nil {
	//	self.player.SendStringMsg(800, "意外错误，请重新输入id")
	//}
	uid := self.player.GetUserID()
	content, _ := json.Marshal(self)
	var test DBIcon
	DB.GormDB.Find(&test, "user_id", uid)
	test.IconMapData = content
	DB.GormDB.Save(test)
}

func (self *ModIcon) LoadData() {
	pid, err := self.player.Conn.GetProperty("PID")
	uid := pid.(int) + 100000000
	if err != nil {
		self.player.SendStringMsg(800, "意外错误，请重新输入id")
	}
	var test DBIcon

	if errors.Is(DB.GormDB.First(&test, "user_id", uid).Error, gorm.ErrRecordNotFound) {
		//fmt.Println(DB.GormDB.Find(&test, "user_id", uid).Error)
		fmt.Println("No Icon map, create new record")
		content, _ := json.Marshal(self)
		tmp := DBIcon{
			UserId:      uid,
			IconMapData: content,
		}
		DB.GormDB.Create(&tmp)
	} else {
		configFile := test.IconMapData
		err = json.Unmarshal(configFile, &self)
		//println("test Icon Result:", self.player.Conn.GetConnID())
		// @Modified By WangYuding 2022/4/9 22:24:00
		// @Modified description 一个奇怪的bug，来自于数据库没有迁移好，所以加了个err判断
		if err != nil {
			fmt.Println("Icon json empty")
			return
		}
	}
}

func (self *ModIcon) init(player *Player) {
	self.player = player
	self.IconInfo = make(map[int]*Icon)
}
