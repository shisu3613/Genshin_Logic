package game

import (
	"encoding/json"
	"fmt"
	"log"
	DB "server/DB/GORM"
	"server/csvs"
	"server/msgJson"
	"server/zinx/ziface"
)

const (
	TaskStateInit   = 0
	TaskStateDoing  = 1
	TaskStateFinish = 2
	ModPlay         = "player"
	IconMod         = "icon"
)

//var PIDGen int = 1    //用于生成玩家ID的计数器
//var IDLock sync.Mutex //保护PIDGen的互斥机制

// ModBase 泛型改写模块
type ModBase interface {
	// LoadData
	// @Description: load the data from mysql
	LoadData()
	// SaveData
	// @Description: when the mod finished, rewrite to database
	SaveData()
	// init
	// @Description: 通过装饰模式将指针赋给模块,初始化模块不与数据库交互
	// @param player
	init(player *Player)
}

type Player struct {
	//客户端连接功能当前
	Conn ziface.IConnection //当前玩家连接

	//ModPlayer     *ModPlayer //modplayer包含玩家的基本面板信息，UID即是玩家当前ID
	//ModIcon       *ModIcon //解耦：包含头像信息，链接数据库或者客户端本地缓存中的id和图片
	ModCard       *ModCard
	ModUniqueTask *ModUniqueTask //任务模块
	ModRole       *ModRole       //人物模块
	ModBag        *ModBag        //背包模块，核心模块
	ModWeapon     *ModWeapon     //武器背包模块
	ModRelic      *ModRelic      //初始化圣遗物模块
	ModCook       *ModCook       //初始化烹饪技能背包
	ModHome       *ModHome       //家园模块
	ModWish       *ModWish
	ModMap        *ModMap //地图逻辑模块

	modManage map[string]ModBase
}

// InitClientPlayer @Modified By WangYuding 2022/4/5 18:06:00
// @Modified description 功能进一步分割这个部分只管初始化play模块，与数据库无关
func InitClientPlayer(conn ziface.IConnection) *Player {
	player := new(Player)
	//绑定客户端连接
	player.Conn = conn

	////生成Player Uid
	//IDLock.Lock()
	//ID := PIDGen
	//PIDGen++
	//IDLock.Unlock()

	//player.GetMod(ModPlay).(*ModPlayer) = new(ModPlayer)
	//playerMod里面绑定了UID
	//player.GetMod(ModPlay).(*ModPlayer).UserId = ID

	//player.ModIcon = new(ModIcon)
	//player.ModIcon.IconInfo = make(map[int]*Icon)
	player.ModCard = new(ModCard)
	player.ModCard.CardInfo = make(map[int]*Card)
	player.ModUniqueTask = new(ModUniqueTask)
	player.ModUniqueTask.MyTaskInfo = make(map[int]*TaskInfo)
	//player.ModUniqueTask.Locker = new(sync.RWMutex)
	player.ModRole = new(ModRole)
	player.ModRole.RoleInfo = make(map[int]*RoleInfo)
	player.ModBag = new(ModBag)
	player.ModBag.BagInfo = make(map[int]*ItemInfo)
	player.ModWeapon = new(ModWeapon)
	player.ModWeapon.WeaponInfo = make(map[int]*Weapon)

	player.ModRelic = new(ModRelic)
	player.ModRelic.RelicInfo = make(map[int]*Relic)

	player.ModCook = new(ModCook)
	player.ModCook.CookInfo = make(map[int]*Cook)

	player.ModHome = new(ModHome)
	player.ModHome.HomeItemInfo = make(map[int]*HomeItem)

	player.ModMap = new(ModMap)
	player.ModMap.InitData()
	//抽卡掉落模块
	player.ModWish = new(ModWish)
	player.ModWish.UPWishPool = new(WishPool)
	player.ModWish.NormalWishPool = new(WishPool)

	//****************************************
	//player.GetMod(ModPlay).(*ModPlayer).PlayerLevel = 1
	//player.GetMod(ModPlay).(*ModPlayer).Name = "旅行者"
	//player.GetMod(ModPlay).(*ModPlayer).WorldLevel = 1
	//player.GetMod(ModPlay).(*ModPlayer).WorldLevelNow = 1
	//player.GetMod(ModPlay).(*ModPlayer).player = player
	//****************************************

	//******************泛型更新部分*********************
	//初始化管理模块
	player.modManage = make(map[string]ModBase)

	player.modManage = map[string]ModBase{
		ModPlay: new(ModPlayer),
		IconMod: new(ModIcon),
	}
	player.initMod()

	return player
}

