package models

import (
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/orm"
)

// AddGroup ...
func AddGroup(m *bridageModels.Group) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
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
