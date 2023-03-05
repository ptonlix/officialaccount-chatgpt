package chatgpt

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type ChatgptConfig struct {
	Key         string
	Model       string
	MaxTokens   int
	Temperature float64
	Url         string
}

var ChatGptConf ChatgptConfig

func init() {
	// 获取ChatGpt配置
	var err error
	ChatGptConf.Key, err = beego.AppConfig.String("gptkey")
	if err != nil {
		logs.Error("ChatGpt配置加载失败: Key")
	}
	ChatGptConf.Model, err = beego.AppConfig.String("gptmodel")
	if err != nil {
		logs.Error("ChatGpt配置加载失败: Model")
	}
	ChatGptConf.MaxTokens, err = beego.AppConfig.Int("gptmaxtoken")
	if err != nil {
		logs.Error("ChatGpt配置加载失败: MaxTokens")
	}
	ChatGptConf.Temperature, err = beego.AppConfig.Float("gpttemperature")
	if err != nil {
		logs.Error("ChatGpt配置加载失败: Temperature")
	}
	ChatGptConf.Url, err = beego.AppConfig.String("gpturl")
	if err != nil {
		logs.Error("ChatGpt配置加载失败: Url")
	}

}
