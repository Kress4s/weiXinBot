package models

import "github.com/astaxie/beego/orm"

// GroupPlan ...
type GroupPlan struct {
	ID   int64  `orm:"auto;column(id)"`       //
	Name string `orm:"size(50);column(name)"` //
	// Type    int      `orm:"column(type)"`          // 群管类型(0:欢迎新人；1:关键字新人... )
	PlanIds string   `orm:"column(plan_id)"` // 群管类型的IDs(用“，”连接)
	Manager *Manager `orm:"rel(fk)"`         //属于哪个用户创建的群管方案
	Groups  []*Group `orm:"reverse(many)"`   // 多对多(群)
}

func init() {
	orm.RegisterModel(new(GroupPlan))
}
