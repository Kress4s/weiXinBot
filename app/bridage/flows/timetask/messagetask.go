package timetask

import (
	"fmt"
	"weiXinBot/app/bridage/constant"
	flow "weiXinBot/app/bridage/flows"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
)

// MessageTask 消息任务
type MessageTask struct{}

// SendImmediately ...
func (c *MessageTask) SendImmediately(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	defer func() {
		if verr := recover(); verr != nil {
			logs.Error("AnnouncementTask SendImmediately: err is ", verr)
		}
	}()
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	if err = bridageModels.ExecuteTask(v); err != nil {
		logs.Error("MessageTask SendImmediately send failed, err is ", err.Error())
	}
	return
}

// TimingSend 定时发送
func (c *MessageTask) TimingSend(p interface{}) (err error) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	defer func() {
		if verr := recover(); verr != nil {
			logs.Error("AnnouncementTask SendImmediately: err is ", verr)
		}
	}()
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	taskIns := toolbox.NewTask(fmt.Sprintf("task-%d", v.ID), bridageModels.SetUpTimeFormatString(v.SendType,
		v.SetUpFormat), func() error {
		return bridageModels.GenerateTask(v)
	})
	toolbox.AddTask(fmt.Sprintf("task-%d", v.ID), taskIns)
	return
}

// ModifyTimeTask ...
func (c *MessageTask) ModifyTimeTask(p interface{}) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	// 删除原来的任务
	toolbox.DeleteTask(fmt.Sprintf("task-%d", v.ID))
	// 新建新的任务
	taskIns := toolbox.NewTask(fmt.Sprintf("task-%d", v.ID), bridageModels.SetUpTimeFormatString(v.SendType,
		v.SetUpFormat), func() error {
		return bridageModels.GenerateTask(v)
	})
	toolbox.AddTask(fmt.Sprintf("task-%d", v.ID), taskIns)
	return
}

// DeleteTimeTask ...
func (c *MessageTask) DeleteTimeTask(p interface{}) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	toolbox.DeleteTask(fmt.Sprintf("task-%d", v.ID))
	return
}

func init() {
	flow.Register(constant.MESSAGE_TASK, &MessageTask{})
}
