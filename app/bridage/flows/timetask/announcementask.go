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

// AnnouncementTask 公告任务
type AnnouncementTask struct{}

// TaskImmediately ...
func (c *AnnouncementTask) TaskImmediately(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	defer func() {
		if verr := recover(); verr != nil {
			logs.Error("AnnouncementTask SendImmediately: err is ", verr)
		}
	}()
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Announcement SendImmediately: v is not TimeTask struct")
	}
	if err = c.TaskGenerate(v); err != nil {
		logs.Error("AnnouncementTask SendImmediately send failed, err is ", err.Error())
	}
	return
}

// TaskGenerate ...
func (c *AnnouncementTask) TaskGenerate(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("Announcement TaskGenerate: v is not TimeTask struct")
	}
	if err = c.TaskExcute(v); err != nil {
		logs.Error("GenerateTask: taskID[%v], err is ", v.ID, err.Error())
	}
	return
}

// TaskSetting 定时发送
func (c *AnnouncementTask) TaskSetting(p interface{}) (err error) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	defer func() {
		if verr := recover(); verr != nil {
			logs.Error("AnnouncementTask TaskSetting: err is ", verr)
		}
	}()
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("AnnouncementTask TaskSetting: v is not TimeTask struct")
	}
	taskIns := toolbox.NewTask(fmt.Sprintf("task-%d", v.ID), bridageModels.SetUpTimeFormatString(v.SendType,
		v.SetUpFormat), func() error {
		return c.TaskGenerate(v)
	})
	toolbox.AddTask(fmt.Sprintf("task-%d", v.ID), taskIns)
	return err
}

// TaskExcute ...
func (c *AnnouncementTask) TaskExcute(p interface{}) (err error) {
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("AnnouncementTask TaskExcute: v is not TimeTask struct")
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
		tr.BotNickName = v.BotNickName
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
						logs.Info("发送群公告")
						if err = bridageModels.SendAnnouncement(v.BotWXID, st, m.Data); err != nil {
							return
						}
						// 停 Interval 秒
						time.Sleep(time.Duration(v.Interval) * time.Second)
					}
				case 2:
					// 图片信息
					logs.Error("AnnouncementTask: task[%d] have images, It is a error", v.ID)
				default:
					fmt.Println("无意义类型")
				}
			}
		}
	}
	return
}

// ModifyTimeTask ...
func (c *AnnouncementTask) ModifyTimeTask(p interface{}) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("AnnouncementTask ModifyTimeTask: v is not TimeTask struct")
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
func (c *AnnouncementTask) DeleteTimeTask(p interface{}) {
	// 定时任务
	var v bridageModels.TimeTask
	var ok bool
	if v, ok = p.(bridageModels.TimeTask); !ok {
		panic("AnnouncementTask DeleteTimeTask: v is not TimeTask struct")
	}
	toolbox.DeleteTask(fmt.Sprintf("task-%d", v.ID))
	return
}

// TasksHooked ...
func (c *AnnouncementTask) TasksHooked(p []interface{}) {
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
	flow.Register(constant.ANNOUNCEMENT_TASK, &AnnouncementTask{})
}
