package models

import (
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/orm"
)

// AddConfiguration ...
func AddConfiguration(m *bridageModels.Configuration) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}
