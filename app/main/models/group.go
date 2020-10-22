package models

import (
	"strconv"
	"strings"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddGroup ...
func AddGroup(m *bridageModels.Group) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

//MultiAddGroup ...
func MultiAddGroup(m []*bridageModels.Group) (err error) {
	o := orm.NewOrm()
	for _, _m := range m {
		_, err = o.Insert(_m)
	}
	return
}

// GetGroupByGID ...
func GetGroupByGID(GID string) (v *bridageModels.Group, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Group{GID: GID}
	if err = o.Read(v, "GID"); err != nil {
		return nil, err
	}
	return v, nil
}

// UpdateGrouByID ...
func UpdateGrouByID(m *bridageModels.Group) (err error) {
	var v = bridageModels.Group{GID: m.GID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, "IsNeedServe", "GroupPlan"); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}

// MultiUpdateGrouByID ...
func MultiUpdateGrouByID(m []*bridageModels.Group, delgroupsIDs string) (err error) {
	/*
		思路:
		1. 根据机器人和所属GID去查询(存在表明这个机器人管理的群有配置)
		2. 有配置的情况再去更新群管方案()
		3. 没有配置直接插入
		4. 把移除该分组的群GroupPlan都置为null(RelatedSel()方法会有问题，目前没用到不考虑)
	*/
	o := orm.NewOrm()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	if err = o.Begin(); err != nil {
		return
	}
	for _, _m := range m {
		if !o.QueryTable(new(bridageModels.Group)).Filter("GID", _m.GID).Filter("Bots", _m.Bots).Exist() {
			// 未存在
			if _, err = o.Insert(_m); err != nil {
				logs.Error("MultiUpdateGrouByID: insert Group failed, err is ", err.Error())
				return
			}
		} else {
			// 已存在
			var _M = bridageModels.Group{GID: _m.GID, Bots: _m.Bots}
			if err = o.Read(&_M, "GID", "Bots"); err == nil {
				var num int64
				if num, err = o.Update(_m); err == nil {
					logs.Debug("Number of Group update in database:", num)
				}
			}
		}
	}
	// 解析删除群组的参数gid:botid,gid2:botid2...
	delGroupSlice := strings.Split(delgroupsIDs, ",")
	for _, m := range delGroupSlice {
		gbid := strings.Split(m, ":")
		gid := gbid[0]
		bid, verr := strconv.ParseInt(gbid[1], 0, 64)
		if verr != nil {
			err = verr
			return
		}
		if _, err = o.QueryTable(new(bridageModels.Group)).Filter("GID", gid).Filter("Bots__ID", bid).Update(orm.Params{
			"GroupPlan": nil,
		}); err != nil {
			return
		}
	}
	return
}
