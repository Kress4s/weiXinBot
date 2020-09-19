package models

import (
	bridage "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/orm"
)

// AddUser ...
func AddUser(user *bridage.User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(user)
	return
}
