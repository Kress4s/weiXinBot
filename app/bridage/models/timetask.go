package models

import (
	"time"
)

// TimeTask 定时任务
type TimeTask struct {
	ID        int64     `orm:"auto;column(id)"`                       //
	Switch    bool      `orm:"column(switch)"`                        // 开关
	Type      int       `orm:"column(type)"`                          // 类型(0:单次执行;1:按天发送;2:按周发送;3:按月发送;)
	BotWXID   string    `orm:"size(30);column(botwxid)"`              //
	Objects   string    `orm:"size(300);column(objects)"`             // 群组或者联系人
	Status    string    `orm:"size(20);column(status)"`               // 任务状态(UnSend;Sended)
	SetUpTime time.Time `orm:"type(datetime);column(setuptime);null"` // 设置发送时间
	Resource  string    `orm:"size(20);column(resource)"`             // 发送内容(多个)
}

// func init() {
// 	orm.RegisterModel(new(TimeTask))
// }
