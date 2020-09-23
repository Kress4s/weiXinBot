package models

// WhiteList ...
type WhiteList struct {
	WXID string `orm:"size(50);column(wx_id)"` // 微信id
}
