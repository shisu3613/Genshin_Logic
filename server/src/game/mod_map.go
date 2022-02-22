package game

import (
	"errors"
	"fmt"
	"math/rand"
	"server/csvs"
)

type Map struct {
	MapId          int
	MapName        string
	MapType        int
	EventInfo      map[int]*Event
	DropItemsOnMap map[int]int
}

type Event struct {
	EventId       int
	State         int
	NextResetTime int64
}

type ModMap struct {
	MapInfo map[int]*Map //地图模块的内容概述，主地图，
}

func (mm *ModMap) InitData() {
	mm.MapInfo = make(map[int]*Map)
	//先将地图模块初始化
	for _, v := range csvs.ConfigMapMap {
		_, ok := mm.MapInfo[v.MapId]
		if !ok {
			mm.MapInfo[v.MapId] = mm.NewMapInfo(v)
		}
	}

	//再初始化每张地图内部有的事件模块
	for _, v := range csvs.ConfigMapEventMap {
		_, ok := mm.MapInfo[v.MapId]
		if !ok {
			fmt.Printf("当前交互事件%s所处地图模块尚未开放\n", v.Name)
			continue
		}
		_, ok = mm.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			mm.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			mm.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			mm.MapInfo[v.MapId].EventInfo[v.EventId].State = csvs.EventStart
		}

	}
}

func (mm *ModMap) NewMapInfo(configMap *csvs.ConfigMap) *Map {
	mapinfo := new(Map)
	mapinfo.MapId = configMap.MapId
	mapinfo.MapName = configMap.MapName
	mapinfo.MapType = configMap.MapType
	mapinfo.EventInfo = make(map[int]*Event)
	mapinfo.DropItemsOnMap = make(map[int]int)
	return mapinfo
}

// SetEventState 变更状态，返回error和如果有掉落物获得该掉落物
func (mm *ModMap) SetEventState(mapID int, eventId int, eventState int, player *Player) error {
	_, ok := mm.MapInfo[mapID]
	if !ok {
		return errors.New("地图不存在")
	}
	_, ok = mm.MapInfo[mapID].EventInfo[eventId]
	if !ok {
		return errors.New("事件不存在")
	}
	if mm.MapInfo[mapID].EventInfo[eventId].State > eventState {
		return errors.New("状态异常（状态减少）")
	}
	eventConfig := csvs.GetEventConfig(eventId)
	if eventConfig == nil {
		return errors.New("配置表中事件不存在")
	}
	if !player.ModBag.HasEnoughItem(eventConfig.CostItem, eventConfig.CostNum) {
		return errors.New(fmt.Sprintf("%s不足!", csvs.GetItemName(eventConfig.CostItem)))
	}
	//根据玩家刷新(秘境地图)

	//根据系统刷新（大地图上的事件）
	mm.MapInfo[mapID].EventInfo[eventId].State = eventState
	if eventState == csvs.EventEnd {
		//变成end就是说明东西都已经拾取了或者东西没拾取便离开地图了
		for i := 1; i < eventConfig.EventDropTimes; i++ {
			DropItemConfigs := csvs.GetItemDrop(eventConfig.EventDrop) //生成需要拾取的物品config map
			if len(DropItemConfigs) == 0 {
				continue
			}
			for _, v := range DropItemConfigs {
				randNum := (rand.Intn(v.ItemNumMax-v.ItemNumMin+1) + v.ItemNumMin) * (player.ModPlayer.WorldLevelNow*v.WorldAdd + csvs.DropWeightAll) / csvs.DropWeightAll
				mm.MapInfo[mapID].DropItemsOnMap[v.ItemId] += randNum
			}
		}
		//这里是直接得到的东西，比如蒲公英种子，宝箱里的原石
		for i := 1; i < eventConfig.EventGainTime; i++ {
			DropItemConfigs := csvs.GetItemDrop(eventConfig.EventGain) //生成需要拾取的物品config map
			if len(DropItemConfigs) == 0 {
				continue
			}
			for _, v := range DropItemConfigs {
				randNum := (rand.Intn(v.ItemNumMax-v.ItemNumMin+1) + v.ItemNumMin) * (player.ModPlayer.WorldLevelNow*v.WorldAdd + csvs.DropWeightAll) / csvs.DropWeightAll
				player.ModBag.AddItem(v.ItemId, int64(randNum), player)
			}
		}
		fmt.Printf("当前事件%s交互完成", eventConfig.Name)
	}
	return nil
}

func (mm *ModMap) RefreshWhenCome(mapId int) {
	curMap, ok := mm.MapInfo[mapId]
	if !ok {
		fmt.Printf("当前地图不存在")
		return
	}
	//当前并不是秘境地图，返回
	if curMap.MapType != csvs.RefreshPlayer {
		return
	}
	curMap.DropItemsOnMap = make(map[int]int)
	for _, v := range curMap.EventInfo {
		v.State = csvs.EventStart
	}

}

//检查当前进入地图的时候有没有遗留物
func (mm *ModMap) checkAnyDropOnMap(mapId int, player *Player) {
	MapDrop := mm.MapInfo[mapId].DropItemsOnMap
	if len(MapDrop) == 0 {
		return
	}
	var action int
	fmt.Println("当前地图上有未拾取的物品，按1拾取(模拟一次性拾取动作)")
	switch action {
	case 1:
		for k, v := range MapDrop {
			player.ModBag.AddItemToBag(k, int64(v))
		}
	default:
		return
	}
}

//todo:事件时间刷新机制
