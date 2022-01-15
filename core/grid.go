package core

import (
	"fmt"
	"sync"
)

//一个AOI地图中的格子类型

type Grid struct {

	//格子ID
	GID int
	//格子的边际坐标
	MinX int
	MaxX int
	MinY int
	MaxY int

	//当前格子内成员的ID集合
	playsIDs map[int]bool

	//保护当前集合的锁
	playsIDsLock sync.RWMutex
}

// NewGrid 初始化当前格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:          gID,
		MinX:         minX,
		MaxX:         maxX,
		MinY:         minY,
		MaxY:         maxY,
		playsIDs:     make(map[int]bool),
		playsIDsLock: sync.RWMutex{},
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.playsIDsLock.Lock()
	defer g.playsIDsLock.Unlock()

	g.playsIDs[playerID] = true
}

// Remove 给当前格子删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.playsIDsLock.Lock()
	defer g.playsIDsLock.Unlock()

	delete(g.playsIDs, playerID)
}

// GetPlayerIDs 得到格子中所有玩家的ID
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.playsIDsLock.RLock()
	defer g.playsIDsLock.RUnlock()

	for k, _ := range g.playsIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 调试用，打印格子的基本信息
func (g Grid) String() string {
	return fmt.Sprintf("Grid id :%d,minX :%d,maxX :%d,,minY :%d,maxY :%d,playIds:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playsIDs)
}
