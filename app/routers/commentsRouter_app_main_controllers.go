package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:BotsController"],
        beego.ControllerComments{
            Method: "DeleteList",
            Router: `/deletelist`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GrouPlanController"],
        beego.ControllerComments{
            Method: "DeleteList",
            Router: `/deletelist`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GroupController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GroupController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GroupController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:GroupController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:gid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:ResouceController"],
        beego.ControllerComments{
            Method: "DeleteList",
            Router: `/deletelist`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: `/:id`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:id`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:id`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "GetOneWelcome",
            Router: `/category`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"] = append(beego.GlobalControllerRouter["weiXinBot/app/main/controllers:WelcomeController"],
        beego.ControllerComments{
            Method: "DeleteList",
            Router: `/deletelist`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
