package models

// AutoKick ...
type AutoKick struct {
	ID        int64      `orm:"auto;column(id)"`           //
	Type      int        `orm:"column(type);default(3)"`   // 所属功能类型 （默认3）
	Switch    bool       `orm:"column(switch);default(1)"` //功能总开关
	GroupPlan *GroupPlan `orm:"rel(fk)"`                   //
}

// KickOnly 仅踢人
type KickOnly struct {
	ID                    int64               `orm:"auto;column(id)"`                        //
	MessageKeyWordsSwitch bool                `orm:"column(msg_key_word_switch);default(1)"` //触发消息关键词开关
	KickOtherRule         *KickOtherRule      `orm:"null;rel(one)"`                          //其他规则
	MessageKeyWords       []*MessageKeyWords  `orm:"reverse(many)"`
	NickNameKeyWords      []*NickNameKeyWords `orm:"reverse(many)"`
}

// MessageKeyWords 触发关键字表
type MessageKeyWords struct {
	ID       int64     `orm:"auto;column(id)"`           //
	KWord    string    `orm:"size(20);column(key_word)"` // 触发的关键词
	Switch   bool      `orm:"column(switch);default(0)"` //开关
	KickOnly *KickOnly `orm:"rel(fk)"`
}

// NickNameKeyWords 昵称关键字表
type NickNameKeyWords struct {
	ID       int64     `orm:"auto;column(id)"`           //
	KWord    string    `orm:"size(20);column(key_word)"` // 触发的关键词
	Switch   bool      `orm:"column(switch);default(0)"` //开关
	KickOnly *KickOnly `orm:"rel(fk)"`
}

// KickOtherRule 仅踢人的其他规则
type KickOtherRule struct {
	ID                 int64     `orm:"auto;column(id)"`                         //
	Types              string    `orm:"size(50);column(types)"`                  //已勾选的类型ids(二维码/名片)
	WordLimit          int       `orm:"column(word_limit)"`                      //字数上限
	WordLimitSwitch    bool      `orm:"column(word_limit_switch);default(0)"`    // 字数上限开关
	ContinueSendSwitch bool      `orm:"column(continue_send_switch);default(0)"` //连续发送踢人开关
	ContinueInfo       string    `orm:"size(50);column(continue_info)"`          //10,20  10秒发送20条消息
	KickOnly           *KickOnly `orm:"reverse(one)"`                            //
}
