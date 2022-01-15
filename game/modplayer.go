package game

import "fmt"

type ModPlayer struct {
	UserID         int
	Icon           int //头像模块
	Card           int
	Name           string
	Sign           string
	PlayerLevel    int
	PlayerExp      int
	WorldLevel     int
	WorldLevelCool int64 //世界等级冷却时间
	Birth          int
	ShowTeam       []int //展示的页面的队伍人物ID
	ShowCard       []int //名片ID
	//看不见的字段：远比上面多
	IsProhibit int //是否封号
	IsGM       int //是否是GM号
}

// SetIcon 对内的逻辑接口，收到客户端发送的头像id
func (mp *ModPlayer) SetIcon(iconId int, pl *Player) {
	if !pl.ModIcon.IsHasIcon(iconId) {
		//通知客户端操作非法
		return
	}
	pl.ModPlayer.Icon = iconId
	fmt.Println("当前图标：", pl.ModPlayer.Icon)
}

// SetCard 对内的逻辑接口，收到客户端发送的名片id
func (mp *ModPlayer) SetCard(iconId int, pl *Player) {
	if !pl.ModCard.IsHasCard(iconId) {
		//通知客户端操作非法
		return
	}
	pl.ModPlayer.Icon = iconId
	fmt.Println("当前名片：", pl.ModPlayer.Icon)
}

// SetName 对内的逻辑接口，收到客户端发送的名字的字符串
func (mp *ModPlayer) SetName(sign string, pl *Player) {
	//名字的处理模块
	//对于文字的违禁字处理，正则表达式处理
	//主要做法有两种：
	//1.用外部处理违禁字，调用HTTP地址接口
	//2。使用管理类公共模块(单例模式 )
	if GetManageBanWord().IsBanWord(sign) {
		return
	}
	pl.ModPlayer.Name = sign
	fmt.Println("当前名字：", pl.ModPlayer.Name)
}

// SetSign 对内的逻辑接口，收到客户端发送的签名的字符传
func (mp *ModPlayer) SetSign(name string, pl *Player) {

	if GetManageBanWord().IsBanWord(name) {
		return
	}
	pl.ModPlayer.Name = name
	fmt.Println("当前签名：", pl.ModPlayer.Name)
}
