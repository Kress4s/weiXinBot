package timetask

import (
	"fmt"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
)

// AnnouncementTask 公告任务
type AnnouncementTask struct{}

// SendImmediately ...
func (c *AnnouncementTask) SendImmediately(p interface{}) (err error) {
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
		logs.Error("AnnouncementTask SendImmediately send failed, err is ", err.Error())
	}
	return
}

// TimingSend 定时发送
func (c *AnnouncementTask) TimingSend(p interface{}) (err error) {
	// 定时任务
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
	fmt.Println(v)
	return err
}

func init() {
	Register(constant.ANNOUNCEMENT_TASK, &AnnouncementTask{})
}
