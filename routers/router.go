package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/ptonlix/officialaccount-chatgpt/controllers"
)

func init() {
	beego.Router("/officialaccount", &controllers.OfficialaccountController{}, "*:AnyHandle")

}
