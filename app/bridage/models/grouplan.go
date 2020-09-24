package models

import "github.com/astaxie/beego/orm"

// GroupPlan ...
type GroupPlan struct {
	ID     int64    `orm:"auto;column(id)"`       //
	Name   string   `orm:"size(50);column(name)"` //
	Type   int      `orm:"column(type)"`          // 群管类型
	Plan   int64    `orm:"column(plan_id)"`       // 群管类型的ID
	Groups []*Group `orm:"reverse(many)"`         // 多对多(群)
}

func init() {
	orm.RegisterModel(new(GroupPlan))
}
