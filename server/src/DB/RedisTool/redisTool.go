package RedisTool

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"time"
)

/**
    @author: WangYuding
    @since: 2022/4/27
    @desc: //Redis的工具库
**/

type RedisConfig struct {
	Addr     string
	Password string
}

var ctx = context.Background()

func NewRedis(db int) *redis.Client { //将数据库连接操作打包为方法使用newRdis(0)方法带入数据库名调用即可
	var configure *RedisConfig
	data, err := ioutil.ReadFile("json/RedisConfig.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &configure)
	rdb := redis.NewClient(&redis.Options{
		Addr:     configure.Addr,     //目前数据库默认安装在本机，监听localhost，默认端口为6379
		Password: configure.Password, // simply password set
		DB:       db,                 // use default DB
	})
	return rdb //返回数据库客户端
}

// CheckKeyExists
// @Description: check Key exist or not
// @param db
// @param key
// @return bool
func CheckKeyExists(db int, key string) bool {
	rdb := NewRedis(db)
	n, err := rdb.Exists(ctx, key).Result()
	CheckError(err)
	if n == 0 {
		return false
	}
	return true

}

func GetLastListVal(db int, key string) string {
	client := NewRedis(db)
	val, err := client.LIndex(ctx, key, -1).Result()
	CheckError(err)
	return val
}

// GetAllKeys 获取该数据库里所有的key
func GetAllKeys(db int) []string {
	rdb := NewRedis(db)
	defer rdb.Close()
	keys, err := rdb.Keys(ctx, "*").Result()
	CheckError(err)
	return keys
}

func AddSetValByKey(db int, key string, val string) {
	rdb := NewRedis(db)
	err := rdb.SAdd(ctx, key, val).Err()
	if err != nil {
		panic(err)
	}
}

func GetSetByKey(db int, key string) []string {
	rdb := NewRedis(db)
	es, _ := rdb.SMembers(ctx, key).Result()
	return es
}

func GetValueByKey(db int, key string) (string, error) {
	rdb := NewRedis(db)
	defer rdb.Close()
	val, err := rdb.Get(ctx, key).Result() //使用IdTime+UserId获取Message
	CheckError(err)
	return val, err
}

func GetListByKey(db int, key string) ([]string, error) {
	rdb := NewRedis(db)
	defer rdb.Close()
	val, err := rdb.LRange(ctx, key, 0, -1).Result() //获取list的长度
	CheckError(err)
	return val, err
}

func SetRecord(db int, key string, data []byte) bool {
	rdb := NewRedis(db)
	err := rdb.Set(ctx, key, data, 24*time.Hour).Err()
	if err != nil {
		log.Println(err)
		return false
	} //保存Message,保存24小时24*time.Hour
	err = rdb.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func SetRecordList(db int, key string, data string) bool {
	rdb := NewRedis(db)
	err := rdb.RPush(ctx, key, data).Err()
	if err != nil {
		log.Println(err)
		return false
	}
	err = rdb.Close()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// CheckError error处理,可以使用客户端log处理的逻辑将错误信息收集保存到数据库,这里不在展开
func CheckError(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
