package models

import (
	taskflow "weiXinBot/app/bridage/flows/timetask"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddTimeTask ...
func AddTimeTask(m *bridageModels.TimeTask) (id int64, err error) {
	var timetasker taskflow.TaskFactory
	o := orm.NewOrm()
	// if m.SendType == -1 {
	// 	// 立即推送
	// 	if err = bridageModels.SendImmediately(m); err != nil {
	// 		logs.Error("[%s] SendImmediately failed, err is ", m.BotWXID)
	// 		m.Status = constant.FAILEDSEND
	// 	}
	// 	m.Status = constant.SENDED
	// 	m.SetUpTime = time.Now()
	// 	id, err = o.Insert(m)
	// } else {
	// 	if err = bridageModels.TimingSend(m); err != nil {
	// 		logs.Error("[%s] TimingSend failed, err is ", m.BotWXID, err.Error())
	// 		m.Status = constant.FAILEDSEND
	// 	}
	// 	m.Status = constant.SENDED
	// 	id, err = o.Insert(m)
	// }
	if timetasker, err = taskflow.GetTaskIns(m.Type); err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	if id, err = o.Insert(m); err != nil {
		logs.Error("insert task failed, err is ", err.Error())
	}
	switch m.SendType {
	case -1:
		go timetasker.SendImmediately(*m)
	}
	return
}
