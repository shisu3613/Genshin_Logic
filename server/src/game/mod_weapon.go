package game

import (
	"fmt"
	"server/csvs"
)

type Weapon struct {
	WeaponId int
	KeyId    int
}

type ModWeapon struct {
	WeaponInfo map[int]*Weapon
	MaxKey     int //永远递增的编号
}

func (mw *ModWeapon) AddItem(itemId int, num int) {
	//武器表的验证环节
	config := csvs.GetWeaponConfig(itemId)
	if config == nil {
		fmt.Println("非法武器")
		return
	}
	if len(mw.WeaponInfo)+num > csvs.MaxWeaponSize {
		fmt.Println("武器背包已满！！！")
		return
	}
	for i := 0; i < num; i++ {
		weapon := new(Weapon)
		weapon.WeaponId = itemId
		mw.MaxKey++
		weapon.KeyId = mw.MaxKey
		mw.WeaponInfo[weapon.KeyId] = weapon
		fmt.Println("获得武器：", csvs.GetItemName(itemId), "------武器星级：", config.Star, "-----武器编号：", weapon.KeyId)
	}
}

func (mw *ModWeapon) RemoveItem(weaponId int) {
	if _, ok := mw.WeaponInfo[weaponId]; !ok {
		fmt.Println("当前编号不存在")
		return
	}
	fmt.Println("移除武器id为：", weaponId, "武器名称为", csvs.GetItemName(mw.WeaponInfo[weaponId].WeaponId))
	delete(mw.WeaponInfo, weaponId)
}
