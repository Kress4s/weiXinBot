package models

import (
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddGroup ...
func AddGroup(m *bridageModels.Group) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

//MultiAddGroup ...
func MultiAddGroup(m []*bridageModels.Group) (err error) {
	o := orm.NewOrm()
	for _, _m := range m {
		_, err = o.Insert(_m)
	}
	return
}

// GetGroupByGID ...
func GetGroupByGID(GID string) (v *bridageModels.Group, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Group{GID: GID}
	if err = o.Read(v, "GID"); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateGrouByID ...
func UpdateGrouByID(m *bridageModels.Group) (err error) {
	var v = bridageModels.Group{GID: m.GID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "IsNeedServe", "GroupPlan"); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}
