package game

import (
	"fmt"
	"server/csvs"
)

type WishPool struct {
	PoolId int

	//为了区分这是测试用保底统计
	FiveStarTimesTest int
	FourStarTimesTest int

	//保底设计
	FiveStarTimes int
	FourStarTimes int

	//真正抽卡统计信息
	StatFiveTotal   int
	StatFourRole    int
	StatFourWeapon  int
	StatTotalWishes int
}

type ModWish struct {
	UPWishPool     *WishPool
	NormalWishPool *WishPool
}

// DoPoolTest 抽卡模拟
func (wish *ModWish) DoPoolTest(times int) (res string) {
	result := make(map[int]int)
	for i := 0; i < times; i++ {
		dropGroup := csvs.ConfigDropGroupMap[1000] //drop id 为1000
		if dropGroup == nil {
			return
		}
		wish.UPWishPool.FourStarTimesTest++
		wish.UPWishPool.FiveStarTimesTest++
		if wish.UPWishPool.FiveStarTimesTest > csvs.FiveStarLimit || wish.UPWishPool.FourStarTimesTest > csvs.FourStarLimit {
			NewDropGroup := new(csvs.DropGroup)
			NewDropGroup.DropId = dropGroup.DropId
			NewDropGroup.WeightAll = dropGroup.WeightAll
			//五星权值变动,先减三星再减去四星,五星优先级更高
			addFiveWeight, remainWeight := (wish.UPWishPool.FiveStarTimesTest-csvs.FiveStarLimit)*csvs.FiveStarLimitIncrement, 0
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//四星保底权重设置
			addFourWeight := (wish.UPWishPool.FourStarTimesTest - csvs.FourStarLimit) * csvs.FourStarLimitIncrement
			//三星和五星总和权重
			FiveThreeWeight := 0

			if addFourWeight < 0 {
				addFourWeight = 0
			} else if addFourWeight > 10000 {
				addFourWeight = 10000
			}
			//修改参数设置
			fourStarConfig := new(csvs.ConfigWishes)
			for _, v := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigWishes)
				newConfig.IsEnd = v.IsEnd
				newConfig.Result = v.Result
				newConfig.DropId = v.DropId
				if v.Result == 10001 { //即当目前设置是五星的设置时
					newConfig.Weight = v.Weight + addFiveWeight
					NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, newConfig)
					FiveThreeWeight += newConfig.Weight
				} else if v.Result == 10003 {
					if v.Weight >= addFiveWeight {
						newConfig.Weight += v.Weight - addFiveWeight
					} else {
						remainWeight = addFiveWeight - v.Weight
						newConfig.Weight = 0
					}
					if newConfig.Weight > addFourWeight {
						newConfig.Weight -= addFourWeight
					} else {
						newConfig.Weight = 0
					}
					NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, newConfig)
					FiveThreeWeight += newConfig.Weight
				} else {
					fourStarConfig = newConfig
				}
			}
			fourStarConfig.Weight = csvs.TotalWishWeight - FiveThreeWeight - remainWeight
			//fmt.Println(fourStarConfig.Weight, FiveThreeWeight, remainWeight)
			NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, fourStarConfig)
			dropGroup = NewDropGroup
		}
		config := csvs.GetRandDrop(dropGroup)
		if config != nil {
			//抽到五星后归零
			roleConfig := csvs.GetRoleConfig(config.Result)
			//如果不是武器
			if roleConfig != nil && roleConfig.Star == 5 {
				//抽到五星了
				wish.UPWishPool.FiveStarTimesTest = 0
			} else {
				WeaponConfig := csvs.GetWeaponConfig(config.Result)
				if WeaponConfig == nil || WeaponConfig.Star == 4 {
					wish.UPWishPool.FourStarTimesTest = 0
				}
			}
			result[config.Result]++
		}
	}
	FiveStar, FourStarRole, FourStarWeapon, ThreeStar := 0, 0, 0, 0
	for k, v := range result {
		fmt.Printf("抽中%s次数：%d\n", csvs.GetItemName(k), v)
		res += fmt.Sprintf("抽中%s次数：%d\n", csvs.GetItemName(k), v)
		if csvs.GetItemConfig(k).SortType == csvs.ItemTypeRole {
			if csvs.GetRoleConfig(k).Star == 4 {
				FourStarRole += v
			} else {
				FiveStar += v
			}
		} else {
			if csvs.GetWeaponConfig(k).Star == 3 {
				ThreeStar += v
			} else {
				FourStarWeapon += v
			}
		}
	}
	fmt.Printf("本次您一共进行了%d次祈愿，共获得五星角色%d位，占总数的%.4f%%,四星角色%d位，四星武器%d把，四星综合概率为%.4f%%\n", times, FiveStar, 100*float32(FiveStar)/float32(times), FourStarRole,
		FourStarWeapon, 100*float32(FourStarRole+FourStarWeapon)/float32(times))
	res += fmt.Sprintf("本次您一共进行了%d次祈愿，共获得五星角色%d位，占总数的%.4f%%,四星角色%d位，四星武器%d把，四星综合概率为%.4f%%\n", times, FiveStar, 100*float32(FiveStar)/float32(times), FourStarRole,
		FourStarWeapon, 100*float32(FourStarRole+FourStarWeapon)/float32(times))
	return res
}

