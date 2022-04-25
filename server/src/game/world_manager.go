package game

import "sync"

type WorldManager struct {
	pLock   sync.RWMutex //保护玩家的互斥机制
	Players map[int]*Player
}

// WorldMgrObj 提供了一个对外的世界管理模块的句柄
var WorldMgrObj *WorldManager

//提供了worldManager初始化的方法
func init() {
	WorldMgrObj = &WorldManager{
		pLock:   sync.RWMutex{},
		Players: make(map[int]*Player),
	}
}

// AddPlayer 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	wm.pLock.Lock()
	UserID, _ := player.GetUserID()
	wm.Players[UserID] = player
	wm.pLock.Unlock()
}

// RemovePlayerByPID 从玩家信息表中移除一个玩家
func (wm *WorldManager) RemovePlayerByPID(pID int) {
	wm.pLock.Lock()
	delete(wm.Players, pID)
	wm.pLock.Unlock()
}

// GetPlayerByPID 通过玩家ID 获取对应玩家信息
func (wm *WorldManager) GetPlayerByPID(pID int) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	return wm.Players[pID]
}

// GetAllPlayers 获取所有玩家的信息
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range wm.Players {
		players = append(players, v)
	}

	//返回
	return players
}
