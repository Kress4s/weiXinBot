package models

import (
	"github.com/astaxie/beego/orm"
)

// Welcome ...
type Welcome struct {
	ID                  int64            `orm:"auto;column(id)"`                           //
	Switch              bool             `orm:"column(switch);default(1)"`                 // 功能总开关
	WaitSeconds         int              `orm:"column(waitseconds)"`                       // 等待时间
	WaitTimeSwitch      bool             `orm:"column(wait_time_switch);default(1)"`       // 进群多少秒发送
	WaitNewersNumSwitch int              `orm:"column(wait_newers_num_switch);default(1)"` // 新进群多少人发送
	ReplyContents       []*SourceWelcome `orm:"reverse(many)"`                             // 空 前端默认展示
	Resources           string           `orm:"size(300); column(resouces)"`               // 从素材库导入的回复内容(ids, ","连接)
}

// SourceWelcome 欢迎语的回复内容(可以多个)
type SourceWelcome struct {
	ID      int64    `orm:"auto;column(id)"`
	Type    int      `orm:"column(type)"`               //消息内容类型(0:文本 1:图片....)
	Content string   `orm:"size(300); column(content)"` //回复内容
	Welcome *Welcome `orm:"rel(fk)"`                    //
}

func init() {
	orm.RegisterModel(new(Welcome), new(SourceWelcome))
}
