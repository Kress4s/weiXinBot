package models

import (
	taskflow "weiXinBot/app/bridage/flows"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddTimeTask ...
func AddTimeTask(m *bridageModels.TimeTask) (id int64, err error) {
	var timetasker taskflow.TaskFactory
	o := orm.NewOrm()
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
