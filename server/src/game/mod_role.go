package game

import (
	"fmt"
	"server/csvs"
)

type RoleInfo struct {
	RoleId   int
	GetTimes int
	Level    int
	//等级 经验 圣遗物 等数据
}

type ModRole struct {
	RoleInfo map[int]*RoleInfo
	player   *Player
}

func (mr *ModRole) IsHasRole(roleId int) bool {
	_, ok := mr.RoleInfo[roleId]
	return ok
}

func (mr *ModRole) GetRoleLevel(roleId int) int {
	data, _ := mr.RoleInfo[roleId]
	return data.Level
}

func (mr *ModRole) AddItem(roleId int, num int64) {
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
			data.Level = 1
			mr.RoleInfo[roleId] = data
		} else {
			//超过七命的时候角色是有25个星辉，应该是吧
			mr.RoleInfo[roleId].GetTimes++
			if mr.RoleInfo[roleId].GetTimes >= csvs.AddRoleTimeNormalMin &&
				mr.RoleInfo[roleId].GetTimes <= csvs.AddRoleTimeNormalMax {
				mr.player.AddBagItem(config.Stuff, config.StuffNum)
				mr.player.AddBagItem(config.StuffItem, config.StuffItemNum)
			} else {
				mr.player.AddBagItem(config.MaxStuffItem, config.MaxStuffItemNum)
			}
		}
	}
	itemConfig := csvs.GetItemConfig(roleId)
	if itemConfig != nil {
		fmt.Println("获得角色", itemConfig.ItemName, "------总计", mr.RoleInfo[roleId].GetTimes, "次")
	}
	mr.player.GetMod(IconMod).(*ModIcon).CheckGetIcon(roleId)
	mr.player.ModCard.CheckGetCard(roleId, 10) //友好度系统不完善
}