func NewTestPlayer() *Player {
	player := new(Player)
	//player.GetMod(ModPlay).(*ModPlayer) = new(ModPlayer)
	//player.ModIcon = new(ModIcon)
	//player.ModIcon.IconInfo = make(map[int]*Icon)
	player.ModCard = new(ModCard)
	player.ModCard.CardInfo = make(map[int]*Card)
	player.ModUniqueTask = new(ModUniqueTask)
	player.ModUniqueTask.MyTaskInfo = make(map[int]*TaskInfo)
	//player.ModUniqueTask.Locker = new(sync.RWMutex)
	player.ModRole = new(ModRole)
	player.ModRole.RoleInfo = make(map[int]*RoleInfo)
	player.ModBag = new(ModBag)
	player.ModBag.BagInfo = make(map[int]*ItemInfo)
	player.ModWeapon = new(ModWeapon)
	player.ModWeapon.WeaponInfo = make(map[int]*Weapon)

	player.ModRelic = new(ModRelic)
	player.ModRelic.RelicInfo = make(map[int]*Relic)

	player.ModCook = new(ModCook)
	player.ModCook.CookInfo = make(map[int]*Cook)

	player.ModHome = new(ModHome)
	player.ModHome.HomeItemInfo = make(map[int]*HomeItem)

	player.ModMap = new(ModMap)
	player.ModMap.InitData()
	//抽卡掉落模块
	player.ModWish = new(ModWish)
	player.ModWish.UPWishPool = new(WishPool)
	player.ModWish.NormalWishPool = new(WishPool)

	//****************************************
	//player.GetMod(ModPlay).(*ModPlayer).PlayerLevel = 1
	//player.GetMod(ModPlay).(*ModPlayer).Name = "旅行者"
	//player.GetMod(ModPlay).(*ModPlayer).WorldLevel = 1
	//player.GetMod(ModPlay).(*ModPlayer).WorldLevelNow = 1
	//player.GetMod(ModPlay).(*ModPlayer).player = player

	//****************************************
	player.modManage = make(map[string]ModBase)

	player.modManage = map[string]ModBase{
		ModPlay: new(ModPlayer),
		IconMod: new(ModIcon),
	}
	player.initMod()
	return player
}

// CreateRoleInDB
// @Description 将新的角色数据插入数据库
// @Author WangYuding 2022-04-05 18:13:42 ${time}
// @Return {
func (pr *Player) CreateRoleInDB() {
	//在数据库中生成对应的记录，根据记录生成对应的user_id
	DB.GormDB.Create(&pr.GetMod(ModPlay).(*ModPlayer).DBPlayer)
	//fmt.Println(player.ModPlayer.DBPlayer.ID)
	//告知客户端pID,同步已经生成的玩家ID给客户端
	DB.GormDB.Model(&pr.GetMod(ModPlay).(*ModPlayer).DBPlayer).Update("user_id", pr.GetMod(ModPlay).(*ModPlayer).DBPlayer.ID+100000000)
	pr.SyncPid()
	pr.Conn.SetProperty("PID", pr.GetMod(ModPlay).(*ModPlayer).UserId)

	//将玩家加入世界管理器中
	pr.LoadElse()
	WorldMgrObj.AddPlayer(pr)
	pr.SendStringMsg(2, pr.GetMod(ModPlay).(*ModPlayer).Name+",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
}

// initMod
// @Description 初始化mod模块
// @Author WangYuding 2022-04-05 17:39:26
func (pr *Player) initMod() {
	for _, v := range pr.modManage {
		v.init(pr)
	}
}

// GetModManager
// @Description 获得模块管理模块
// @Author WangYuding 2022-04-05 18:25:14
// @Return map[string]ModBase
func (pr *Player) GetModManager() map[string]ModBase {
	return pr.modManage
}

// GetUserID
// @Description 从当前Conn中获取uid
// @Author WangYuding 2022-04-05 18:28:32
func (pr *Player) GetUserID() (int, error) {
	uid, err := pr.Conn.GetProperty("PID")
	if err != nil {
		return -1, err
	}
	return uid.(int), nil
}

func (pr *Player) GetMod(Name string) ModBase {
	return pr.modManage[Name]
}

// SyncPid 告知客户端pID,同步已经生成的玩家ID给客户端
func (pr *Player) SyncPid() {
	log.Println("SyncPid")

	pidMsg := msgJson.SyncPID{PID: pr.GetMod(ModPlay).(*ModPlayer).UserId}
	data, err := json.Marshal(pidMsg) //
	if err != nil {
		log.Println(err)
		return
	}
	//调用zinx框架的SendMsg发包
	if err := pr.Conn.SendMsg(1, data); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
}

func (pr *Player) SendStringMsg(msgId uint32, msg string) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	if err := pr.Conn.SendMsg(msgId, data); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
}

