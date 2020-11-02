package models

import (
	"fmt"
	"strings"
	"time"
)

// TimeTask 定时任务
type TimeTask struct {
	ID         int64     `orm:"auto;column(id)"`                       //
	Title      string    `orm:"size(30);column(title)"`                // 推送内容的标题
	Switch     bool      `orm:"column(switch)"`                        // 开关
	Type       int       `orm:"column(type)"`                          // 类型(-1:立刻推送;0:单次执行;1:按天发送;2:按周发送;3:按月发送;)
	BotWXID    string    `orm:"size(30);column(botwxid)"`              //
	ObjectsIDS string    `orm:"size(300);column(objects)"`             // 群组或者联系人
	Status     string    `orm:"size(20);column(status)"`               // 任务状态(UnSend;Sended;)
	SetUpTime  time.Time `orm:"type(datetime);column(setuptime);null"` // 设置发送时间
	Resource   string    `orm:"size(20);column(resource)"`             // 发送内容(多个)
	Remark     string    `orm:"size(50);column(remark)"`               // 任务备注
}

// func init() {
// 	orm.RegisterModel(new(TimeTask))
// }

// SendImmediately 立即推送
func SendImmediately(tt *TimeTask) (err error) {
	var resources []*Resource
	var sendTo []string
	if sendTo = strings.Split(tt.ObjectsIDS, ","); len(sendTo) == 0 {
		err = fmt.Errorf("sendTo is null")
		return
	}
	/*
		立即推送
		1. 查看资源，找出要发送的资源类型
		2. 发送，选择对应的方法
	*/
	if resources, err = GetResourceByIds(tt.Resource); err == nil && len(resources) > 0 {
		for _, r := range resources {
			for _, m := range r.Material {
				switch m.Type {
				case 1:
					// 文本信息
					// switch tt.Type {
					// case -1: // 立即推送
					for _, st := range sendTo {
						if err = SendText(tt.BotWXID, st, m.Data); err != nil {
							return
						}
						// 停0.5秒
						time.Sleep(500 * time.Microsecond)
					}
				case 2:
					// 图片信息
					for _, st := range sendTo {
						if err = SendImage(tt.BotWXID, st, m.Data); err != nil {
							return
						}
						// 停0.5秒
						time.Sleep(500 * time.Microsecond)
					}
				default:
					fmt.Println("未定义类型")
				}
			}
		}
	}
	return
}

// TimingSend 定时发送
func TimingSend(tt *TimeTask) (err error) {
	var resources []*Resource
	var sendTo []string
	if sendTo = strings.Split(tt.ObjectsIDS, ","); len(sendTo) == 0 {
		err = fmt.Errorf("sendTo is null")
		return
	}
	if resources, err = GetResourceByIds(tt.Resource); err == nil && len(resources) > 0 {
		for _, r := range resources {
			for _, _m := range r.Material {
				switch tt.Type {
				case 0:
					fmt.Println("单次执行")
					fmt.Println(_m.Data)
					// 定时任务

				case 1:
					fmt.Println("按天发送")
					// 定时任务

				case 2:
					fmt.Println("按周发送")
					// 定时任务

				case 3:
					fmt.Println("按月发送")
					// 定时任务

				default:
					fmt.Println("未知任务类型")
				}
			}
		}
	}
	return
}
