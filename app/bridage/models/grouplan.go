package models

import "github.com/astaxie/beego/orm"

// GroupPlan ...
type GroupPlan struct {
	ID      int64    `orm:"auto;column(id)"`       //
	Name    string   `orm:"size(50);column(name)"` // 名称
	PlanIds string   `orm:"column(plan_id)"`       // 群管类型的IDs(用“，”连接) 接口加个type参数，指定id的位置对应关系，点的是哪个配置的id（0:欢迎新人；1:关键字新人...）
	Manager *Manager `orm:"rel(fk)"`               // 属于哪个用户创建的群管方案
	Groups  []*Group `orm:"reverse(many)"`         //
}

func init() {
	orm.RegisterModel(new(GroupPlan))
}