//DoPool UP池子抽卡
func (wish *ModWish) DoPool(times int, player *Player) {
	result := make(map[int]int)
	for i := 0; i < times; i++ {
		dropGroup := csvs.ConfigDropGroupMap[1000] //drop id 为1000
		if dropGroup == nil {
			fmt.Println("数据错误，请检查配置表")
			return
		}
		wish.UPWishPool.FourStarTimes++
		wish.UPWishPool.FiveStarTimes++
		if wish.UPWishPool.FiveStarTimes > csvs.FiveStarLimit || wish.UPWishPool.FourStarTimes > csvs.FourStarLimit {
			NewDropGroup := new(csvs.DropGroup)
			NewDropGroup.DropId = dropGroup.DropId
			NewDropGroup.WeightAll = dropGroup.WeightAll
			//五星权值变动,先减三星再减去四星,五星优先级更高
			addFiveWeight, remainWeight := (wish.UPWishPool.FiveStarTimes-csvs.FiveStarLimit)*csvs.FiveStarLimitIncrement, 0
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//四星保底权重设置
			addFourWeight := (wish.UPWishPool.FourStarTimes - csvs.FourStarLimit) * csvs.FourStarLimitIncrement
			//三星和五星总和权重
			FiveThreeWeight := 0

			if addFourWeight < 0 {
				addFourWeight = 0
			} else if addFourWeight > 10000 {
				addFourWeight = 10000
			}
			//修改参数设置
			fourStarConfig := new(csvs.ConfigWishes)
			for _, v := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigWishes)
				newConfig.IsEnd = v.IsEnd
				newConfig.Result = v.Result
				newConfig.DropId = v.DropId
				if v.Result == 10001 { //即当目前设置是五星的设置时
					newConfig.Weight = v.Weight + addFiveWeight
					NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, newConfig)
					FiveThreeWeight += newConfig.Weight
				} else if v.Result == 10003 {
					if v.Weight >= addFiveWeight {
						newConfig.Weight += v.Weight - addFiveWeight
					} else {
						remainWeight = addFiveWeight - v.Weight
						newConfig.Weight = 0
					}
					if newConfig.Weight > addFourWeight {
						newConfig.Weight -= addFourWeight
					} else {
						newConfig.Weight = 0
					}
					NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, newConfig)
					FiveThreeWeight += newConfig.Weight
				} else {
					fourStarConfig = newConfig
				}
			}
			fourStarConfig.Weight = csvs.TotalWishWeight - FiveThreeWeight - remainWeight
			//fmt.Println(fourStarConfig.Weight, FiveThreeWeight, remainWeight)
			NewDropGroup.DropConfigs = append(NewDropGroup.DropConfigs, fourStarConfig)
			dropGroup = NewDropGroup
		}
		config := csvs.GetRandDrop(dropGroup)
		if config != nil {
			//抽到五星后归零
			roleConfig := csvs.GetRoleConfig(config.Result)
			//如果不是武器
			if roleConfig != nil && roleConfig.Star == 5 {
				//抽到五星了
				wish.UPWishPool.FiveStarTimes = 0
			} else {
				WeaponConfig := csvs.GetWeaponConfig(config.Result)
				if WeaponConfig == nil || WeaponConfig.Star == 4 {
					wish.UPWishPool.FourStarTimes = 0
				}
			}
			player.ModBag.AddItemToBag(config.Result, 1)
			result[config.Result]++
		}
	}
	FiveStar, FourStarRole, FourStarWeapon, ThreeStar := 0, 0, 0, 0
	for k, v := range result {
		fmt.Printf("抽中%s次数：%d\n", csvs.GetItemName(k), v)
		if csvs.GetItemConfig(k).SortType == csvs.ItemTypeRole {
			if csvs.GetRoleConfig(k).Star == 4 {
				FourStarRole += v
			} else {
				FiveStar += v
			}
		} else {
			if csvs.GetWeaponConfig(k).Star == 3 {
				ThreeStar += v
			} else {
				FourStarWeapon += v
			}
		}
	}
	wish.UPWishPool.StatTotalWishes += times
	wish.UPWishPool.StatFiveTotal += FiveStar
	wish.UPWishPool.StatFourRole += FourStarRole
	wish.UPWishPool.StatFourWeapon += FourStarWeapon
	fmt.Printf("本次您一共进行了%d次祈愿，共获得五星角色%d位，占总数的%.4f%%,四星角色%d位，四星武器%d把，四星物品占总数的%.4f%%\n当前您的五星保底为%d抽，四星保底为%d抽\n", times, FiveStar, 100*float32(FiveStar)/float32(times), FourStarRole,
		FourStarWeapon, 100*float32(FourStarRole+FourStarWeapon)/float32(times), wish.UPWishPool.FiveStarTimes, wish.UPWishPool.FourStarTimes)
}
