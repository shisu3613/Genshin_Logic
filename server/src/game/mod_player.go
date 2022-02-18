package game

import (
	"fmt"
	"server/csvs"
	"time"
)

type ShowRole struct {
	RoleId    int
	RoleLevel int
}
type ModPlayer struct {
	DBPlayer

	ShowCard *Cards      //展示名片
	ShowTeam []*ShowRole //展示阵容
}

type Cards struct {
	Cards []int
}

func (mp *ModPlayer) SetIcon(iconId int, player *Player) {
	if !player.ModIcon.IsHasIcon(iconId) {
		//通知客户端，操作非法
		fmt.Println("没有头像:", iconId)
		return
	}

	player.ModPlayer.Icon = iconId
	fmt.Println("变更头像为:", csvs.GetItemName(iconId), player.ModPlayer.Icon)
}

func (mp *ModPlayer) SetCard(cardId int, player *Player) {
	if !player.ModCard.IsHasCard(cardId) {
		//通知客户端，操作非法
		return
	}

	player.ModPlayer.Card = cardId
	fmt.Println("当前名片", player.ModPlayer.Card)
}

func (mp *ModPlayer) SetName(name string, player *Player) {
	if GetManageBanWord().IsBanWord(name) {
		return
	}

	player.ModPlayer.Name = name
	fmt.Println("设置成功,名字变更为:", player.ModPlayer.Name)
}

func (mp *ModPlayer) SetSign(sign string, player *Player) {
	if GetManageBanWord().IsBanWord(sign) {
		return
	}

	player.ModPlayer.Sign = sign
	fmt.Println("设置成功,签名变更为:", player.ModPlayer.Sign)
}

func (mp *ModPlayer) AddExp(exp int, player *Player) {
	mp.PlayerExp += exp
	for {
		config := csvs.GetNowLevelConfig(mp.PlayerLevel)
		if config == nil {
			break
		}
		if config.PlayerExp == 0 {
			break
		}
		//是否完成任务
		if config.ChapterId > 0 && !player.ModUniqueTask.IsTaskFinish(config.ChapterId) {
			break
		}
		if mp.PlayerExp >= config.PlayerExp {
			mp.PlayerLevel += 1
			mp.PlayerExp -= config.PlayerExp
		} else {
			break
		}
	}
	fmt.Println("当前等级:", mp.PlayerLevel, "---当前经验：", mp.PlayerExp)
}

func (mp *ModPlayer) ReduceWorldLevel(player *Player) {
	if mp.WorldLevel < csvs.ReduceWorldLevelStart {
		fmt.Println("操作失败:, ---当前世界等级：", mp.WorldLevel)
		return
	}

	if mp.WorldLevel-mp.WorldLevelNow >= csvs.ReduceWorldLevelMax {
		fmt.Println("操作失败:, ---当前世界等级：", mp.WorldLevel, "---真实世界等级：", mp.WorldLevelNow)
		return
	}

	if time.Now().Unix() < mp.WorldLevelCool {
		fmt.Println("操作失败:, ---冷却中")
		return
	}

	mp.WorldLevelNow -= 1
	mp.WorldLevelCool = time.Now().Unix() + csvs.ReduceWorldLevelCoolTime
	fmt.Println("操作成功:, ---当前世界等级：", mp.WorldLevel, "---真实世界等级：", mp.WorldLevelNow)
	return
}

func (mp *ModPlayer) ReturnWorldLevel(player *Player) {
	if mp.WorldLevelNow == mp.WorldLevel {
		fmt.Println("操作失败:, ---当前世界等级：", mp.WorldLevel, "---真实世界等级：", mp.WorldLevelNow)
		return
	}

	if time.Now().Unix() < mp.WorldLevelCool {
		fmt.Println("操作失败:, ---冷却中")
		return
	}

	mp.WorldLevelNow += 1
	mp.WorldLevelCool = time.Now().Unix() + csvs.ReduceWorldLevelCoolTime
	fmt.Println("操作成功:, ---当前世界等级：", mp.WorldLevel, "---真实世界等级：", mp.WorldLevelNow)
	return
}

// SetBirth 月份判断，已经设置过生日也要判断
func (mp *ModPlayer) SetBirth(birth int, player *Player) {
	if mp.Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}

	month := birth / 100
	day := birth % 100

	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if day <= 0 || day > 31 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 4, 6, 9, 11:
		if day <= 0 || day > 30 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 2:
		if day <= 0 || day > 29 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	default:
		fmt.Println("没有", month, "月！")
		return
	}

	mp.Birth = birth
	fmt.Println("设置成功，生日为:", month, "月", day, "日")

	if mp.IsBirthDay() {
		fmt.Println("今天是你的生日，生日快乐！") //定时的礼物代码 物品icon
	} else {
		fmt.Println("期待你生日的到来!")
	}

}

// IsBirthDay 当前服务器时间判断
func (mp *ModPlayer) IsBirthDay() bool {
	month := time.Now().Month()
	day := time.Now().Day()
	if int(month) == mp.Birth/100 && day == mp.Birth%100 {
		return true
	}
	return false
}

func (mp *ModPlayer) SetShowCard(showCard []int, player *Player) {

	if len(showCard) > csvs.ShowSize {
		return
	}

	cardExist := make(map[int]int)
	newList := make([]int, 0)
	for _, cardId := range showCard {
		_, ok := cardExist[cardId]
		if ok {
			continue
		}
		if !player.ModCard.IsHasCard(cardId) { //判断玩家有没有这个名片
			continue
		}
		newList = append(newList, cardId) //切片来保证数据
		cardExist[cardId] = 1
	}
	mp.ShowCard.Cards = newList
	fmt.Println(mp.ShowCard.Cards)
}

func (mp *ModPlayer) SetShowTeam(showRole []int, player *Player) {
	if len(showRole) > csvs.ShowSize {
		fmt.Println("消息结构错误")
		return
	}

	roleExist := make(map[int]int)
	newList := make([]*ShowRole, 0)
	for _, roleId := range showRole {
		_, ok := roleExist[roleId]
		if ok {
			continue
		}
		if !player.ModRole.IsHasRole(roleId) {
			continue
		}
		showRole := new(ShowRole)
		showRole.RoleId = roleId
		showRole.RoleLevel = player.ModRole.GetRoleLevel(roleId)
		newList = append(newList, showRole)
		roleExist[roleId] = 1
	}
	mp.ShowTeam = newList
	fmt.Println(mp.ShowCard)
}

func (mp *ModPlayer) SetHideShowTeam(isHide int, player *Player) {
	if isHide != csvs.LogicFalse && isHide != csvs.LogicTrue {
		return
	}
	mp.HideShowTeam = isHide
}

func (mp *ModPlayer) SetProhibit(prohibit int) {
	mp.Prohibit = prohibit
}

func (mp *ModPlayer) SetIsGM(isGm int) { //布尔值尽量用int来代替
	mp.IsGM = isGm
}

func (mp *ModPlayer) IsCanEnter() bool {
	return int64(mp.Prohibit) < time.Now().Unix()
}