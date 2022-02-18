package csvs

import (
	"fmt"
	"math/rand"
)

var ConfigDropGroupMap map[int]*DropGroup

type DropGroup struct {
	DropId      int
	WeightAll   int
	DropConfigs []*ConfigWishes
}

func CheckLoadCsv() {
	//二次处理,更新表结构数组结构为map结构
	MakeDropGroupMap()

	fmt.Println("csv配置读取完成---ok")
}
func MakeDropGroupMap() {
	ConfigDropGroupMap = make(map[int]*DropGroup)
	for _, v := range ConfigWishesSlice {
		dropGroup, ok := ConfigDropGroupMap[v.DropId]
		if !ok {
			dropGroup = &DropGroup{
				DropId: v.DropId,
			}
			ConfigDropGroupMap[v.DropId] = dropGroup
		}
		dropGroup.WeightAll += v.Weight
		dropGroup.DropConfigs = append(dropGroup.DropConfigs, v)
	}
	fmt.Println("抽卡掉落模块数据结构加载完成")
	//RandDropTest()
	return
}

func RandDropTest() {
	dropGroup := ConfigDropGroupMap[1000] //drop id 为1000
	if dropGroup == nil {
		return
	}
	num := 0
	for {
		config := GetRandDrop(dropGroup)
		if config == nil {
			break
		} else {
			fmt.Println(GetItemName(config.Result))
			dropGroup = ConfigDropGroupMap[1000]
			num++
			if num == 100 {
				break
			}
		}
	}
}

func GetRandDrop(dropGroup *DropGroup) *ConfigWishes {
	//rand.Seed(time.Now().UnixNano())
	//time.Sleep(time.Millisecond * 100)
	//
	//for i := 0; i < 10; i++ {
	//	randNum := rand.Intn(dropGroup.WeightAll)
	//	fmt.Println(randNum)
	//}
	randNum := rand.Intn(dropGroup.WeightAll)
	//fmt.Println(randNum)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if randNum < randNow {
			if v.IsEnd == LogicTrue {
				return v
			}
			dropGroup = ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				fmt.Println(" 当前resultID:", v.Result, "不存在")
				return nil
			}
			return GetRandDrop(dropGroup)
		}
	}
	return nil
}
