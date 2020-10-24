package common

import "encoding/xml"

//RestResult Rest接口返回值
type RestResult struct {
	Code    int         // 0 表示成功，其他失败
	Message string      // 错误信息
	Data    interface{} // 数据体
}

// StandardRestResult 标准的 rest 返回接口，字符小写化
type StandardRestResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// RecieveGroupList ...
// 从协议获取群ID列表数据结构
type RecieveGroupList struct {
	Code int `json:"code"`
	Data struct {
		CurrentWxContactSeq       int      `json:"current_wx_contact_seq"`
		CurrentChatRoomContactSeq int      `json:"current_chat_room_contact_seq"`
		IDs                       []string `json:"ids"`
	} `json:"data"`
}

// DetailGroupInfo ...
// 从协议获取群详细信息的数据结构
type DetailGroupInfo struct {
	ID        string `json:"wx_id"`
	NickName  string `json:"nick_name"`
	Owner     string `json:"owner"`
	MemberNum int    `json:"member_num"`
	// AliasName         string `json:"alias_name"`
	// Sex               int    `json:"sex"`
	// Country           string `json:"country"`
	// Province          string `json:"province"`
	// City              string `json:"city"`
	// Signature         string `json:"signature"`
	// HeadBigImageURL   string `json:"head_big_image_url"`
	HeadSmallImageURL string `json:"head_small_image_url"`
	// Status            int    `json:"status"` (群组的属性 不知道干嘛的先保留)
	// Label             string `json:"label"`
}

// WxSysMsg 目前踢人和新人入群发送的消息通知
type WxSysMsg struct {
	XMLName        xml.Name `xml:"sysmsg"`
	SysmsgTemplate struct {
		ContenTemplate struct {
			// Plain    string `xml:"plain"`
			Template string `xml:"template"`
			Linklist struct {
				Link []struct {
					Type       string `xml:"name,attr"` //新人入群：adder；踢人：kickoutname
					MemberList struct {
						Member []struct {
							Username string `xml:"username"`
							NickName string `xml:"nickname"`
						} `xml:"member"`
					} `xml:"memberlist"`
				} `xml:"link"`
			} `xml:"link_list"`
		} `xml:"content_template"`
	} `xml:"sysmsgtemplate"`
}

// ProtoMessage 底层协议推送的消息
type ProtoMessage struct {
	FromUserName struct {
		Str string `json:"str"`
	} `json:"from_user_name"` //
	ToUserName struct {
		Str string `json:"str"`
	} `json:"to_user_name"` //
	MsgType int `json:"msg_type"` // 消息类型 10002(踢人、加人的消息类型(xml))
	Content struct {
		Str string `json:"str"`
	} `json:"content"` // 内容(我发：{"str":"程序监控你"}；别人发：{"str":"aaaa520jj:\nG吐总冠军"})
	Status      int    `json:"status"`       //貌似群的消息都是
	CreateTime  int    `json:"create_time"`  //消息时间戳
	MsgSource   string `json:"msg_source"`   // ?
	PushContent string `json:"push_content"` //提示消息(聊天输入框提示) (别人发有这个字段，我发没有这个字段)
}

// 查询条件常量
const (
	MultiSelect      = "multi-select"       // 多选
	MultiText        = "multi-text"         // 模糊多选
	NumRange         = "num-range"          // 数字范围或者数字查询
	CommaMultiSelect = "comma-multi-select" //数据库中存放是逗号分隔的值
)

// QueryConditon 表格查询条件对象
type QueryConditon struct {
	QueryKey    string
	QueryType   string // multi-select  multi-text  num-range  comma-multi-select
	QueryValues []string
}
