package services

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	wechatConfig "github.com/ptonlix/officialaccount-chatgpt/init/wechat"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type OfficialAccountService struct {
	OA     *officialaccount.OfficialAccount
	ReqMsg string //发送消息
	OpenId string //发送者的openid
}

func NewOfficialAccountService(w *wechatConfig.WechatConfig) *OfficialAccountService {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:          w.Appid,
		AppSecret:      w.AppSecret,
		Token:          w.Token,
		EncodingAESKey: w.EncodingAESKey,
		Cache:          memory,
	}
	return &OfficialAccountService{OA: wc.GetOfficialAccount(cfg)}
}
func (oa *OfficialAccountService) Handle(msg *message.MixMessage) *message.Reply {
	//记录收到的信息
	oa.ReqMsg = msg.Content
	//回复消息：演示回复用户发送的消息
	text := message.NewText("以下内容由ChatGPT提供:")
	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
}

func (oa *OfficialAccountService) SendSyncMessage() {
	gptService := new(ChatGptService)
	replyText, err := DoRetryHTTP(func() (interface{}, error) {
		return gptService.Completions(oa.ReqMsg)
	}, 3, time.Second*2)
	if err != nil {
		logs.Error("Retry GPT HTTP API Failed!")
		replyText = "GPT服务器异常，请稍后重试。。。"
	}
	offAccService := NewOfficialAccountService(&wechatConfig.WxConf)
	cm := offAccService.OA.GetCustomerMessageManager()
	data := message.NewCustomerTextMessage(oa.OpenId, replyText.(string))
	err = cm.Send(data)
	if err != nil {
		logs.Debug(err.Error())
	}
}
