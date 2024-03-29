package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/constant"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Group ...
type Group struct {
	ID             int64      `orm:"auto;column(id)"`
	GID            string     `orm:"size(50);column(g_id)" json:"wx_id" `                                // json:wx_id
	NickName       string     `orm:"size(50);column(nick_name)" json:"nick_name"`                        //
	Owner          string     `orm:"size(50);column(owner)" json:"owner" `                               //群主
	MemberNum      int        `orm:"column(member_num)" json:"member_num"`                               //
	HeadSmallImage string     `orm:"size(200);column(head_small_image_url)" json:"head_small_image_url"` //
	Listers        string     `orm:"size(500);column(listers)" json:"listers"`                           //成员微信号的IDs，”，“连接 接口返回值[]不好处理 记录1
	IsNeedServe    bool       `orm:"column(isneedserve);default(0)" json:"isneedserver"`                 // 是否有服务功能
	Bots           *Bots      `orm:"rel(fk)"`                                                            //
	GroupPlan      *GroupPlan `orm:"null;rel(fk);on_delete(set_null)"`                                   //群方案
	Messages       []*Message `orm:"reverse(many)"`                                                      //
}

// TableUnique ...
// 多字段唯一键
func (u *Group) TableUnique() [][]string {
	return [][]string{
		{"GID", "Bots"},
	}
}

func init() {
	orm.RegisterModel(new(Group))
}

// ProtoGiveGroup ...
func ProtoGiveGroup(Authorization string) (ret interface{}, err error) {
	var roomContactSeq = "0"
	var wxContactSeq = "0"

	type GetMultiDetailGroupInfo struct {
		Code int                       `json:"code"`
		Data []*common.DetailGroupInfo `json:"data"`
	}
	var retGroup []common.DetailGroupInfo
	for {
		var resp *http.Response
		if resp, err = httplib.Get(constant.CONTACT_GROUP_LIST_URL).Header(constant.H_AUTHORIZATION, Authorization).
			Param("room_contact_seq", roomContactSeq).Param("wx_contact_seq", wxContactSeq).DoRequest(); err != nil {
			logs.Error("get response[%s] failed, err is ", constant.CONTACT_GROUP_LIST_URL, err.Error())
			return nil, err
		}
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			logs.Error("get URL[%s] body failed, err is ", constant.CONTACT_GROUP_LIST_URL, err.Error())
			return nil, err
		}
		var restBody common.RecieveGroupList
		if err = json.Unmarshal(body, &restBody); err != nil {
			logs.Error("json Unmarshal failed, err is ", err.Error())
			return nil, err
		}
		/*
			1. 判断GroupList是否为空 空跳出
			2. 非空， 赋值参数，继续请求
			3. 返回所有群详情列表
		*/
		if len(restBody.Data.IDs) == 0 {
			// 群列表获取完全，退出返回
			break
		}
		// 拿到群ID去查看群详情(批量获取)
		var queryIDList []string
		for _, id := range restBody.Data.IDs {
			queryIDList = append(queryIDList, id)
		}
		query := "?group_id=" + strings.Join(queryIDList, "&group_id=")
		var gresp *http.Response
		if gresp, err = httplib.Get(constant.GROUP_INFO_URL+query).Header(constant.H_AUTHORIZATION, Authorization).DoRequest(); err != nil {
			logs.Error("get response[%s] failed, err is ", constant.GROUP_INFO_URL, err.Error())
			return nil, err
		}
		var allgbody []byte
		if allgbody, err = ioutil.ReadAll(gresp.Body); err != nil {
			logs.Error("get URL[%s] body failed, err is ", constant.GROUP_INFO_URL, err.Error())
			return nil, err
		}
		var allgrouprestBody GetMultiDetailGroupInfo
		if err = json.Unmarshal(allgbody, &allgrouprestBody); err != nil {
			logs.Error("json Unmarshal failed, err is ", err.Error())
			return nil, err
		}
		for _, _v := range allgrouprestBody.Data {
			retGroup = append(retGroup, *_v)
		}
		roomContactSeq = strconv.Itoa(restBody.Data.CurrentChatRoomContactSeq) //
		wxContactSeq = strconv.Itoa(restBody.Data.CurrentWxContactSeq)         //
	}
	return retGroup, nil
}

// GetGroupInfoByGIDs ... eg: gids 20102797431@chatroom,22475302355@chatroom
func GetGroupInfoByGIDs(gids string) (v interface{}, err error) {
	type GroupInfo struct {
		GId               string `json:"wx_id"`
		NickName          string `json:"nick_name"`
		HeadSmallImageUrl string `json:"head_small_image_url"`
	}
	if len(gids) == 0 {
		err := fmt.Errorf("GetGroupInfoByGIDs gids cant be null")
		return nil, err
	}
	var groups []*GroupInfo
	o := orm.NewOrm()
	if _, err = o.Raw("select DISTINCT(g_id), nick_name, head_small_image_url from `group` where ? like concat('%,', g_id)"+
		"or ? like concat('%,', g_id, ',%')"+"or ? like concat(g_id, ',%')"+"or ? = g_id;", gids, gids, gids, gids).QueryRows(&groups); err != nil {
		logs.Error("GetGroupInfoByGIDs failed, err is ", err.Error())
		return nil, err
	}
	return groups, nil
}
