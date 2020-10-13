package models

import (
	"github.com/astaxie/beego/orm"
)

// KeyWords ...
type KeyWords struct {
	ID        int64       `orm:"auto;column(id)"`             //
	Type      int         `orm:"column(type);default(2)"`     // 所属功能类型 （默认2）
	Switch    bool        `orm:"column(switch);default(1)"`   //功能总开关
	Content   string      `orm:"size(20);column(title)"`      //
	Resources string      `orm:"size(300); column(resouces)"` // 来自资源库的具体回复内容(ids, ","连接,有多个)
	Questions []*Question `orm:"reverse(many)"`               //
	GroupPlan *GroupPlan  `orm:"rel(fk)"`                     //
}

// Question 关键字回复的配置的问题
type Question struct {
	ID         int64        `orm:"auto;column(id)"`            //
	Title      string       `orm:"size(50);column(title)"`     //
	ExactWords []*ExactWord `orm:"reverse(many)"`              //
	FuzzWords  []*FuzzWord  `orm:"reverse(many)"`              //
	Resources  string       `orm:"size(30); column(resouces)"` // 来自资源库的具体回复内容(ids, ","连接,有多个)
	KeyWords   *KeyWords    `orm:"rel(fk)"`                    //
}

// ExactWord 精准
type ExactWord struct {
	ID       int64     `orm:"auto;column(id)"`        //
	Word     string    `orm:"size(20);column(words)"` //精准关键词内容
	Question *Question `orm:"rel(fk)"`                //
}

// FuzzWord 模糊
type FuzzWord struct {
	ID       int64     `orm:"auto;column(id)"`        //
	Word     string    `orm:"size(20);column(words)"` // 模糊关键词内容
	Question *Question `orm:"rel(fk)"`                //
}

func init() {
	orm.RegisterModel(new(KeyWords), new(ExactWord), new(FuzzWord), new(Question))
}
