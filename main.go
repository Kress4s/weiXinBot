package main

import (
	"fmt"
	"os"
	"weiXinBot/app/bridage/common"
	_ "weiXinBot/app/bridage/common/dbmysql"

	_ "weiXinBot/app/bridage/flows/timetask"
	_ "weiXinBot/app/bridage/grpc"
	_ "weiXinBot/app/bridage/models"
	_ "weiXinBot/app/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/astaxie/beego/toolbox"
)

func main() {
	// var checkAcess = func(ctx *context.Context) {
	// 	if token := ctx.Input.Header("Authorization"); token != "" {
	// 		fmt.Println(ctx.Input.URL())
	// 		// ctx.ResponseWriter.WriteHeader(constant.EXPEIRE_ACCOUNT_CODE)
	// 		fmt.Println(token)
	// 	}
	// }
	// 接口认证(调试阶段屏蔽)
	// beego.InsertFilter("^(?!.*/manager)", beego.BeforeRouter, checkAcess)

	//允许跨站访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Cookie", "WX_TOKEN"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	//}
	createDIR() //初始化必要目录
	//设置日志规则
	logs.SetLogger(logs.AdapterMultiFile, `{"filename":"logs/beta.log","separate":["error", "warning", "notice", "info", "debug"]}`)
	logs.EnableFuncCallDepth(true)
	// beego.AddViewPath("template")
	args := os.Args //获取用户输入的所有参数
	if args == nil || len(args) < 2 {
		//如果用户没有输入,或参数个数不够,则调用该函数提示用户
		//beego.Run()
		cmdUsage()
	} else if len(args) == 2 {
		switch args[1] {
		case "start":
			//runtime.GOMAXPROCS(runtime.NumCPU()) // 预置运行规格协程信息
			common.ListenGrpcStatus()
			toolbox.StartTask()      // 开始全局定时任务
			defer toolbox.StopTask() //
			beego.Run()
		case "orm":
			orm.RunCommand()
		default:
			cmdUsage()
		}
	} else if len(args) > 2 {
		if args[1] == "orm" {
			orm.RunCommand()
		} else {
			cmdUsage()
		}
	}
}

//cmdUsage 显示命令行帮助
func cmdUsage() {
	fmt.Println(`
		USAGE
			WeixinBot [commond]
		AVAILABLE COMMANDS
			start                     Start beta store server node.
			orm                       Operate the database.
			`)
}

//初始化公共目录
func createDIR() {
	var err error
	//初始化日志目录
	if _, err = os.Stat("logs"); err != nil {
		os.Mkdir("logs", os.ModePerm)
	}
}