// RecvSetIcon 对外接口
func (pr *Player) RecvSetIcon(iconId int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetIcon(iconId)
}

func (pr *Player) RecvSetCard(cardId int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetCard(cardId)
}

func (pr *Player) RecvSetName(name string) {
	pr.GetMod(ModPlay).(*ModPlayer).SetName(name)
}

func (pr *Player) RecvSetSign(sign string) {
	pr.GetMod(ModPlay).(*ModPlayer).SetSign(sign)
}

func (pr *Player) ReduceWorldLevel() {
	pr.GetMod(ModPlay).(*ModPlayer).ReduceWorldLevel()
}

func (pr *Player) ReturnWorldLevel() {
	pr.GetMod(ModPlay).(*ModPlayer).ReturnWorldLevel()
}

func (pr *Player) SetBirth(birth int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetBirth(birth)
}

func (pr *Player) SetShowCard(showCard []int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetShowCard(showCard)
}

func (pr *Player) SetShowTeam(showRole []int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetShowTeam(showRole)
}

func (pr *Player) SetHideShowTeam(isHide int) {
	pr.GetMod(ModPlay).(*ModPlayer).SetHideShowTeam(isHide)
}

func (pr *Player) Run() {
	//pr.ModWish.DoPoolTest()
	//ticker := time.NewTicker(1*time.Second)
	fmt.Println("Test Tools by YudingWang Learn from B站刘丹冰Aceld,大海葵,一棵平衡树")
	fmt.Println("模拟用户创建成功OK------开始测试")
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	//fmt.Println(pr.GetMod(ModPlay).(*ModPlayer).Name, ",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
	for {
		fmt.Println(pr.GetMod(ModPlay).(*ModPlayer).Name, ",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘） 5.地图")
		var modChoose int
		_, err := fmt.Scan(&modChoose)
		if err != nil {
			fmt.Println("Scan error!")
			return
		}
		switch modChoose {
		case 1:
			pr.HandleBase()
		case 2:
			pr.HandleBag()
		case 3:
			pr.HandleWishTest()
		case 4:
			pr.HandleWishUp()
		case 5:
			pr.HandleMap()
		}
		//fmt.Println(pr.GetMod(ModPlay).(*ModPlayer).Name, ",欢迎来到提瓦特大陆,请选择功能：1.基础信息 2.背包 3.up池抽卡模拟 4.up池抽卡（消耗相遇之缘）5.地图")
	}
}

// HandleBase 基础信息，测试模块
func (pr *Player) HandleBase() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作：0返回1查询信息2设置名字3设置签名4头像5名片6设置生日")
		var action int
		_, _ = fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			pr.HandleBaseGetInfo()
		case 2:
			pr.HandleBagSetName()
		case 3:
			pr.HandleBagSetSign()
		case 4:
			pr.HandleBagSetIcon()
		case 5:
			pr.HandleBagSetCard()
		case 6:
			pr.HandleBagSetBirth()
		}
	}
}

