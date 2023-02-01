# Officialaccount-ChatGPT
简体中文
<p>
	<p align="center">
		<img height=280 src="./doc/image/gopher.png">
	</p>
	<p align="center">
		<font size=6 face="宋体">微信公众号体验ChatGPT</font>
	<p>
</p>
<p align="center">
<img alt="Go" src="https://img.shields.io/badge/Go-1.7%2B-blue">
<img alt="Mysql" src="https://img.shields.io/badge/OpenAI-API-brightgreeni">
<img alt="Redis" src="https://img.shields.io/badge/Redis-6.2%2B-yellowgreen">
<img alt="influxDB" src="https://img.shields.io/badge/
WechatSDK-2.0%2B-orange">
<img alt="license" src="https://img.shields.io/badge/license-Apache-lightgrey">
</p>

> ChatGPT是由人工智能研究实验室[OpenAI](https://baike.baidu.com/item/OpenAI/19758408?fromModule=lemma_inlink)在2022年11月30日发布的全新聊天机器人模型，一款[人工智能](https://baike.baidu.com/item/%E4%BA%BA%E5%B7%A5%E6%99%BA%E8%83%BD/9180?fromModule=lemma_inlink)技术驱动的自然语言处理工具。它能够通过学习和理解人类的语言来进行对话，还能根据聊天的上下文进行互动，真正像人类一样来聊天交流，甚至能完成撰写邮件、视频脚本、文案、翻译、代码等任务。同时也引起无数网友沉迷与ChatGPT聊天，成为大家讨论的火爆话题。

# 01.前言

大家好，我是baird～

去年12月份大火的ChatGPT，让人印象深刻，尤其是不在像之前的AI客服那样答非所问。只要你提问合理，它回答的问题基本上就是你想要的答案。比如，通过ChatGPT可以让你写代码更轻松，加以引导，它就能为你写出一段符合要求的函数代码，你修改一下就能达到质量要求。

目前网上体验ChatGPT的方式主要有两种
1. 直接官网体验 [ChatGPT官网](https://openai.com/blog/chatgpt/)  
   这种体验效果较好，但需要梯子访问
2. API访问  
   通过比如微信机器人或微信公众号等作为前端入口，后台调用openai的API接口返回数据。 这种体验稍差，但在国内每个人都体验无需梯子。

独乐乐不如众乐乐，为了让更多人能快速体验到ChatGPT，我选择第二种方式，做一个简单的开发，共享出来让大家访问。

对于微信机器人的方式，网上已有很多教程，这里我直接选择微信公众号的方式进行开发。**（有公众号开发需求的朋友可以重点关注）**

> 直接官网体验参考地址 [体验教程](https://zhuanlan.zhihu.com/p/589201080)  
> 微信机器人开发参考 [参考地址](https://blog.csdn.net/zhayujie5200/article/details/128359055)  

# 02.准备阶段

## 准备OpenAI账号

1. 自己注册  直接上链接供大家参考 [注册参考地址](https://51.ruyo.net/18252.html)
2. 万能淘宝  淘宝直接搜ChatGPT账号，9.9带回家

自己注册会遇到手机验证码验证问题，淘宝也有此服务，但为了省时间，我直接选择了购买账号。

## 获取API KEY

登陆账号则前往 [API管理页面](https://beta.openai.com/account/api-keys) 创建一个 API Key 并保存下来，后面需要在项目中配置这个key。

![WX20230131-090207@2x.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/640f16305b634e2cafbb3dfc41034e95~tplv-k3u1fbpfcp-watermark.image?)

账号创建有免费的 $18 额度，使用完可以更换邮箱重新注册。目前体验期是到4月1日截止，没来得及体验的朋友尽快体验噢～🚗

![WX20230131-090155@2x.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a2f40b26056e4f91a3686f0721144a7e~tplv-k3u1fbpfcp-watermark.image?)

## 微信公众号

微信公众平台上注册或使用测试公众号

# 03.开发阶段

## 流程框架

![frame.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/806510eb22fc4bb88d123f35b1990c52~tplv-k3u1fbpfcp-watermark.image?)

1. 根据微信公众号开发规范，我们需要开发一个后台服务，这个服务需要提供一个统一API接口，微信公众号会去调用这个API接口，将用户信息转发到该服务。
2. 后台服务将用户消息再包装，调用OpenAI的API接口，该接口则返回ChatGPT回复的信息，该服务再返回给用户。

整体流程比较简单，主要关注的点是微信公众号开发流程，需要验证AppSecret等信息，详情可参考[微信开发文档](https://developers.weixin.qq.com/doc/offiaccount/Getting_Started/Overview.html)

## 软件框架
这里我采用 Go语言进行开发，使用Beego+WechatSDK快速开发
>Beego: https://github.com/beego/beego  
>Go Wechat SDK: https://github.com/silenceper/wechat

后台服务流程图如下：

![software.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1d30b314de9a406da65da9f9d974ab9c~tplv-k3u1fbpfcp-watermark.image?)

## 开发细节

- 主要处理逻辑如下，这里有个问题需要注意。微信公众号回复的超时时间是5s，而访问OpenAI的API处理时间往往是超过5秒的，这里需要用到异步处理。获取用户的openid，调用微信API再把信息返回给该用户

```go
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

```
- 配置文件说明,这里主要说明一下OpenAI的API接口配置  
gptmodel 选择哪个模型，默认可以参考text-davinci-003这个模型即可  
gptmaxtoken 最大输入的字符数，默认最大为2048。输入越长，收费越高，具体参考OpenAI的API文档  
gpttemperature 数据范围在0～1 数据越高代表回答的随机性越大，感兴趣的朋友可以自行调整体验
```ini
#chatgpt
gptkey = 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
gptmodel =  "text-davinci-003"
gptmaxtoken = 1024
gpttemperature = 0.7
```

- 微信公众号调试时，使用内网穿透软件能方便在本地进行调试   
这里我推荐一款我在用的免费软件[netapp](https://natapp.cn/),大家可以自行查阅

# 04.部署运行
项目部署采用容器部署,命令如下
```shell
# 生成镜像
docker build -t officialaccount-chatgpt:v1 -f Dockerfile .
# 运行容器
docker run -itd -p 8080:8080 --restart always -v /root/officialaccount/conf:/app/conf ptonlix/officialaccount-chatgpt:v1
# 查看日志
docker logs -f 6b05e1c81380
```

# 05.体验环境
附上我自己的公众号，关注后直接输入即可体验ChatGPT，欢迎大家体验
<p>
	<p align="center">
		<img height=280 src="./doc/image/wx.jpg">
	</p>
</p>

# 06.后续

目前该项目只是一个Demo，还有很多可以优化的地方。比如，通过识别特定前缀判断是访问ChatGPT的请求，避免所有请求都访问，这块需要结合业务去考虑。  
最后希望能给大家一点启发和帮助，感谢阅读。

<p align="center">
  <b>SPONSORED BY</b>
</p>
<p align="center">
   <a href="https://www.gogeek.com.cn/" title="gogeek" target="_blank">
      <img height="200px" src="https://img.gejiba.com/images/96b6d150bd758b13d66aec66cb18044e.jpg" title="gogeek">
   </a>
</p>