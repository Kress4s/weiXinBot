package models

import (
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// KeyWords ...
type KeyWords struct {
	ID               int64       `orm:"auto;column(id)"`             //
	Type             int         `orm:"column(type);default(2)"`     // 所属功能类型 （默认2）
	Switch           bool        `orm:"column(switch);default(1)"`   //功能总开关
	IsAt             bool        `orm:"column(isAt)"`                // 回复是否@对方
	IsAttachQuestion bool        `orm:"Isattachquestion"`            // 回复是否携带问题
	Resources        string      `orm:"size(300); column(resouces)"` // 来自资源库的具体回复内容(ids, ","连接,有多个)
	Questions        []*Question `orm:"reverse(many)"`               //
	GroupPlan        *GroupPlan  `orm:"rel(fk)"`                     //
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
	ID       int64     `orm:"auto;column(id)"`       //
	Word     string    `orm:"size(20);column(word)"` //精准关键词内容
	Question *Question `orm:"rel(fk)"`               //
}

// FuzzWord 模糊
type FuzzWord struct {
	ID       int64     `orm:"auto;column(id)"`       //
	Word     string    `orm:"size(20);column(word)"` // 模糊关键词内容
	Question *Question `orm:"rel(fk)"`               //
}

func init() {
	orm.RegisterModel(new(KeyWords), new(ExactWord), new(FuzzWord), new(Question))
}

// MultiDeleteFuzzWordByIDs multi delete FuzzWord
func MultiDeleteFuzzWordByIDs(ids string) (err error) {
	var idslice []interface{}
	o := orm.NewOrm()
	s := strings.Split(ids, ",")
	for _, v := range s {
		idslice = append(idslice, v)
	}
	var num int64
	if num, err = o.QueryTable(new(FuzzWord)).Filter("ID__in", idslice...).Delete(); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}

// MultiDeleteExactWordByIDs multi delete ExactWord
func MultiDeleteExactWordByIDs(ids string) (err error) {
	var idslice []interface{}
	o := orm.NewOrm()
	s := strings.Split(ids, ",")
	for _, v := range s {
		idslice = append(idslice, v)
	}
	var num int64
	if num, err = o.QueryTable(new(ExactWord)).Filter("ID__in", idslice...).Delete(); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}

// KeyWordsService ...
// @Params  keyContent: "push_content":"🛫张 : G吐总冠军"
func KeyWordsService(keyContent string) (isNeedReply bool, replyContent []interface{}) {
	/*
		1. 判断开关
		2. 查找精准
		3. 匹配模糊
		4. 是否@; 是否attach上问题
		5. 返回回复内容
	*/
	return
}
