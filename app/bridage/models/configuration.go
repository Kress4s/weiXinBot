package models

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Configuration 功能设置(消息监听使用)
type Configuration struct {
	ID        int64  `orm:"auto;column(id)"`              //
	Type      int    `orm:"column(type)"`                 // 配置对象 0: 群组; 1: 联系人 ...可拓展
	FuncType  int    `orm:"column(function_type)"`        // 功能配置类型 1:入群欢迎语 2:关键词回复 3:自动踢人...可拓展
	FuncID    int64  `orm:"column(function_id)"`          // 配置ID
	BotID     string `orm:"size(30);column(bot_id)"`      // 机器人微信ID，执行消息回复、踢人等操作的微信号(保证机器人是正确的)
	ObjectIDS string `orm:"size(200);column(object_ids)"` // 要执行对象的IDs,多个用”,“连接(群、联系人...可拓展)
}

func init() {
	orm.RegisterModel(new(Configuration))
}

// GroupIsNeedServer 查看此群是否需要机器人服务
// @Param fromUserName: {Str:22475302355@chatroom}
// @Param toUserName: {Str:wxid_vao3ptfez4p22}}
func GroupIsNeedServer(fromUserName, toUserName string) (isServer bool, err error) {
	o := orm.NewOrm()
	// is group message
	if !strings.Contains(fromUserName, "@chatroom") {
		return false, nil
	}
	if !o.QueryTable(new(Configuration)).Filter("Type", 0).Filter("ObjectIDS__icontains", fromUserName).Filter("BotID", toUserName).Exist() {
		return false, nil
	}
	return true, nil
}

// GroupService 需要的服务
// keyContent: "push_content":"🛫张 : G吐总冠军"
func GroupService(fromUserName, toUserName, keyContent string) {
	o := orm.NewOrm()
	var err error
	var configs []*Configuration
	if _, err = o.QueryTable(new(Configuration)).Filter("Type", 0).Filter("ObjectIDS__icontains", fromUserName).Filter("BotID", toUserName).All(&configs); err != nil {
		logs.Error("get Configuration by fromUserName and toUserName failed, err is ", err.Error())
	}
	for _, v := range configs {
		switch v.FuncType {
		// is welcome function config
		case 1:
			//确定新人进群的数据结构再做处理
			fmt.Println("新人进来了")

		// is keywords function config
		case 2:
			fmt.Println("关键词回复")
		// is autokick function config
		case 3:
			fmt.Println("自动踢人")
		// is white function config
		case 4:
			fmt.Println("白名单")
		default:
			logs.Error("function config Type[%d] is not right, please cheak it and modify it", v.FuncType)
		}
	}
}
