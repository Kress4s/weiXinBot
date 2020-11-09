package models

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"weiXinBot/app/bridage/common"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddGrouplan ...
func AddGrouplan(grouplan *bridageModels.GroupPlan) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(grouplan)
	return
}

// GetGouplanByID ...
func GetGouplanByID(id int64) (v *bridageModels.GroupPlan, err error) {
	o := orm.NewOrm()
	v = &bridageModels.GroupPlan{ID: id}
	if err = o.QueryTable(new(bridageModels.GroupPlan)).Filter("ID", id).One(v); err != nil {
		return nil, err
	}
	o.LoadRelated(v, "Groups")
	for i := range v.Groups {
		o.QueryTable(new(bridageModels.Group)).Filter("ID", v.Groups[i].ID).RelatedSel("Bots").One(v.Groups[i])
	}
	return v, nil
}

// GetAllGrouplan get all bots
func GetAllGrouplan(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.GroupPlan))
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

	var l []bridageModels.GroupPlan
	qs = qs.OrderBy(sortFields...).RelatedSel()
	// add IsDelete filter
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				o.LoadRelated(&v, "Groups")
				for i := range v.Groups {
					o.QueryTable(new(bridageModels.Group)).Filter("ID", v.Groups[i].ID).RelatedSel("Bots").One(v.Groups[i])
				}
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

// UpdateGrouplanByID ...
func UpdateGrouplanByID(m *bridageModels.GroupPlan) (err error) {
	var v = bridageModels.GroupPlan{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}

// DeleteGrouplanByID deletes bot by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGrouplanByID(id int64) (err error) {
	o := orm.NewOrm()
	/*
	 1. 找出该方案下所有的群（当groupPlan被删除时自动置空）
	 2. 找到这些群被哪些机器人托管
	 3. 删除这些
	*/
	o.Begin()
	defer func() {
		if err == nil {
			o.Commit()
		} else {
			o.Rollback()
		}
	}()
	var groups []*bridageModels.Group
	if _, err = o.QueryTable(new(bridageModels.Group)).Filter("GroupPlan__ID", id).RelatedSel().All(&groups); err != nil {
		logs.Error("DeleteGrouplanByID: get all groups by planid failed, err is ", err.Error())
		return
	}
	var wxids []string
	for _, group := range groups {
		wxids = append(wxids, group.Bots.WXID)
	}
	if err = bridageModels.DeleteConifgForWxMigration(strings.Join(wxids, ",")); err != nil {
		return
	}
	v := bridageModels.GroupPlan{ID: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&bridageModels.GroupPlan{ID: id}); err == nil {
			logs.Debug("Number of Bots deleted in database:", num)
		}
	}
	return
}

// MultiDeleteGrouplanByIDs multi delete Bot
func MultiDeleteGrouplanByIDs(ids []interface{}) (err error) {
	o := orm.NewOrm()
	var num int64
	if num, err = o.QueryTable(new(bridageModels.GroupPlan)).Filter("ID__in", ids...).Delete(); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}
