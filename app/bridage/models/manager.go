package models

import (
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
func GetManagerByAccount(account string) (v *Manager, err error) {
	o := orm.NewOrm()
	v = &Manager{Tel: account}
	if err = o.Read(v, "Tel"); err != nil {
		return nil, err
	}
	return v, nil
}

// AddManager ...
func AddManager(manager *Manager) (id int64, err error) {
	o := orm.NewOrm()
	_, err = o.Insert(manager)
	return
}

//根据手机号查询是否存在注册记录
func FindManagerByTel(tel string) bool {

	o := orm.NewOrm()
	exist := o.QueryTable("manager").Filter("Tel", tel).Exist()

	return exist
}
