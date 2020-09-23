package models

// AutoKick ...
type AutoKick struct {
	ID     int64 `orm:"auto;column(id)"`           //
	Switch bool  `orm:"column(switch);default(1)"` //功能总开关
}

// KickOnly 仅踢人
type KickOnly struct {
	ID                    int64              `orm:"auto;column(id)"`                        //
	MessageKeyWordsSwitch bool               `orm:"column(msg_key_word_switch);default(1)"` //触发消息关键词开关
	MessageKeyWords       []*MessageKeyWords `orm:"reverse(many)"`
}

// MessageKeyWords 触发关键字表
type MessageKeyWords struct {
	ID       int64     `orm:"auto;column(id)"`               //
	Word     string    `orm:"size(20);column(msg_key_word)"` // 触发的关键词
	KickOnly *KickOnly `orm:"rel(fk)"`
}
