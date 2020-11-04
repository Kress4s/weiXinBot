package models

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/constant"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
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

// SendTextMessage receive struct
// type SendTextMessage struct {
// 	Type    int         // 0: text 1:imgae 2:video 3:card  4: emoji
// 	From    string      //发送者
// 	To      string      // 接受者
// 	Content interface{} // data
// }
type SendTextMessage struct {
	At      []string `json:"at"`      // 发送者
	To      string   `json:"to"`      // 接受者
	Content string   `json:"content"` // word content
}

// SendImageMessage 发送图片消息
type SendImageMessage struct {
	To    string // 接受者
	URL   string //image url
	Token string //发送者的token
}

// AnnounceMessage 发送公告
type AnnounceMessage struct {
	Announcement string `json:"announcement"`
	GroupID      string `json:"group_id"`
}

// SendText ...
// @Param At: send wxid(origin WXID or modified wxid); To: received wxid ; content: send content
// @Param proto method [post]
func SendText(At, To, Content string) (err error) {
	var bot *Bots
	if bot, err = GetBotByWXID(At); err != nil {
		return
	}
	// json 发送的数据
	var sendText = new(SendTextMessage)
	sendText.At = append(sendText.At, At)
	sendText.To = To
	sendText.Content = Content
	res, verr := httplib.Post(constant.SEND_TEXT).Header(constant.H_AUTHORIZATION, bot.Token).JSONBody(&sendText)
	if verr != nil {
		logs.Error("[%+v] send message to [%s] faield, err is %s", sendText.At, sendText.To, err.Error())
		return verr
	}
	var response common.StandardRestResult
	if err = res.ToJSON(&response); err != nil {
		logs.Error("ToJSON: send text interface response failed, err is", err.Error())
	}
	// 目前地底层协议发送成功和失败code都是0，没明确提示
	if response.Code != 0 {
		logs.Error("sender[%s] send message[%s] to receiver[%s] failed, err is ", At, Content, To, err.Error())
		return err
	}
	return
}

// SendImage ...
// @Param At: send wxid(origin WXID or modified wxid); content: send content
// @Param proto method [get]
func SendImage(At, To, Content string) (err error) {
	var bot *Bots
	if bot, err = GetBotByWXID(At); err != nil {
		return
	}
	resp, verr := httplib.Get(constant.SEND_IMAGE).Header(constant.H_AUTHORIZATION, bot.Token).Param("to", To).Param("url", Content).DoRequest()
	if verr != nil {
		logs.Error("send image[%s] to[%s] failed, err is ", Content, To)
		return verr
	}
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	var response common.StandardRestResult
	if err = json.Unmarshal(body, &response); err != nil {
		logs.Error("SendImage: json Unmarshal failed, err is ", err.Error())
		return
	}
	// 目前地底层发送成功和失败code都是0，没明确提示
	if response.Code != 0 {
		logs.Error("send message[%s] to receiver[%s] failed, err is ", Content, To, err.Error())
		return err
	}
	return
}

// SendAnnouncement ...
// @Param At: send wxid(origin WXID or modified wxid); content: send content
// @Param proto method [get]
func SendAnnouncement(At, To, Content string) (err error) {
	var bot *Bots
	if bot, err = GetBotByWXID(At); err != nil {
		return
	}
	announce := new(AnnounceMessage)
	announce.Announcement = Content
	announce.GroupID = To
	res, verr := httplib.Post(constant.SEND_TEXT).Header(constant.H_AUTHORIZATION, bot.Token).JSONBody(&announce)
	if verr != nil {
		logs.Error("[%+v] send message to [%s] faield, err is %s", bot.WXID, announce.GroupID, err.Error())
		return verr
	}
	var response common.StandardRestResult
	if err = res.ToJSON(&response); err != nil {
		logs.Error("ToJSON: send text interface response failed, err is", err.Error())
	}
	// 目前地底层发送成功和失败code都是0，没明确提示
	if response.Code != 0 {
		logs.Error("send message[%s] to receiver[%s] failed, err is ", Content, To, err.Error())
		return err
	}
	return
}
