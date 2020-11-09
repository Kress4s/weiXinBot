package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"weiXinBot/app/bridage/common"
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

// GetAllGroup get all bots
func GetAllGroup(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.Group))
	// query QueryCondition
	cond := orm.NewCondition()
	for _, query := range querys {
		var k string
		cond1 := orm.NewCondition()
		switch query.QueryType {
		case common.MultiSelect:
			k = query.QueryKey // + "__iexact"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		case common.MultiText:
			k = query.QueryKey + "__icontains"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		case common.NumRange:
			if len(query.QueryValues) == 2 {
				var from, to float64
				if from, err = strconv.ParseFloat(query.QueryValues[0], 64); err != nil {
					logs.Error(err.Error())
					return
				}
				if to, err = strconv.ParseFloat(query.QueryValues[1], 64); err != nil {
					logs.Error(err.Error())
					return
				}
				k = query.QueryKey + "__gte"
				cond1 = cond1.Or(k, from)
				k = query.QueryKey + "__lte"
				cond1 = cond1.And(k, to)
				cond = cond.AndCond(cond1)
			} else {
				k = query.QueryKey + "__icontains"
				for _, v := range query.QueryValues {
					cond1 = cond1.Or(k, v)
				}
				cond = cond.AndCond(cond1)
			}
		default:
			k = query.QueryKey + "__icontains"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		}
	}
	qs = qs.SetCond(cond)
	// get query's total count
	if totalcount, err = qs.Count(); err != nil {
		return
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, 0, errors.New("Error: unused 'order' fields")
		}
	}

	var l []bridageModels.Group
	qs = qs.OrderBy(sortFields...).RelatedSel()
	// add IsDelete filter
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, totalcount, nil
	}
	return nil, 0, err
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

// UpdateGroupByID ...
func UpdateGroupByID(m *bridageModels.Group) (err error) {
	/*
		思路:
		1. 修改群对应的方案(根据群号和微信号判断)
		2. 修改配置信息：找到机器人的配置信息和修改前后的方案id，且objectids中存在该群的配置，修改即可
	*/
	var v = bridageModels.Group{GID: m.GID}
	o := orm.NewOrm()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	o.Begin()
	if err = o.QueryTable(new(bridageModels.Group)).Filter("GID", m.ID).RelatedSel("Bots").One(&v); err == nil {
		var num int64
		if v.GroupPlan != m.GroupPlan {
			if err = bridageModels.UpdateConfigByCutPlan(v.Bots.WXID, v.GID, v.GroupPlan.ID); err != nil {
				return
			}
			// 新方案后面所有的配置objects要补上
			if o.QueryTable(new(bridageModels.Configuration)).Filter("BotWXID", m.Bots.WXID).Filter("GrouplanID", m.GroupPlan.ID).Exist() {
				// 该微信号之前有配置信息，直接加上群号
				if num, err = o.QueryTable(new(bridageModels.Configuration)).Filter("BotWXID", m.Bots.WXID).Filter("GrouplanID", m.GroupPlan.ID).
					Update(orm.Params{
						"ObjectIDS": orm.ColValue(orm.ColAdd, m.GID),
					}); err == nil {
					logs.Debug("Number of Config update in database:", num)
				}
			} else {
				// 该微信号的当前的群之前配置表中没有配置过，新增方案下面的配置信息，批量插入
				var newConfigs []*bridageModels.Configuration
				if _, err = o.QueryTable(new(bridageModels.Configuration)).Filter("GrouplanID", m.GroupPlan.ID).All(&newConfigs); err == nil {
					for _, _v := range newConfigs {
						_v.ID = 0
						_v.BotWXID = m.Bots.WXID
						_v.ObjectIDS = m.GID
					}
					if num, err = o.InsertMulti(1, newConfigs); err == nil {
						logs.Error("UpdateGroupByID: InsertMulti")
					}
				}
			}
		}
		if num, err = o.Update(&m, "IsNeedServe", "GroupPlan"); err == nil {
			logs.Debug("Number of Group update in database:", num)
		}
	}
	return
}

// MultiUpdateGroupByID ...
func MultiUpdateGroupByID(m []*bridageModels.Group, delgroupsIDs string) (err error) {
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
		if _m.NickName == "" {
			_m.NickName = fmt.Sprintf("群聊(%d)", _m.MemberNum)
		}
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
				if _M.GroupPlan != nil && _M.GroupPlan.ID != _m.GroupPlan.ID {
					// 该群之前有配置且当前配置切换了方案
					// 修改配置表的objectids
					if err = bridageModels.UpdateConfigByCutPlan(_M.Bots.WXID, _M.GID, _M.GroupPlan.ID); err != nil {
						return
					}
				}
				var num int64
				_m.ID = _M.ID
				if num, err = o.Update(_m); err == nil {
					logs.Debug("Number of Group update in database:", num)
				}
			}
		}
	}
	// 解析删除群组的参数gid:botid,gid2:botid2...
	if len(delgroupsIDs) == 0 {
		return nil
	}
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

// DeleteGroupMgrByID 我的群管删除
func DeleteGroupMgrByID(gid string) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	group := bridageModels.Group{GID: gid}
	if err = o.QueryTable(new(bridageModels.Group)).Filter("GID", gid).RelatedSel("Bots").One(&group); err != nil {
		logs.Error("DeleteGroupMgrByID: get group failed, err is", err.Error())
		return
	}
	if err = bridageModels.UpdateConfigByCutPlan(group.Bots.WXID, group.GID, group.GroupPlan.ID); err != nil {
		return
	}
	var num int64
	group.GroupPlan = nil
	if num, err = o.Update(&group, "GroupPlan"); err == nil {
		logs.Debug("Number of Group update in database:", num)
	}
	return
}
