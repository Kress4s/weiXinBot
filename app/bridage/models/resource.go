package models

import (
	"github.com/astaxie/beego/orm"
)

// Resource send source
// 素材库
type Resource struct {
	ID   int64  `orm:"auto;column(id)"`
	Type int    `orm:"column(type)"`            //消息内容类型
	Text string `orm:"type(text)"`              //消息文本
	Data string `orm:"size(300); column(data)"` //其他类型的地址(后期可拆)
}

func init() {
	orm.RegisterModel(new(Resource))
}
