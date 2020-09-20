package models

import (
	"github.com/astaxie/beego/orm"
)

// Manager ...
type Manager struct {
	ID       int64   `orm:"auto;column(id)"`
	IsAdmin  bool    `orm:"column(is_admin);default(1)"`
	Account  string  `orm:"size(50);column(account)" `
	PassWord string  `orm:"size(50);column(pwd)"`
	NickName string  `orm:"size(50);column(nickname)"`
	Tel      string  `orm:"size(50);column(tel)"`
	Avatar   string  `orm:"size(100);column(avatar)"`
	BotsNum  int     `orm:"column(botsnum)"`
	Bots     []*Bots `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Manager))
}
