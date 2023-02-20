package services

import (
	"context"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/go-redis/redis/v8"
	"github.com/ptonlix/officialaccount-chatgpt/init/chatgpt"
	wechatRedis "github.com/ptonlix/officialaccount-chatgpt/init/redis"
	wechatConfig "github.com/ptonlix/officialaccount-chatgpt/init/wechat"
	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

type OfficialAccountService struct {
	OA          *officialaccount.OfficialAccount
	ReqMsg      string //发送消息
	ChatMsg     string //聊天记录上下文
	OpenId      string //发送者的openid
	RedisClient *redis.Client
}

func NewOfficialAccountService(w *wechatConfig.WechatConfig) *OfficialAccountService {
	wc := wechat.NewWechat()
	redisOpts := &cache.RedisOpts{
		Host:        w.RedisAddr,
		Password:    w.RedisPasswd,
		Database:    w.RedisDb,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: 60, //second
	}
	//这里本地内存保存access_token，也可选择redis，memcache或者自定cache
	redisCache := cache.NewRedis(context.Background(), redisOpts)
	cfg := &offConfig.Config{
		AppID:          w.Appid,
		AppSecret:      w.AppSecret,
		Token:          w.Token,
		EncodingAESKey: w.EncodingAESKey,
		Cache:          redisCache,
	}
	return &OfficialAccountService{OA: wc.GetOfficialAccount(cfg), RedisClient: wechatRedis.Client}
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
	//获取当前用户上下文信息,检查聊天记录上下文长度
	oa.CheckChatContext()
	//开始请求API
	replyText, err := DoRetryHTTP(func() (interface{}, error) {
		return gptService.Completions(oa.ChatMsg)
	}, 3, time.Second*2)
	if err != nil {
		logs.Error("Retry GPT HTTP API Failed!")
		replyText = "GPT服务器异常,请稍后重试。。。"
	} else {
		// 记录上下文信息
		if ret := strings.LastIndex(replyText.(string), oa.ReqMsg); ret != -1 {
			arrStr := strings.Split(replyText.(string)[ret:], "\n\n")
			if len(arrStr) > 1 {
				replyText = arrStr[1]
			} else {
				replyText = arrStr[0]
			}
		}
		chatRedisMsg := oa.ChatMsg + "\n\n" + replyText.(string)
		if err := oa.SetChatContext(chatRedisMsg); err != nil {
			replyText = "GPT服务器异常,请稍后重试。。。"
		}
	}
	offAccService := NewOfficialAccountService(&wechatConfig.WxConf)
	cm := offAccService.OA.GetCustomerMessageManager()
	data := message.NewCustomerTextMessage(oa.OpenId, replyText.(string))
	err = cm.Send(data)
	if err != nil {
		logs.Debug("发送消息失败%v", err)
	}

}

// 获取用户聊天记录上下文
func (oa *OfficialAccountService) GetChatContext() string {

	val, err := oa.RedisClient.Get(context.Background(), oa.OpenId).Result()
	if err != nil {
		logs.Warn("查询Redis获取上下文错误:%v", err)
		return ""
	}
	return val
}

//设置用户聊天记录上下文
func (oa *OfficialAccountService) SetChatContext(chatmsg string) error {

	_, err := oa.RedisClient.Set(context.Background(), oa.OpenId, chatmsg, time.Minute*10).Result()
	if err != nil {
		logs.Warn("Redis保存用户上下文错误:%v", err)
		return err
	}
	return nil
}

//检查聊天记录上下文长度
func (oa *OfficialAccountService) CheckChatContext() {
	msg := oa.GetChatContext()
	if msg == "" {
		oa.ChatMsg = oa.ReqMsg
	} else {
		oa.ChatMsg = msg + "\n\n" + oa.ReqMsg
	}
	if len(oa.ChatMsg) < chatgpt.ChatGptConf.MaxTokens/2 {
		return
	}
	oa.ChatMsg = oa.ReqMsg
	//清除Redis数据
	_, err := oa.RedisClient.Del(context.Background(), oa.OpenId).Result()
	if err != nil {
		logs.Warn("删除Redis用户上下文失败:%v", err)
		oa.RedisClient.Del(context.Background(), oa.OpenId)
	}
}
