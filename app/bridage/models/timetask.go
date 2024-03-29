package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// TimeTask 定时任务
type TimeTask struct {
	ID          int64  `orm:"auto;column(id)"`              //
	Title       string `orm:"size(50);column(title)"`       // 推送内容的标题
	Type        string `orm:"size(30);column(type)"`        // 任务内容类型(messagetask:发消息;announcementask:公告...)
	Switch      bool   `orm:"column(switch);default(true)"` // 开关
	Interval    int    `orm:"column(interval)"`             // 发送多群间隔时间
	SendType    int    `orm:"column(sendtype)"`             // 类型(-1:立刻推送; 0:间隔时间执行; 1:单次执行; 2:按天发送; 3:按周发送; 4:按月发送;)
	SetUpFormat string `orm:"size(20);column(setupformat)"` // 设置定时格式的表达式
	BotWXID     string `orm:"size(30);column(botwxid)"`     // 设置的发送的微信号
	BotNickName string `orm:"size(50);column(botnickname)"` // 发送信号的昵称
	ObjectsIDS  string `orm:"size(1000);column(objectids)"` // 群组或者联系人
	Manager     string `orm:"size(30);column(manager)"`     // 属于哪个用户的任务(用户Tel)
	Resource    string `orm:"size(20);column(resource)"`    // 发送内容(多个)
	Remark      string `orm:"size(50);column(remark)"`      // 任务备注
}

// TimeTaskRecode ...
type TimeTaskRecode struct {
	ID          int64     `orm:"auto;column(id)"`                      //
	Title       string    `orm:"size(50);column(title)"`               // 推送内容的标题
	Type        string    `orm:"size(30);column(type)"`                // 任务内容类型(message:发消息;announcement:公告...)
	BotWXID     string    `orm:"size(30);column(botwxid)"`             // 设置的发送的微信号
	BotNickName string    `orm:"size(50);column(botnickname)"`         // 发送信号的昵称
	ObjectsIDS  string    `orm:"size(300);column(objectids)"`          // 群组或者联系人
	SendTime    time.Time `orm:"type(datetime);column(sendtime);null"` // 设置发送时间
	Status      string    `orm:"size(20);column(status)"`              // 任务状态(UnSend;Sended;)
	Manager     string    `orm:"size(30);column(manager)"`             // 属于哪个用户的任务(用户Tel)
	Resource    string    `orm:"size(20);column(resource)"`            // 发送内容(多个)
	Remark      string    `orm:"size(50);column(remark)"`              // 任务备注
}

func init() {
	orm.RegisterModel(new(TimeTask), new(TimeTaskRecode))
}

// AddTimeTaskRecode ...
func AddTimeTaskRecode(v *TimeTaskRecode) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(v)
	return
}

// GetTimeTaskByID ...
func GetTimeTaskByID(id int64) (v *TimeTask, err error) {
	o := orm.NewOrm()
	v = &TimeTask{ID: id}
	if err = o.Read(v); err != nil {
		logs.Error("GetTimeTaskByID read timetask failed, err is ", err.Error())
		return nil, err
	}
	return v, nil
}

// GetHMSBySecond ...
func GetHMSBySecond(p string) string {
	s, _ := strconv.Atoi(p)
	var m int
	var h int
	// s, _ := strconv.ParseInt(p, 0, 64)
	if s >= 60 {
		// _s := s % 60 // 秒数
		m = s / 60 // 分钟数
		if m >= 60 {
			h = m / 60
			// _m := m / 60
			return fmt.Sprintf("* * */%v", h)
		}
		return fmt.Sprintf("* */%v *", m)
	}
	return fmt.Sprintf("0/%v * *", s)
}

// SetUpTimeFormatString for toolbox time set format string
func SetUpTimeFormatString(sendType int, setUpString string) (SetUpFormat string) {
	switch sendType {
	case 0:
		/*
			间隔时间
			秒数：60的倍数
		*/
		subformat := GetHMSBySecond(setUpString)
		logs.Debug("SetUpFormat case 0: %s", fmt.Sprintf("%s * * *", subformat))
		return fmt.Sprintf("%s * * *", subformat)
	case 1:
		/*
			单次执行
			01-01 00:00
		*/
		reslice := strings.Split(setUpString, " ")
		if strings.Contains(reslice[1], ":") && strings.Contains(reslice[0], "-") {
			MD := strings.Split(reslice[0], "-")
			hm := strings.Split(reslice[1], ":")
			if len(MD) == 2 && len(hm) == 2 {
				logs.Debug("SetUpFormat case 1: %s", fmt.Sprintf("00 %s %s %s %s *", hm[1], hm[0], MD[1], MD[0]))
				return fmt.Sprintf("00 %s %s %s %s *", hm[1], hm[0], MD[1], MD[0])
			}
		}
	case 2:
		/*
			按天发送(每天)
			00:00
		*/
		hm := strings.Split(setUpString, ":")
		if len(hm) == 2 {
			logs.Debug("SetUpFormat case 2: %s", fmt.Sprintf("00 %s %s */1 * *", hm[1], hm[0]))
			return fmt.Sprintf("00 %s %s */1 * *", hm[1], hm[0])
		}
	case 3:
		/*
			每周几时间点发送
			0,1,2,3,4,5,6 00:00
		*/
		reslice := strings.Split(setUpString, " ")
		if len(reslice) == 2 {
			hm := strings.Split(reslice[1], ":")
			if len(hm) == 2 {
				logs.Debug("SetUpFormat case 3: %s", fmt.Sprintf("00 %s %s * * %s", hm[1], hm[0], reslice[0]))
				return fmt.Sprintf("00 %s %s * * %s", hm[1], hm[0], reslice[0])
			}
		}
	case 4:
		/*
			按月发送(每月的 哪些日期 时间点发送)
			1,2,3 00:00
		*/
		reslice := strings.Split(setUpString, " ")
		if len(reslice) == 2 {
			hm := strings.Split(reslice[1], ":")
			if len(hm) == 2 {
				logs.Debug("SetUpFormat case 4: %s", fmt.Sprintf("00 %s %s * * %s", hm[1], hm[0], reslice[0]))
				return fmt.Sprintf("00 %s %s %s */1 *", hm[1], hm[0], reslice[0])
			}
		}
	}
	return
}
