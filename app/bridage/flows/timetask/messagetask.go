package timetask

import (
	"fmt"
	"strings"
	"time"
	"weiXinBot/app/bridage/constant"
	flow "weiXinBot/app/bridage/flows"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
)

// MessageTask 消息任务
type MessageTask struct{}

// TaskImmediately ...
func (c *MessageTask) TaskImmediately(p interface{}) (err error) {
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
	if err = c.TaskGenerate(v); err != nil {
		logs.Error("MessageTask SendImmediately send failed, err is ", err.Error())
	}
	return
}

// TaskGenerate ...
func (c *MessageTask) TaskGenerate(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	if err = c.TaskExcute(v); err != nil {
		logs.Error("GenerateTask: taskID[%v], err is %s", v.ID, err.Error())
	}
	return
}

// TaskSetting 定时发送
func (c *MessageTask) TaskSetting(p interface{}) (err error) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	defer func() {
		if verr := recover(); verr != nil {
			logs.Error("Message SendImmediately: err is ", verr)
		}
	}()
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	taskIns := toolbox.NewTask(fmt.Sprintf("task-%d", v.ID), bridageModels.SetUpTimeFormatString(v.SendType,
		v.SetUpFormat), func() error {
		return c.TaskGenerate(v)
	})
	logs.Info("设置定时任务 ", bridageModels.SetUpTimeFormatString(v.SendType, v.SetUpFormat))
	toolbox.AddTask(fmt.Sprintf("task-%d", v.ID), taskIns)
	return
}

// TaskExcute ...
func (c *MessageTask) TaskExcute(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Message SendImmediately: v is not TimeTask struct")
	}
	var resources []*bridageModels.Resource
	var sendTo []string
	defer func() {
		tr := new(bridageModels.TimeTaskRecode)
		tr.Type = v.Type
		tr.Title = v.Title
		tr.SendTime = time.Now()
		tr.ObjectsIDS = v.ObjectsIDS
		tr.BotWXID = v.BotWXID
		tr.Resource = v.Resource
		tr.Manager = v.Manager
		if err != nil {
			tr.Status = constant.FAILEDSEND
			tr.Remark = err.Error()
		} else {
			tr.Status = constant.SENDED
			tr.Remark = "send successful"
		}
		bridageModels.AddTimeTaskRecode(tr)
	}()
	if sendTo = strings.Split(v.ObjectsIDS, ","); len(sendTo) == 0 {
		err = fmt.Errorf("sendTo is null")
		return
	}
	if resources, err = bridageModels.GetResourceByIds(v.Resource); err == nil {
		for _, r := range resources {
			for _, m := range r.Material {
				switch m.Type {
				case 1:
					// 文本信息
					for _, st := range sendTo {
						fmt.Println("开始发送消息")
						if err = bridageModels.SendText(v.BotWXID, st, m.Data); err != nil {
							return
						}
						// 停 Interval 秒
						time.Sleep(time.Duration(v.Interval) * time.Second)
					}
				case 2:
					// 图片信息
					for _, st := range sendTo {
						fmt.Println("开始发送图片")
						if err = bridageModels.SendImage(v.BotWXID, st, m.Data); err != nil {
							return
						}
						// 停 Interval 秒
						time.Sleep(time.Duration(v.Interval) * time.Second)
					}
				default:
					fmt.Println("未定义类型")
				}
			}
		}
	}
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
		return c.TaskGenerate(v)
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

// TasksHooked ...
func (c *MessageTask) TasksHooked(p []interface{}) {
	var err error
	for _, v := range p {
		if err = c.TaskSetting(v); err != nil {
			logs.Error("TaskSetting, accoured err is ", err.Error())
			continue
		}
	}
	return
}

func init() {
	flow.Register(constant.MESSAGE_TASK, &MessageTask{})
}
