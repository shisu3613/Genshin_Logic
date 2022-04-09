package game

import (
	"fmt"
	"server/csvs"
)

type RoleInfo struct {
	RoleId   int
	GetTimes int
	//等级 经验 圣遗物 等数据
}

type ModRole struct {
	RoleInfo map[int]*RoleInfo
}

func (mr *ModRole) IsHasRole(roleId int) bool {
	return true
}

func (mr *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

func (mr *ModRole) AddItem(roleId int, num int64, player *Player) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}
	for i := 0; i < int(num); i++ {
		_, ok := mr.RoleInfo[roleId]
		if !ok {
			data := new(RoleInfo)
			data.RoleId = roleId
			data.GetTimes = 1
			mr.RoleInfo[roleId] = data
		} else {
			//超过七命的时候角色是有25个星辉，应该是吧
			mr.RoleInfo[roleId].GetTimes++
			if mr.RoleInfo[roleId].GetTimes >= csvs.AddRoleTimeNormalMin &&
				mr.RoleInfo[roleId].GetTimes <= csvs.AddRoleTimeNormalMax {
				player.ModBag.AddItemToBag(config.Stuff, config.StuffNum)
				player.ModBag.AddItemToBag(config.StuffItem, config.StuffItemNum)
			} else {
				player.ModBag.AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
			}
		}
	}
	itemConfig := csvs.GetItemConfig(roleId)
	if itemConfig != nil {
		fmt.Println("获得角色", itemConfig.ItemName, "------总计", mr.RoleInfo[roleId].GetTimes, "次")
	}
	player.GetMod(IconMod).(*ModIcon).CheckGetIcon(roleId)
	player.ModCard.CheckGetCard(roleId, 10) //友好度系统不完善
}
