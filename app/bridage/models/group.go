package models

import (
	"github.com/astaxie/beego/orm"
)

// Group ...
type Group struct {
	GID            string     `orm:"pk;size(50);column(g_id)" `              // json:wx_id
	NickName       string     `orm:"size(50);column(nick_name)" `            //
	Owner          string     `orm:"size(50);column(owner)" `                //群主
	MemberNum      int        `orm:"column(member_num)"`                     //
	HeadSmallImage string     `orm:"size(200);column(head_small_image_url)"` //
	Listers        string     `orm:"size(500);column(listers)"`              //成员微信号的IDs，”，“连接 接口返回值[]不好处理 记录1
	IsNeedServe    bool       `orm:"column(isneedserve);default(0)"`         // 是否有服务功能
	Bots           *Bots      `orm:"rel(fk)"`                                //
	GroupPlan      *GroupPlan `orm:"null;rel(fk)"`                           //群方案
	Messages       []*Message `orm:"reverse(many)"`                          //
}

func init() {
	orm.RegisterModel(new(Group))
}
