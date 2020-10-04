package models

import (
	"github.com/astaxie/beego/orm"
)

// Resource send source
// 素材库(单个组)
type Resource struct {
	ID       int64       `orm:"auto;column(id)"`              //
	IsPublic int         `orm:"column(is_public);default(0)"` //资源是否是公有（区分素材公共库展示和增加方案功能的具体回复内容）
	Material []*Material `orm:"reverse(many)"`                //
	Manager  string      `orm:"size(30);column(manager)"`     //属于哪个用户的素材库(用户Tel)
}

// Material 单个素材类型
type Material struct {
	ID       int64     `orm:"auto;column(id)"`         //
	Type     int       `orm:"column(type)"`            //消息内容类型(0:文本 1:图片....)
	Data     string    `orm:"size(300); column(data)"` //类型的具体内容(后期可拆)
	Resource *Resource `orm:"rel(fk)"`                 //所属素材组
}

func init() {
	orm.RegisterModel(new(Resource), new(Material))
}
