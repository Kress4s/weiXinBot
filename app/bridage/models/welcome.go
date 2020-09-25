package models

import (
	"github.com/astaxie/beego/orm"
)

// Welcome ...
type Welcome struct {
	ID                  int64  `orm:"auto;column(id)"`                           //
	Switch              bool   `orm:"column(switch);default(1)"`                 //功能总开关
	Words               string `orm:"size(200);column(words)"`                   //回复语
	Type                int    `orm:"column(type);default(0)"`                   //回复类型
	WaitSeconds         int    `orm:"column(waitseconds)"`                       //等待时间
	WaitTimeSwitch      bool   `orm:"column(wait_time_switch);default(1)"`       //进群多少秒发送
	WaitNewersNumSwitch int    `orm:"column(wait_newers_num_switch);default(1)"` //新进群多少人发送
	SourceURL           string `orm:"size(300); column(source_url)"`             //回复文件地址
}

func init() {
	orm.RegisterModel(new(Welcome))
}
