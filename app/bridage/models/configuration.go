package models

import "github.com/astaxie/beego/orm"

// Configuration 功能设置(消息监听使用)
type Configuration struct {
	ID        int64  `orm:"auto;column(id)"`              //
	Type      int    `orm:"column(type)"`                 // 配置对象 0: 群组; 1: 联系人 ...可拓展
	FuncType  int    `orm:"column(function_type)"`        // 功能配置类型 1:入群欢迎语 2:关键词回复 3:自动踢人...可拓展
	FuncID    int64  `orm:"column(function_id)"`          // 配置ID
	BotID     string `orm:"size(30);column(bot_id)"`      // 机器人微信ID，执行消息回复、踢人等操作的微信号
	ObjectIDS string `orm:"size(200);column(object_ids)"` // 要执行对象的IDs,多个用”,“连接(群、联系人...可拓展)
}

func init() {
	orm.RegisterModel(new(Configuration))
}
