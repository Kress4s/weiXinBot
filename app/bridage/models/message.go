package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// Message ...
type Message struct {
	ID          int64     `orm:"auto;column(id)"`
	From        string    `orm:"size(50);column(from)" json:"at"`               //发送者
	To          string    `orm:"size(50);column(to)" json:"to"`                 //接受者
	BelongTo    int       `orm:"column(belong_type)" json:"belong_type"`        //消息来源类型 0:好友 1:群 2:公众号 3：系统消息 4...
	MessageType int       `orm:"column(msg_type)" json:"msg_type"`              //消息内容类型 0: text 1:imgae 2:video 3:card  4: emoji
	IsMe        bool      `orm:"column(isme)" json:"isme"`                      //是否我发送的内容
	SendTime    time.Time `orm:"auto_now_add;type(datetime)" json:"sendtime"`   //发送时间
	Text        string    `orm:"type(text)"`                                    //消息文本
	ImageURL    string    `orm:"size(300); column(image_url)" json:"image_url"` //图片地址
	VideoURL    string    `orm:"size(300); column(video_url)" json:"video_url"` //视频地址
	Card        *Card     `orm:"rel(fk);null"`                                  //名片信息
	Contact     *Contact  `orm:"rel(fk);null"`                                  //联系人信息
	Group       *Group    `orm:"rel(fk);null"`                                  //群信息
	WxID        string    `orm:"size(50);column(wx_id)"`                        //所属微信号消息
	Emoji       string    `orm:"size(50); column(emoji)" json:"emoji"`          //emoji的md5值
}

// Card ...
type Card struct {
	ID           int64  `orm:"auto;column(id)"`
	CardAlias    string `orm:"size(50);column(card_alias)" json:"card_alias"`
	CardID       string `orm:"size(50);column(card_id)" json:"card_id"`
	CardNickName string `orm:"size(50);column(card_nick_name)" json:"card_nick_name"`
}

func init() {
	orm.RegisterModel(new(Message), new(Card))
}

// SendMessage receive struct
type SendMessage struct {
	Type    int         // 0: text 1:imgae 2:video 3:card  4: emoji
	Content interface{} // data
}
