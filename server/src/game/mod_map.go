package game

import (
	"errors"
	"fmt"
	"math/rand"
	"server/csvs"
	"time"
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
	////当前物品适量不足
	//if !player.ModBag.HasEnoughItem(eventConfig.CostItem, eventConfig.CostNum) {
	//	name := csvs.GetItemName(eventConfig.CostItem)
	//	return errors.New(fmt.Sprintf("%s不足!,请从背包系统中通过物品id:%d增加%s", name, eventConfig.CostItem, name))
	//}
	//根据玩家刷新(秘境地图)
	//秘境地图的刷新之前已经做了，这里是更为具体的在完成所有之前的事件前无法领取奖励的机制
	if mm.MapInfo[mapID].MapType == csvs.RefreshPlayer && eventConfig.EventType == csvs.EventTypeReward {
		for _, v := range mm.MapInfo[mapID].EventInfo {
			eventConfigNow := csvs.GetEventConfig(v.EventId)
			if eventConfigNow == nil || eventConfigNow.EventType != csvs.EventTypeNormal {
				continue
			}
			if v.State != csvs.EventEnd {
				return fmt.Errorf("领取奖励前有事件尚未完成:%d", v.EventId)
			}
		}
	}

	//根据系统刷新（大地图上的事件）
	//mm.MapInfo[mapID].EventInfo[eventId].State = eventState
	if eventState == csvs.EventEnd {
		//先扣除需要消耗的物品
		if err := player.ModBag.RemoveItem(eventConfig.CostItem, eventConfig.CostNum, player); err != nil {
			return err
		}
		//变成end就是说明东西都已经拾取了或者东西没拾取便离开地图了
		for i := 0; i < eventConfig.EventDropTimes; i++ {
			DropItemConfigs := csvs.GetItemDrop(eventConfig.EventDrop) //生成需要拾取的物品config map
			if len(DropItemConfigs) == 0 {
				continue
			}
			for _, v := range DropItemConfigs {
				randNum := (rand.Intn(v.ItemNumMax-v.ItemNumMin+1) + v.ItemNumMin) * (player.GetMod(ModPlay).(*ModPlayer).WorldLevelNow*v.WorldAdd + csvs.DropWeightAll) / csvs.DropWeightAll
				mm.MapInfo[mapID].DropItemsOnMap[v.ItemId] += randNum
			}
		}
		//这里是直接得到的东西，比如蒲公英种子，宝箱里的原石
		for i := 0; i < eventConfig.EventGainTime; i++ {
			DropItemConfigs := csvs.GetItemDrop(eventConfig.EventGain) //生成需要拾取的物品config map
			if len(DropItemConfigs) == 0 {
				continue
			}
			for _, v := range DropItemConfigs {
				randNum := (rand.Intn(v.ItemNumMax-v.ItemNumMin+1) + v.ItemNumMin) * (player.GetMod(ModPlay).(*ModPlayer).WorldLevelNow*v.WorldAdd + csvs.DropWeightAll) / csvs.DropWeightAll
				player.ModBag.AddItem(v.ItemId, int64(randNum), player)
			}
		}
		fmt.Printf("当前事件%s交互完成\n", eventConfig.Name)
		//拾取物品更新刷新时间
		switch eventConfig.RefreshType {
		case csvs.MapRefreshCant:
			//永远不刷新
			//return nil
		case csvs.MapRefreshTwoDay:
			//48小时刷新，稀有植物,为了测试这里显示两分钟刷新一次
			mm.MapInfo[mapID].EventInfo[eventId].NextResetTime = time.Now().Add(time.Minute * 2).Unix()
		case csvs.MapRefreshHalfDay:
			mm.MapInfo[mapID].EventInfo[eventId].NextResetTime = time.Now().Add(time.Second * 30).Unix()
		case csvs.MapRefreshThreeDay:
			mm.MapInfo[mapID].EventInfo[eventId].NextResetTime = time.Now().Add(time.Minute * 3).Unix()
		case csvs.MapRefreshWeek:
			//一星期的刷新机制和其他的都不一样，是打完之后下星期固定时间刷新的
			//这里设置为下个小时的开始
			currentTime := (time.Now().Unix()/(60*60) + 1) * (60 * 60)
			mm.MapInfo[mapID].EventInfo[eventId].NextResetTime = currentTime
		}
	}
	//当奖励领取完后，更新地图状态退出秘境
	mm.MapInfo[mapID].EventInfo[eventId].State = eventState
	if eventConfig.EventType == csvs.EventTypeReward {
		return errors.New("退出当前秘境")
	}
	return nil
}

// RefreshWhenCome 进入地图时的刷新机制检查刷新状态
func (mm *ModMap) RefreshWhenCome(mapId int) {
	curMap, ok := mm.MapInfo[mapId]
	if !ok {
		fmt.Printf("当前地图不存在")
		return
	}
	//当前并不是秘境地图，根据时间判断更新
	if curMap.MapType != csvs.RefreshPlayer {
		for _, v := range curMap.EventInfo {
			if v.NextResetTime <= time.Now().Unix() {
				v.State = csvs.EventStart
			}
		}
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
		//fmt.Println("testPoint")
		return
	}
	var action int
	fmt.Println("当前地图上有未拾取的物品，按1拾取(模拟一次性拾取动作)")
	fmt.Scan(&action)
	switch action {
	case 1:
		for k, v := range MapDrop {
			player.ModBag.AddItemToBag(k, int64(v))
			delete(MapDrop, k)
		}

	default:
		return
	}
}

// GetEventList 生成当前可选事件列表
func (mm *ModMap) GetEventList(mapId int) {
	curMap := mm.MapInfo[mapId]
	for _, v := range curMap.EventInfo {
		if v.State == csvs.EventStart {
			fmt.Printf("当前事件%s,状态可以交互,触发事件Id为%d\n", csvs.GetEventName(v.EventId), v.EventId)
		} else {
			if csvs.GetEventConfig(v.EventId).RefreshType == csvs.MapRefreshCant {
				fmt.Printf("当前事件%s,状态不可以交互,永不刷新\n", csvs.GetEventName(v.EventId))
			} else {
				fmt.Printf("当前事件%s,状态不可以交互,下次刷新时间为%s\n", csvs.GetEventName(v.EventId), getTimeForm(v.NextResetTime))
			}
		}
	}
}

func getTimeForm(strTime int64) string {
	//记12345,3那个位置的数这里我使用的15，也就是用24小时格式来显示，如果直接写03则是12小时am pm格式。
	timeLayout := "2006-01-02 15:04:05"
	datetime := time.Unix(strTime, 0).Format(timeLayout)
	return datetime
}
