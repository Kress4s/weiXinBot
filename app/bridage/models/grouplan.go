package models

// GroupPlan ...
type GroupPlan struct {
	ID     int64    `orm:"auto;column(id)"` //
	Type   int      `orm:"column(type)"`    // 群管类型
	Plan   int64    `orm:"column(plan_id)"` // 群管类型的ID
	Groups []*Group `orm:"rel(m2m)"`        // 多对多(群)
}
