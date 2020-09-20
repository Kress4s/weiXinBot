package models

import (
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/orm"
)

// AddBot ...
func AddBot(bot *bridageModels.Bots) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(bot)
	return
}

// GetBotByID ...
func GetBotByID(id int64) (v *bridageModels.Bots, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Bots{ID: id}
	if err = o.Read(v); err != nil {
		return nil, err
	}
	return v, nil
}
