package models

import (
	"github.com/astaxie/beego/orm"
)

// KeyWords ...
type KeyWords struct {
	ID         int64        `orm:"auto;column(id)"`           //
	Switch     bool         `orm:"column(switch);default(1)"` //功能总开关
	ExactWords []*ExactWord `orm:"reverse(many)"`             //
	FuzzWords  []*FuzzWord  `orm:"reverse(many)"`             //
	Content    string       `orm:"size(20);column(title)"`    //
}

// ExactWord 精准
type ExactWord struct {
	ID       int64     `orm:"auto;column(id)"`        //
	Words    string    `orm:"size(20);column(title)"` //精准关键词内容
	KeyWords *KeyWords `orm:"rel(fk)"`                //
}

// FuzzWord 模糊
type FuzzWord struct {
	ID       int64     `orm:"auto;column(id)"`        //
	Words    string    `orm:"size(20);column(title)"` // 模糊关键词内容
	KeyWords *KeyWords `orm:"rel(fk)"`                //
}

func init() {
	orm.RegisterModel(new(KeyWords), new(ExactWord), new(FuzzWord))
}