// HandleBaseGetInfoServer 用于服务器生成具体信息的列表
func (pr *Player) HandleBaseGetInfoServer() (res string) {
	res += fmt.Sprintln("名字:", pr.GetMod(ModPlay).(*ModPlayer).Name)
	res += fmt.Sprintln("等级:", pr.GetMod(ModPlay).(*ModPlayer).PlayerLevel)
	res += fmt.Sprintln("大世界等级:", pr.GetMod(ModPlay).(*ModPlayer).WorldLevelNow)
	if pr.GetMod(ModPlay).(*ModPlayer).Sign == "" {
		res += fmt.Sprintln("签名:", "未设置")
	} else {
		res += fmt.Sprintln("签名:", pr.GetMod(ModPlay).(*ModPlayer).Sign)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Icon == 0 {
		res += fmt.Sprintln("头像:", "未设置")
	} else {
		res += fmt.Sprintln("头像:", csvs.GetItemConfig(pr.GetMod(ModPlay).(*ModPlayer).Icon), pr.GetMod(ModPlay).(*ModPlayer).Icon)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Card == 0 {
		res += fmt.Sprintln("名片:", "未设置")
	} else {
		res += fmt.Sprintln("名片:", csvs.GetItemConfig(pr.GetMod(ModPlay).(*ModPlayer).Card), pr.GetMod(ModPlay).(*ModPlayer).Card)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Birth == 0 {
		res += fmt.Sprintln("生日:", "未设置")
	} else {
		res += fmt.Sprintln("生日:", pr.GetMod(ModPlay).(*ModPlayer).Birth/100, "月", pr.GetMod(ModPlay).(*ModPlayer).Birth%100, "日")
	}
	return res
}

func (pr *Player) HandleBaseGetInfo() {
	fmt.Println("名字:", pr.GetMod(ModPlay).(*ModPlayer).Name)
	fmt.Println("等级:", pr.GetMod(ModPlay).(*ModPlayer).PlayerLevel)
	fmt.Println("大世界等级:", pr.GetMod(ModPlay).(*ModPlayer).WorldLevelNow)
	if pr.GetMod(ModPlay).(*ModPlayer).Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", pr.GetMod(ModPlay).(*ModPlayer).Sign)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", csvs.GetItemConfig(pr.GetMod(ModPlay).(*ModPlayer).Icon), pr.GetMod(ModPlay).(*ModPlayer).Icon)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", csvs.GetItemConfig(pr.GetMod(ModPlay).(*ModPlayer).Card), pr.GetMod(ModPlay).(*ModPlayer).Card)
	}

	if pr.GetMod(ModPlay).(*ModPlayer).Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", pr.GetMod(ModPlay).(*ModPlayer).Birth/100, "月", pr.GetMod(ModPlay).(*ModPlayer).Birth%100, "日")
	}
}

func (pr *Player) HandleBagSetName() {
	fmt.Println("请输入名字:")
	var name string
	fmt.Scan(&name)
	pr.RecvSetName(name)
}

func (pr *Player) HandleBagSetSign() {
	fmt.Println("请输入签名:")
	var sign string
	fmt.Scan(&sign)
	pr.RecvSetSign(sign)
}

func (pr *Player) HandleBagSetIcon() {
	for {
		fmt.Println("当前处于基础信息--头像界面,请选择操作：0返回1查询头像背包2设置头像")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			pr.HandleBagSetIconGetInfo()
		case 2:
			pr.HandleBagSetIconSet()
		}
	}
}

func (pr *Player) HandleBagSetIconGetInfo() {
	fmt.Println("当前拥有头像如下:")
	for _, v := range pr.GetMod(IconMod).(*ModIcon).IconInfo {
		config := csvs.GetItemConfig(v.IconId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (pr *Player) HandleBagSetIconSet() {
	fmt.Println("请输入头像id:")
	var icon int
	fmt.Scan(&icon)
	pr.RecvSetIcon(icon)
}

func (pr *Player) HandleBagSetCard() {
	for {
		fmt.Println("当前处于基础信息--名片界面,请选择操作：0返回1查询名片背包2设置名片")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			pr.HandleBagSetCardGetInfo()
		case 2:
			pr.HandleBagSetCardSet()
		}
	}
}

func (pr *Player) HandleBagSetCardGetInfo() {
	fmt.Println("当前拥有名片如下:")
	for _, v := range pr.ModCard.CardInfo {
		config := csvs.GetItemConfig(v.CardId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (pr *Player) HandleBagSetCardSet() {
	fmt.Println("请输入名片id:")
	var card int
	fmt.Scan(&card)
	pr.RecvSetCard(card)
}

func (pr *Player) HandleBagSetBirth() {
	if pr.GetMod(ModPlay).(*ModPlayer).Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	pr.GetMod(ModPlay).(*ModPlayer).SetBirth(month*100 + day)
}

//
func (pr *Player) HandleBag() {
	for {
		fmt.Println("当前处于背包界面,请选择操作：0返回1增加物品2扣除物品3使用物品")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			pr.HandleBagAddItem()
		case 2:
			pr.HandleBagRemoveItem()
		case 3:
			pr.HandleBagUseItem()
		}
	}
}

func (pr *Player) HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	pr.ModBag.UseItem(itemId, int64(itemNum), pr)
}

func (pr *Player) HandleBagAddItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	pr.ModBag.AddItem(itemId, int64(itemNum), pr)
}

func (pr *Player) HandleBagRemoveItem() {
	BagId := 0
	fmt.Println("当前处于删除物品界面界面,请选择操作：0返回1普通背包2圣遗物背包3武器背包")
	fmt.Scan(&BagId)
	switch BagId {
	case 0:
		return
	case csvs.NormalBagId:
		itemId := 0
		itemNum := 0
		fmt.Println("物品ID")
		fmt.Scan(&itemId)
		fmt.Println("物品数量")
		fmt.Scan(&itemNum)
		if err := pr.ModBag.RemoveItem(itemId, int64(itemNum), pr); err != nil {
			fmt.Println(err)
		}
	case csvs.RelicBagId:
		keyId := 0
		fmt.Println("圣遗物编号")
		fmt.Scan(&keyId)
		pr.ModRelic.RemoveItem(keyId)
	case csvs.WeaponBagId:
		keyId := 0
		fmt.Println("武器编号")
		fmt.Scan(&keyId)
		pr.ModWeapon.RemoveItem(keyId)
	}
}

// HandleMap 地图模块
func (pr *Player) HandleMap() {
	fmt.Println("向着星辰与深渊,欢迎来到冒险家协会！")
	//fmt.Println("当前位置:", "蒙德城")
	for {
		fmt.Println("选择交互地图 0，返回 1.蒙德 2.璃月 1001.深入风龙废墟 2001.无妄引咎迷宫")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			pr.HandleMapIn(action)

		}
	}
}

// HandleMapIn 进入地图模块
func (pr *Player) HandleMapIn(mapId int) {
	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		fmt.Println("地图无法识别")
		return
	}
	//重新进入地图时候的刷新工作(秘境地图)
	pr.ModMap.RefreshWhenCome(mapId)
Loop:
	for {
		//检查当前进入地图的时候有没有遗留物
		pr.ModMap.checkAnyDropOnMap(mapId, pr)
		//生成当前可选事件列表
		pr.ModMap.GetEventList(mapId)
		fmt.Println("请选择触发事件Id(0返回)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			eventConfig := csvs.ConfigMapEventMap[action]
			if eventConfig == nil {
				fmt.Println("无法识别的事件")
				break
			}
			if err := pr.ModMap.SetEventState(mapId, eventConfig.EventId, csvs.EventEnd, pr); err != nil {
				fmt.Println(err)
				break Loop
			}
		}
	}
}

