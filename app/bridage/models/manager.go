package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Manager ...
type Manager struct {
	ID       int64   `orm:"auto;column(id)"`
	IsAdmin  bool    `orm:"column(is_admin);default(0)"`
	PassWord string  `orm:"size(50);column(pwd)"`
	Tel      string  `orm:"size(50);column(tel)"`
	Avatar   string  `orm:"size(100);column(avatar)"`
	BotsNum  int     `orm:"column(botsnum)"`
	Bots     []*Bots `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Manager))
}

// GetManagerByAccount ...
func GetManagerByAccount(tel string) (ret *Manager, err error) {
	o := orm.NewOrm()
	v := Manager{Tel: tel}
	if err = o.QueryTable(new(Manager)).Filter("Tel", tel).One(&v); err == nil {
		o.LoadRelated(&v, "Bots")
		return &v, nil
	} else if err == orm.ErrNoRows {
		err = fmt.Errorf("账号或密码错误")
		return nil, err
	}
	logs.Error("Get manager by account failed, err is ", err.Error())
	return nil, err
}

// AddManager ...
func AddManager(manager *Manager) (id int64, err error) {
	o := orm.NewOrm()
	_, err = o.Insert(manager)
	return
}

// FindManagerByTel 根据手机号查询是否存在注册记录
func FindManagerByTel(tel string) bool {
	o := orm.NewOrm()
	exist := o.QueryTable(new(Manager)).Filter("Tel", tel).Exist()
	return exist
}
