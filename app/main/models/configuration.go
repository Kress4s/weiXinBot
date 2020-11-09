package models

import (
	"strconv"
	"strings"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddConfiguration ...
func AddConfiguration(m *bridageModels.Configuration) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// UpdateConfigurationByID ...
func UpdateConfigurationByID(m *bridageModels.Configuration) (err error) {
	var v = bridageModels.Configuration{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "ObjectIDS"); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}

// UpdateOrAddConfig ...
func UpdateOrAddConfig(m bridageModels.MultiDealConfig) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	for k, v := range m.FuncInfoList {
		var config bridageModels.Configuration
		config.Type = m.Type
		config.BotWXID = v.BotWXID
		config.ObjectIDS = v.ObjectsIDS
		config.GrouplanID = v.GrouplanID
		for _k, _v := range m.FuncInfoList[k].Info {
			config.FuncType, _ = strconv.Atoi(_k)
			config.FuncID = _v
			_config := bridageModels.Configuration{Type: config.Type, FuncID: config.FuncID, FuncType: config.FuncType, BotWXID: config.BotWXID}
			if err = o.Read(&_config, "Type", "FuncType", "FuncID", "BotWXID"); err != nil {
				if err == orm.ErrNoRows {
					_config.ObjectIDS = config.ObjectIDS
					if _, err = o.Insert(&_config); err != nil {
						return
					}
					err = nil
				}
			} else {
				var num int64
				config.ID = _config.ID
				if num, err = o.Update(&config, "ObjectIDS"); err != nil {
					logs.Error("update config failed, err is ", err.Error())
					return
				}
				logs.Debug("Number of Bot update in database:", num)
			}
		}
	}
	return
}

// GetConfigRelation ...
func GetConfigRelation(grouplanID int64) (ret interface{}, err error) {
	o := orm.NewOrm()
	var l []bridageModels.GBGRelation
	if _, err = o.QueryTable(new(bridageModels.GBGRelation)).Filter("GrouplanID", grouplanID).All(&l); err != nil {
		logs.Error("Get ConfigRelation failed, err is ", err.Error())
		return nil, err
	}
	return l, nil
}

// UpdateConfigRelation ...
func UpdateConfigRelation(m []bridageModels.GBGRelation, WXID string, grouplanID int64) (err error) {
	o := orm.NewOrm()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	o.Begin()
	if len(WXID) != 0 {
		WXIDSlice := strings.Split(WXID, ",")
		if _, err = o.QueryTable(new(bridageModels.GBGRelation)).Filter("GrouplanID", grouplanID).Filter("BotWXID__in", WXIDSlice).Delete(); err != nil {
			logs.Error("delete config  group wxid failed, err is ", err.Error())
			return
		}
	}
	for _, v := range m {
		var _v = bridageModels.GBGRelation{GrouplanID: v.GrouplanID, BotWXID: v.BotWXID}
		if err = o.Read(&_v, "GrouplanID", "BotWXID"); err != nil {
			if err == orm.ErrNoRows {
				if _, err = o.Insert(&v); err != nil {
					logs.Error("GetConfigRelation: insert GrouplanWXRelation failed, err is ", err.Error())
					return
				}
			}
		} else {
			v.ID = _v.ID
			if _, err = o.Update(&v, "ObjectIDS"); err != nil {
				logs.Error("GetConfigRelation: update GrouplanWXRelation failed, err is ", err.Error())
				return
			}
		}
	}
	return
}
