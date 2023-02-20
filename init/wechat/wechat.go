package wechat

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type WechatConfig struct {
	Appid          string // 微信公众号配置参数
	AppSecret      string
	Token          string
	EncodingAESKey string

	RedisAddr   string
	RedisPasswd string
	RedisDb     int // redis db

}

var WxConf WechatConfig

func init() {
	// 获取微信配置
	var err error
	WxConf.Appid, err = beego.AppConfig.String("wxappid")
	if err != nil {
		logs.Error("微信配置加载失败: wxappid, 请检查app.conf配置文件")
	}
	WxConf.AppSecret, err = beego.AppConfig.String("wxappsecret")
	if err != nil {
		logs.Error("微信配置加载失败: wxappsecret, 请检查app.conf配置文件")
	}
	WxConf.Token, err = beego.AppConfig.String("wxtoken")
	if err != nil {
		logs.Error("微信配置加载失败: wxtoken, 请检查app.conf配置文件")
	}
	WxConf.EncodingAESKey, err = beego.AppConfig.String("wxencodingasekey")
	if err != nil {
		logs.Info("微信配置加载失败: wxencodingasekey, 请检查app.conf配置文件")
	}

	WxConf.RedisAddr, err = beego.AppConfig.String("redisaddr")
	if err != nil {
		logs.Info("微信配置加载失败: redisaddr, 请检查app.conf配置文件")
	}
	WxConf.RedisPasswd, err = beego.AppConfig.String("redispasswd")
	if err != nil {
		logs.Info("微信配置加载失败: redispasswd, 请检查app.conf配置文件")
	}
	WxConf.RedisDb, err = beego.AppConfig.Int("redisdb")
	if err != nil {
		logs.Info("微信配置加载失败: redisdb, 请检查app.conf配置文件")
	}

}
