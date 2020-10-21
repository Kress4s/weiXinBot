package common

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
	ID                string `json:"id"`
	NickName          string `json:"nick_name"`
	AliasName         string `json:"alias_name"`
	Sex               int    `json:"sex"`
	Country           string `json:"country"`
	Province          string `json:"province"`
	City              string `json:"city"`
	Signature         string `json:"signature"`
	HeadBigImageURL   string `json:"head_big_image_url"`
	HeadSmallImageURL string `json:"head_small_image_url"`
	Label             string `json:"label"`
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
