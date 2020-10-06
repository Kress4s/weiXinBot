package models

import (
	"github.com/astaxie/beego/orm"
)

// WhiteList ...
type WhiteList struct {
	ID        int64      `orm:"auto;column(id)"`               //
	WXID      string     `orm:"size(50);column(wx_id);unique"` // 微信id
	Type      int        `orm:"column(type);default(4)"`       // 所属功能类型 （默认4）
	GroupPlan *GroupPlan `orm:"rel(fk)"`                       //
}

func init() {
	orm.RegisterModel(new(WhiteList))
}
