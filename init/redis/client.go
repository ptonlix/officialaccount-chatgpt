package redis

import (
	"context"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func init() {

	RedisAddr, err := beego.AppConfig.String("redisaddr")
	if err != nil {
		logs.Error("Redis配置加载失败: redisaddr, 请检查app.conf配置文件")
	}
	RedisPasswd, err := beego.AppConfig.String("redispasswd")
	if err != nil {
		logs.Error("Redis配置加载失败: redispasswd, 请检查app.conf配置文件")
	}
	RedisDb, err := beego.AppConfig.Int("redisdb")
	if err != nil {
		logs.Error("Redis配置加载失败: redisdb, 请检查app.conf配置文件")
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     RedisAddr,
		Password: RedisPasswd,
		DB:       RedisDb,
	})
	if _, err := Client.Ping(context.Background()).Result(); err != nil {
		logs.Error("Redis连接失败,请检查")
	}
}
