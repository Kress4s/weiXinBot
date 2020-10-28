package models

import (
	"fmt"
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
	IsAttachQuestion bool        `orm:"column(is_attach_question)"`  // 回复是否携带问题
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
func KeyWordsService(id int64, keyContent string) (isNeedReply bool, replyContent []*Resource, err error) {
	/*
		1. 判断开关
		2. 查找精准
		3. 匹配模糊
		4. 返回回复内容
	*/
	o := orm.NewOrm()
	var kWord = KeyWords{ID: id}
	if err = o.Read(&kWord); err != nil {
		logs.Error("KeyWordsService: get keyward failed, err is ", err.Error())
		return false, nil, err
	}
	// 开关关闭
	if !kWord.Switch {
		return false, nil, nil
	}
	var exWords []*ExactWord
	var num int64
	if num, err = o.QueryTable(new(ExactWord)).Filter("Question__KeyWords__ID", id).All(&exWords); err != nil {
		logs.Error("KeyWordsService: get all ExactWord by keyword ID failed, err is ", err.Error())
		return false, nil, err
	}
	keyContent = strings.ReplaceAll(keyContent, "\n", "")
	nameContent := strings.Split(keyContent, ":")
	// 设置了精准关键词
	if num > 0 {
		// 先精准匹配(匹配到直接模糊匹配)
		for _, ew := range exWords {
			if nameContent[1] != ew.Word {
				continue
			}
			// 匹配到(查回复内容)
			var question Question
			if err = o.QueryTable(new(Question)).Filter("ExactWords__ID", ew.ID).One(&question); err != nil {
				logs.Error("KeyWordsService: get question by ExactWords__ID failed, err is ", err.Error())
				return false, nil, err
			}
			// 查看是否设置回复资源
			if len(question.Resources) == 0 {
				return false, nil, nil
			}
			// 查找回复的资源
			var replyresource []*Resource
			if replyresource, err = GetResourceByIds(question.Resources); err != nil {
				logs.Error("KeyWordsService: get Get replyresource by Resources failed, err is ", err.Error())
				return false, nil, err
			}
			//检查配置设置
			if kWord.IsAt && kWord.IsAttachQuestion {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = strings.ReplaceAll(replyresource[i].Material[j].Data, "{{@新人}}", "@"+nameContent[0])
							replyresource[i].Material[j].Data = fmt.Sprintf("@%s%s\n\n%s", nameContent[0], nameContent[1], replyresource[i].Material[j].Data)
						}
					}
				}
			} else if kWord.IsAt && !kWord.IsAttachQuestion {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = fmt.Sprintf("@%s\n\n%s", nameContent[0], replyresource[i].Material[j].Data)
						}
					}
				}
			} else {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = fmt.Sprintf("%s\n%s", nameContent[1], replyresource[i].Material[j].Data)
						}
					}
				}
			}
			return true, replyresource, nil
		}
	}
	var fuzzWords []*FuzzWord
	if num, err = o.QueryTable(new(FuzzWord)).Filter("Question__KeyWords__ID", id).All(&fuzzWords); err != nil {
		logs.Error("KeyWordsService: get all FuzzWord by keyword ID failed, err is ", err.Error())
		return false, nil, err
	}
	// 设置了模糊关键词
	if num > 0 {
		// 先模糊匹配(匹配到直接模糊匹配)
		for _, fw := range fuzzWords {
			fmt.Println(fw.Word, nameContent[1])
			if !strings.Contains(nameContent[1], fw.Word) {
				continue
			}
			// 匹配到(查回复内容)
			var question Question
			if err = o.QueryTable(new(Question)).Filter("FuzzWords__ID", fw.ID).One(&question); err != nil {
				logs.Error("KeyWordsService: get question by FuzzWords__ID failed, err is ", err.Error())
				return false, nil, err
			}
			// 查看是否设置回复资源
			if len(question.Resources) == 0 {
				return false, nil, nil
			}
			// 查找回复的资源
			var replyresource []*Resource
			if replyresource, err = GetResourceByIds(question.Resources); err != nil {
				logs.Error("KeyWordsService: get Get replyresource by Resources failed, err is ", err.Error())
				return false, nil, err
			}
			if kWord.IsAt && kWord.IsAttachQuestion {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = strings.ReplaceAll(replyresource[i].Material[j].Data, "{{@新人}}", "@"+nameContent[0])
							replyresource[i].Material[j].Data = fmt.Sprintf("@%s%s\n\n%s", nameContent[0], nameContent[1], replyresource[i].Material[j].Data)
						}
					}
				}
			} else if kWord.IsAt && !kWord.IsAttachQuestion {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = fmt.Sprintf("@%s\n%s", nameContent[0], replyresource[i].Material[j].Data)
						}
					}
				}
			} else {
				for i := range replyresource {
					for j := range replyresource[i].Material {
						if replyresource[i].Material[j].Type == 1 {
							replyresource[i].Material[j].Data = fmt.Sprintf("%s\n%s", nameContent[1], replyresource[i].Material[j].Data)
						}
					}
				}
			}
			return true, replyresource, nil
		}
	}
	return false, nil, nil
}