func (pr *Player) HandleWishTest() {
	times := 0
	fmt.Println("请输出抽卡次数")
	fmt.Scan(&times)
	pr.ModWish.DoPoolTest(times)
}

func (pr *Player) HandleWishUp() {
	for {
		var choice int
		fmt.Println("您现在在抽卡界面 按0返回 按1祈愿1次 按2祈愿10次 按3查询抽卡信息")
		fmt.Scan(&choice)
		switch choice {
		case 0:
			return
		case 1:
			fmt.Println("如果祈愿之缘数量不足，请通过背包功能增加祈愿之缘，物品id为1000005")
			if err := pr.ModBag.RemoveItem(1000005, 1, pr); err != nil {
				fmt.Println(err)
				continue
			}
			pr.ModWish.DoPool(1, pr)
		case 2:
			fmt.Println("如果祈愿之缘数量不足，请通过背包功能增加祈愿之缘，物品id为1000005")
			if err := pr.ModBag.RemoveItem(1000005, 10, pr); err != nil {
				fmt.Println(err)
				continue
			}
			pr.ModWish.DoPool(10, pr)
		case 3:
			fmt.Printf("本次您一共进行了%d次祈愿，共获得五星角色%d位，占总数的%.4f%%,四星角色%d位，四星武器%d把，四星物品占总数的%.4f%%\n当前您的五星保底为%d抽，四星保底为%d抽\n",
				pr.ModWish.UPWishPool.StatTotalWishes, pr.ModWish.UPWishPool.StatFiveTotal,
				100*float32(pr.ModWish.UPWishPool.StatFiveTotal)/float32(pr.ModWish.UPWishPool.StatTotalWishes),
				pr.ModWish.UPWishPool.StatFourRole, pr.ModWish.UPWishPool.StatFourWeapon,
				100*float32(pr.ModWish.UPWishPool.StatFourRole+pr.ModWish.UPWishPool.StatFourWeapon)/float32(pr.ModWish.UPWishPool.StatTotalWishes),
				pr.ModWish.UPWishPool.FiveStarTimes, pr.ModWish.UPWishPool.FourStarTimes)
		default:
			fmt.Println("无效输入")

		}
	}
}

// LoadElse
// @Description load mod after the modplayer have load
// @Author WangYuding 2022-04-09 22:32:50
func (pr *Player) LoadElse() {
	for x, v := range pr.modManage {
		if x != ModPlay {
			v.LoadData()
		}
	}
}
