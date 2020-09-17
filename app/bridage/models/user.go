package models

import (
	"github.com/astaxie/beego/orm"
)

// User ...
type User struct {
	ID             int64      `orm:"auto;column(id)"`
	WXID           string     `orm:"size(200);column(wx_id)" json:"id"`
	BigHeadImage   string     `orm:"size(200);column(big_head_img_url)" json:"big_head_img_url"`
	SmallHeadImage string     `orm:"size(200);column(small_head_img_url)" json:"small_head_img_url"`
	NickName       string     `orm:"size(50);column(nick_name)" json:"nick_name"`
	Country        string     `orm:"size(50);column(country)" json:"country"`
	Province       string     `orm:"size(50);column(province)" json:"province"`
	City           string     `orm:"size(50);column(city)" json:"city"`
	Sex            bool       `orm:"column(sex)" json:"sex"`
	Signature      string     `orm:"size(50);column(signature)" json:"signature"`
	Alias          string     `orm:"size(50);column(alias)" json:"alias"`
	Contacts       []*Contact `orm:"reverse(many)"` //好友
	Groups         []*Group   `orm:"reverse(many)"`
}

// TableIndex ...
// 多字段索引
func (u *User) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "WXID"},
	}
}

func init() {
	orm.RegisterModel(new(User))
}
