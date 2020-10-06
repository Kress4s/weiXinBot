package models

import "github.com/astaxie/beego/orm"

// Contact ...
type Contact struct {
	WXID           string     `orm:"pk;size(200);column(id)"`
	Bots           *Bots      `orm:"rel(fk)"`
	BigHeadImage   string     `orm:"size(200);column(head_big_image_url)"`
	SmallHeadImage string     `orm:"size(200);column(head_small_image_url)"`
	NickName       string     `orm:"size(50);column(nick_name)"`
	Country        string     `orm:"size(50);column(country)"`
	Province       string     `orm:"size(50);column(province)"`
	City           string     `orm:"size(50);column(city)"`
	Sex            bool       `orm:"column(sex)"`
	Signature      string     `orm:"size(50);column(signature)"`
	Alias          string     `orm:"size(50);column(alias_name)"`
	Messages       []*Message `orm:"reverse(many)"`
	// Labels         []*Label `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Contact))
}
