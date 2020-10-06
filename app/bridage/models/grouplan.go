package models

import "github.com/astaxie/beego/orm"

// GroupPlan ...
type GroupPlan struct {
	ID      int64    `orm:"auto;column(id)"`       //
	Name    string   `orm:"size(50);column(name)"` // 名称
	Manager *Manager `orm:"rel(fk)"`               // 属于哪个用户创建的群管方案
	Groups  []*Group `orm:"reverse(many)"`         //
}

func init() {
	orm.RegisterModel(new(GroupPlan))
}
