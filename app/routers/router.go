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
	beego.Router("/manager/getmyinfo", &index.MgrIndexController{}, "get:GetMyInfo")
	beego.Router("/qiniu/getoken", &index.QiniuController{}, "get:GetUploadToken")

	ns1 := beego.NewNamespace("/v1",
		beego.NSNamespace("/bot",
			beego.NSInclude(&controllers.BotsController{}),
		),
		beego.NSNamespace("/group",
			beego.NSInclude(&controllers.GroupController{}),
		),
		beego.NSNamespace("/grouplan",
			beego.NSInclude(&controllers.GrouPlanController{}),
		),
		beego.NSNamespace("/welcome",
			beego.NSInclude(&controllers.WelcomeController{}),
		),
		beego.NSNamespace("/resource",
			beego.NSInclude(&controllers.ResouceController{}),
		),
		beego.NSNamespace("/material",
			beego.NSInclude(&controllers.MaterialController{}),
		),
		beego.NSNamespace("/keyword",
			beego.NSInclude(&controllers.KeyWordsController{}),
		),
		beego.NSNamespace("/question",
			beego.NSInclude(&controllers.QuestionController{}),
		),
		beego.NSNamespace("/config",
			beego.NSInclude(&controllers.ConfigurationController{}),
		),
		beego.NSNamespace("/timetask",
			beego.NSInclude(&controllers.TimeTaskController{}),
		),
	)
	beego.AddNamespace(ns1)
}
