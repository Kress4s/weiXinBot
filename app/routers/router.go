package routers

import (
	"weiXinBot/app/index"
	"weiXinBot/app/main/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/index/login/qr_code", &index.IndexController{}, "get:GetQrCode")
	beego.Router("/index/login/check", &index.IndexController{}, "get:Check")
	beego.Router("/manager/login/?:authtype", &index.MgrIndexController{}, "post:Login")
	beego.Router("/manager/register", &index.MgrIndexController{}, "post:Register")

	ns1 := beego.NewNamespace("v1",
		beego.NSNamespace("bot",
			beego.NSInclude(&controllers.BotsController{}),
		),
		beego.NSNamespace("group",
			beego.NSInclude(&controllers.GroupController{}),
		),
	)
	beego.AddNamespace(ns1)
}
