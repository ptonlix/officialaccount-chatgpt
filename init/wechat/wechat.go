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
}

var WxConf WechatConfig

func init() {
	// 获取微信配置
	var err error
	WxConf.Appid, err = beego.AppConfig.String("wxappid")
	if err != nil {
		logs.Error("微信配置加载失败: Appid")
	}
	WxConf.AppSecret, err = beego.AppConfig.String("wxappsecret")
	if err != nil {
		logs.Error("微信配置加载失败: AppSecret")
	}
	WxConf.Token, err = beego.AppConfig.String("wxtoken")
	if err != nil {
		logs.Error("微信配置加载失败: Token")
	}
	WxConf.EncodingAESKey, err = beego.AppConfig.String("wxencodingasekey")
	if err != nil {
		logs.Info("微信配置加载失败: EncodingAESKey")
	}
}
