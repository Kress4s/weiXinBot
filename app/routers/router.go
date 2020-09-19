package routers

import (
	"weiXinBot/app/index"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/index/login/qr_code", &index.IndexController{}, "get:GetQrCode")
	beego.Router("/index/login/check", &index.IndexController{}, "get:Check")
}
