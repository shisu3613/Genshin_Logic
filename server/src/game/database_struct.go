package game

import "github.com/jinzhu/gorm"

type DBPlayer struct {
	gorm.Model
	UserId         int    `gorm:"unique_index:uid;"` //唯一id `gorm:"unique_index:hash_idx;"`
	Icon           int    //头像   新增icon模块
	Card           int    //名片   新增card模块
	Name           string //名字   新增banword模块
	Sign           string //签名
	PlayerLevel    int    //等级
	PlayerExp      int    //阅历(经验)
	WorldLevel     int    //大世界等级
	WorldLevelNow  int    //大世界等级(当前)
	WorldLevelCool int64  //操作大世界等级的冷却时间
	Birth          int    //生日

	//ShowCard     []int       //展示名片
	//ShowTeam     []*ShowRole //展示阵容
	HideShowTeam int //隐藏开关,是否包含角色属性

	//不可见字段
	Prohibit int //封禁状态
	IsGM     int //GM账号标志
}

func (DBPlayer) TableName() string {
	return "BasicProfiles"
}
