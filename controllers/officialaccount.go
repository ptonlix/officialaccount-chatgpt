package controllers

import (
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	wechatConfig "github.com/ptonlix/officialaccount-chatgpt/init/wechat"
	"github.com/ptonlix/officialaccount-chatgpt/services"
)

type OfficialaccountController struct {
	beego.Controller
}

func (o *OfficialaccountController) AnyHandle() {

	// 传入request和responseWriter
	offAccService := services.NewOfficialAccountService(&wechatConfig.WxConf)
	server := offAccService.OA.GetServer(o.Ctx.Request, o.Ctx.Output.Context.ResponseWriter)
	//设置接收消息的处理方法
	server.SetMessageHandler(offAccService.Handle)
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	offAccService.OpenId = server.GetOpenID() //获取openid
	//发送回复的消息
	server.Send()
	go offAccService.SendSyncMessage()
}
