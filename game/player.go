package game

type Player struct {
	ModPlayer *ModPlayer
	ModIcon   *ModIcon
	ModCard   *ModCard
}

// NewTestPlayer  测试用功能，新建一个玩家
func NewTestPlayer() *Player {
	player := new(Player)
	player.ModPlayer = new(ModPlayer)
	player.ModIcon = new(ModIcon)
	//上面是模块的初始化
	//***********************************//
	//下面是一些写死的功能模块
	//player.ModPlayer.Icon = 0
	return player
}

// RecvSetIcon 对外接口，收到客户端发送的头像id,主要与客户端直接沟通
func (pl *Player) RecvSetIcon(iconId int) {
	pl.ModPlayer.SetIcon(iconId, pl)
}

// RecvName 对外接口，收到客户端发送的名片id,主要与客户端直接沟通
func (pl *Player) RecvName(name string) {
	pl.ModPlayer.SetName(name, pl)
}

// RecvSign 对外接口，收到客户端发送的名片id,主要与客户端直接沟通
func (pl *Player) RecvSign(sign string) {
	pl.ModPlayer.SetSign(sign, pl)
}

// RecvSetCard 对外接口，收到客户端发送的名片id,主要与客户端直接沟通
func (pl *Player) RecvSetCard(cardId int) {
	pl.ModPlayer.SetCard(cardId, pl)
}
