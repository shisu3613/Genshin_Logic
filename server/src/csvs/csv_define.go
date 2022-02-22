package csvs

const (
	LogicFalse = 0
	LogicTrue  = 1
)

//基础面板模块常数
const (
	ReduceWorldLevelStart    = 5  //降低世界等级的要求
	ReduceWorldLevelMax      = 1  //最多能降低多少级
	ReduceWorldLevelCoolTime = 10 //冷却时间
	ShowSize                 = 9
	AddRoleTimeNormalMin     = 2
	AddRoleTimeNormalMax     = 7
	MaxWeaponSize            = 2000
	MaxRelicSize             = 1500
	NormalBagId              = 1
	WeaponBagId              = 3
	RelicBagId               = 2
)

//抽卡常数
const (
	FiveStarLimit          = 73
	FiveStarLimitIncrement = 600
	FourStarLimit          = 8
	FourStarLimitIncrement = 5500
	TotalWishWeight        = 10000
)

//地图模块常数
const (
	// EventStart 事件开始Flag
	EventStart = 0
	// EventFinish 当前事件回事件结束Flag
	EventFinish = 9
	// EventEnd 事件结束Flag
	EventEnd = 10

	// MapRefreshCant 永不刷新
	MapRefreshCant = 0

	// MapRefreshTwoDay 48小时刷新的植物
	MapRefreshTwoDay = 1

	// MapRefreshTwoThreeDay 三天刷新一次的矿物
	MapRefreshTwoThreeDay = 2

	// MapRefreshWeek  一周刷新一次的周本
	MapRefreshWeek = 3

	// MapRefreshSelf 随着玩家进入刷新
	MapRefreshSelf = 4

	// MapRefreshHalfDay 每隔离12小时刷新的怪物
	MapRefreshHalfDay = 5

	RefreshSystem = 1
	RefreshPlayer = 2
)

//物品掉落常数
const (
	DropWeightAll     = 10000
	DropOneItem       = 1
	DropGroupItems    = 2
	DropWeightedItems = 3
)
