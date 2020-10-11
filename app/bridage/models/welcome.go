package models

import (
	"github.com/astaxie/beego/orm"
)

// Welcome ...
type Welcome struct {
	ID                  int64      `orm:"auto;column(id)"`                           //
	Type                int        `orm:"column(type);default(1)"`                   // 所属功能类型 （默认1）
	Switch              bool       `orm:"column(switch);default(1)"`                 // 功能总开关
	WaitSeconds         int        `orm:"column(waitseconds)"`                       // 等待时间
	WaitTimeSwitch      bool       `orm:"column(wait_time_switch);default(1)"`       // 进群多少秒发送
	NewerNum            int        `orm:"column(newernum)"`                          //新进到多少人发送
	WaitNewersNumSwitch bool       `orm:"column(wait_newers_num_switch);default(1)"` // 新进群多少人发送开关
	Resources           string     `orm:"size(300); column(resouces)"`               // 来自资源库的具体回复内容(ids, ","连接,有多个,有素材库导入)
	GroupPlan           *GroupPlan `orm:"rel(fk)"`                                   //
}

func init() {
	orm.RegisterModel(new(Welcome))
}
