package models

import (
	"time"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddTimeTask ...
func AddTimeTask(m *bridageModels.TimeTask) (id int64, err error) {
	o := orm.NewOrm()
	if m.Type == -1 {
		// 立即推送
		if err = bridageModels.SendImmediately(m); err != nil {
			logs.Error("[%s] SendImmediately failed, err is ", m.BotWXID)
			m.Status = constant.FAILEDSEND
		}
		m.Status = constant.SENDED
		m.SetUpTime = time.Now()
		id, err = o.Insert(m)
	} else {
		if err = bridageModels.TimingSend(m); err != nil {
			logs.Error("[%s] TimingSend failed, err is ", m.BotWXID)
			m.Status = constant.FAILEDSEND
		}
		m.Status = constant.SENDED
		id, err = o.Insert(m)
	}
	return
}
