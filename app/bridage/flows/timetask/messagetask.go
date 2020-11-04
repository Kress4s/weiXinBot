package timetask

import (
	"fmt"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
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
	// taskIns := toolbox.NewTask(fmt.Sprintf("%d", v.ID), bridageModels.SetUpTimeFormatString(v.SendType,
	// 	v.SetUpFormat), bridageModels.GenerateTask(v))
	fmt.Println(v)
	return

}

func init() {
	Register(constant.MESSAGE_TASK, &MessageTask{})
}
