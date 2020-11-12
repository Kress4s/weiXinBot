package models

import (
	"strconv"
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
					if _, err = o.Insert(&config); err != nil {
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
