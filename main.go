package main

import (
	_ "github.com/ptonlix/officialaccount-chatgpt/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 初始化微信公众号实例
	
	beego.Run()
}
