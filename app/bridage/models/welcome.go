package models

// Welcome ...
type Welcome struct {
	ID          int64  `orm:"auto;column(id)"` //
	Words       string `orm:"size(200);column(words)"`
	WaitSeconds int    `orm:"column(waitseconds)"`
	StartTime   string `orm:"size(20);column(startime)"`
	EndTime     string `orm:"size(20);column(endtime)"`
	NickName    string `orm:"size(50);column(nick_name)" json:"nick_name"`
}
