package models

import (
	"fmt"
	"strings"
	"weiXinBot/app/bridage/common"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Configuration 功能设置(消息监听使用)
type Configuration struct {
	ID         int64  `orm:"auto;column(id)"`               //
	Type       int    `orm:"column(type)"`                  // 配置对象 0: 群组; 1: 联系人 ...可拓展
	FuncType   int    `orm:"column(function_type)"`         // 功能配置类型 1:入群欢迎语 2:关键词回复 3:自动踢人...可拓展
	FuncID     int64  `orm:"column(function_id)"`           // 配置ID
	BotWXID    string `orm:"size(30);column(bot_wxid)"`     // 机器人微信ID，执行消息回复、踢人等操作的微信号(保证机器人是正确的)
	ObjectIDS  string `orm:"size(1000);column(object_ids)"` // 要执行对象的IDs,多个用”,“连接(群、联系人...可拓展)
	GrouplanID int64  `orm:"column(grouplan_id)"`           // 所属的方案
}

// MultiDealConfig ...
type MultiDealConfig struct {
	Type         int
	FuncInfoList []struct {
		BotWXID    string
		Info       map[string]int64
		ObjectsIDS string
		GrouplanID int64
	}
}

func init() {
	orm.RegisterModel(new(Configuration))
}

// GroupIsNeedServer 查看此群是否需要机器人服务
// @Param fromUserName: {Str:22475302355@chatroom}
// @Param toUserName: {Str:wxid_vao3ptfez4p22}}
// func GroupIsNeedServer(fromUserName, toUserName string) (isServer bool, err error) {
func GroupIsNeedServer(message common.ProtoMessage) (isServer bool, err error) {
	o := orm.NewOrm()
	/*
		这里可以做我的群管状态修改(gid和botid判断)
	*/
	// is group message
	if !strings.Contains(message.FromUserName.Str, "@chatroom") {
		return false, nil
	}
	if !o.QueryTable(new(Configuration)).Filter("Type", 0).Filter("ObjectIDS__icontains", message.FromUserName.Str).Filter("BotWXID", message.ToUserName.Str).Exist() {
		return false, nil
	}
	return true, nil
}

// GroupService 需要的服务(master method)
// message.PushContent: "push_content":"🛫张 : G吐总冠军"
// func GroupService(fromUserName, toUserName, keyContent string) {
func GroupService(message common.ProtoMessage) {
	o := orm.NewOrm()
	var err error
	var configs []*Configuration
	if _, err = o.QueryTable(new(Configuration)).Filter("Type", 0).Filter("ObjectIDS__icontains", message.FromUserName.Str).Filter("BotWXID", message.ToUserName.Str).All(&configs); err != nil {
		logs.Error("get Configuration by fromUserName and toUserName failed, err is ", err.Error())
	}
	for _, v := range configs {
		switch v.FuncType {
		// is welcome function config
		case 1:
			//确定新人进群的数据结构再做处理
			var parsesysmsg *common.WxSysMsg
			if parsesysmsg, err = common.PraseXMLString(message.Content.Str); err != nil {
				logs.Error(err.Error())
				continue
			}
			// 当前文本消息的类型是 MesType = 10002 现在接入语音也是
			if parsesysmsg == nil {
				continue
			}
			if message.MsgType == 10002 && parsesysmsg.Type != "sysmsgtemplate" {
				continue
			}
			if strings.Contains(parsesysmsg.SysmsgTemplate.ContenTemplate.Template, "kickoutname") {
				continue
			}
			var replyResource []*Resource
			var isNeedServer bool
			if isNeedServer, replyResource, err = WelcomeService(v.FuncID, message.PushContent); err != nil {
				logs.Error("%s send some resources to group or contact([%s]) failed...", v.BotWXID, message.FromUserName.Str)
				continue
			} else if isNeedServer && replyResource != nil && err == nil {
				for _, _rR := range replyResource {
					for _, _rM := range _rR.Material {
						switch _rM.Type {
						case 1:
							//回复的文字内容
							logs.Info("开始发送文字...")
							if strings.Contains(_rM.Data, "{{@新人}}") {
								var newAtData = fmt.Sprintf("@%s", parsesysmsg.SysmsgTemplate.ContenTemplate.Linklist.Link.MemberList.Member[0].NickName)
								_rM.Data = strings.ReplaceAll(_rM.Data, "{{@新人}}", newAtData)
							}
							if err = SendText(message.ToUserName.Str, message.FromUserName.Str, _rM.Data); err != nil {
								logs.Error("SendText %s send %s to %s failed, err is ", message.ToUserName.Str, message.FromUserName.Str, _rM.Data)
							}
						case 2:
							// 图片内容
							logs.Info("开始发送图片...")
							if err = SendImage(message.ToUserName.Str, message.FromUserName.Str, _rM.Data); err != nil {
								logs.Error("SendImage %s send %s to %s failed, err is ", message.ToUserName.Str, message.FromUserName.Str, _rM.Data)
							}
						default:
							fmt.Println("等待扩展的类型")
						}
					}
				}
			}
		// is keywords function config
		case 2:
			// 当前文本消息的类型是 MesType = 1
			// nameContent := strings.SplitN(keyContent, ":", 2)
			if message.MsgType != 1 {
				continue
			}
			fmt.Println("开启关键词查询服务..")
			var replyResource []*Resource
			var isNeedServer bool
			if isNeedServer, replyResource, err = KeyWordsService(v.FuncID, message.Content.Str); err != nil {
				// if isNeedServer, replyResource, err = KeyWordsService(v.FuncID, message.PushContent); err != nil {
				logs.Error("KeyWordsService failed, err is ", err.Error())
				continue
			} else if isNeedServer && err == nil {
				logs.Info("找到问题库...开始查找资源...")
				for _, _rR := range replyResource {
					for _, _rM := range _rR.Material {
						switch _rM.Type {
						case 1:
							//回复的文字内容
							logs.Info("开始发送文字...")
							if err = SendText(message.ToUserName.Str, message.FromUserName.Str, _rM.Data); err != nil {
								logs.Error("SendText %s send %s to %s failed, err is ", message.ToUserName.Str, message.FromUserName.Str, _rM.Data)
							}
						case 2:
							// 图片内容
							logs.Info("开始发送图片...")
							if err = SendImage(message.ToUserName.Str, message.FromUserName.Str, _rM.Data); err != nil {
								logs.Error("SendImage %s send %s to %s failed, err is ", message.ToUserName.Str, message.FromUserName.Str, _rM.Data)
							}
						default:
							fmt.Println("等待扩展的类型")
						}
					}
				}
			}
			// fmt.Println("关键词回复")
		// is autokick function config
		case 3:
			fmt.Println("自动踢人")
		// is white function config
		case 4:
			fmt.Println("白名单")
		default:
			logs.Error("function config Type[%d] is not right, please cheak it and modify it", v.FuncType)
		}
	}
}

// DeleteConifgForWxMigration ...
func DeleteConifgForWxMigration(wxidOrPlanid interface{}) (err error) {
	o := orm.NewOrm()
	var num int64
	// 账号迁移删除的配置
	if wxid, ok := wxidOrPlanid.(string); ok {
		if num, err = o.QueryTable(new(Configuration)).Filter("BotWXID", wxid).Delete(); err == nil {
			logs.Debug("Number of Configuration deleted in database:", num)
			return
		}
	} else if grouplanid, ok := wxidOrPlanid.(int64); ok {
		// 删除方案批量清空配置
		if num, err = o.QueryTable(new(Configuration)).Filter("GrouplanID", grouplanid).Delete(); err == nil {
			logs.Debug("Number of Configuration deleted in database:", num)
		}
	}
	return
}

// UpdateConfigByCutPlan 我的群管修改方案 或者 设置微信群到另外一个方案下的情况， 需要针对性的修改config配置
func UpdateConfigByCutPlan(BotWXID, gid string, grouplanID int64) (err error) {
	var configs []*Configuration
	o := orm.NewOrm()
	if _, err = o.QueryTable(new(Configuration)).Filter("BotWXID", BotWXID).Filter("GrouplanID", grouplanID).
		Filter("ObjectIDS__contains", gid).All(&configs); err == nil {
		// 因为一个方案下的一个微信号的配置的群都是一样的,去掉被转移方案的群号,取一个值批量更新
		var _objectids string
		if strings.Contains(configs[0].ObjectIDS, gid+",") {
			_objectids = strings.ReplaceAll(configs[0].ObjectIDS, gid+",", "")
		} else if strings.Contains(configs[0].ObjectIDS, ","+gid) {
			_objectids = strings.ReplaceAll(configs[0].ObjectIDS, ","+gid, "")
		} else if configs[0].ObjectIDS == gid {
			_objectids = strings.ReplaceAll(configs[0].ObjectIDS, gid, "")
		} else {
			logs.Error("UpdateConfigByCutPlan: config objectsids is null, bot[%s], grouplanID[%d]", BotWXID, grouplanID)
			return
		}
		// 防止一个方案下所有的群都被替换了，为空了 此时应该把这个方案下所有的配置都山删掉
		var num int64
		if num, err = o.QueryTable(new(Configuration)).Filter("BotWXID", BotWXID).Filter("GrouplanID", grouplanID).
			Filter("ObjectIDS__contains", gid).Update(orm.Params{
			"ObjectIDS": _objectids,
		}); err == nil {
			logs.Debug("Number of Configuration deleted in database:", num)
		}
	}
	return
}
