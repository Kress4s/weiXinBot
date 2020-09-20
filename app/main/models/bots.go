package models

import (
	bridage "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/orm"
)

// AddBot ...
func AddBot(bot *bridage.Bots) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(bot)
	return
}
