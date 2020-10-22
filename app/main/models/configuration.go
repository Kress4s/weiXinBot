package models

import (
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddConfiguration ...
func AddConfiguration(m *bridageModels.Configuration) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// UpdateConfigurationByID ...
func UpdateConfigurationByID(m *bridageModels.Configuration) (err error) {
	var v = bridageModels.Configuration{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "ObjectIDS"); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}
