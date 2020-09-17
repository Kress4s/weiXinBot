package models

import "github.com/astaxie/beego/orm"

// Contact ...
type Contact struct {
	WXID           string `orm:"pk;size(200);column(id)" json:"id"`
	User           *User  `orm:"rel(fk)"`
	BigHeadImage   string `orm:"size(200);column(head_big_image_url)" json:"head_big_image_url"`
	SmallHeadImage string `orm:"size(200);column(head_small_image_url)" json:"head_small_image_url"`
	NickName       string `orm:"size(50);column(nick_name)" json:"nick_name"`
	Country        string `orm:"size(50);column(country)" json:"country"`
	Province       string `orm:"size(50);column(province)" json:"province"`
	City           string `orm:"size(50);column(city)" json:"city"`
	Sex            bool   `orm:"column(sex)" json:"sex"`
	Signature      string `orm:"size(50);column(signature)" json:"signature"`
	Alias          string `orm:"size(50);column(alias_name)" json:"alias_name"`
	// Labels         []*Label `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Contact))
}
