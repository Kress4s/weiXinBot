package models

import (
	"errors"
	"reflect"
	"strconv"
	"weiXinBot/app/bridage/common"
	taskflow "weiXinBot/app/bridage/flows"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddTimeTask ...
func AddTimeTask(m *bridageModels.TimeTask) (id int64, err error) {
	var timetasker taskflow.TaskFactory
	o := orm.NewOrm()
	if timetasker, err = taskflow.GetTaskIns(m.Type); err != nil {
		logs.Error(err.Error())
		return 0, err
	}
	if id, err = o.Insert(m); err != nil {
		logs.Error("insert task failed, err is ", err.Error())
	}
	switch m.SendType {
	case -1:
		go timetasker.TaskImmediately(*m)
	default:
		go timetasker.TaskSetting(*m)
	}
	return
}

// GetTimeTaskByID ...
func GetTimeTaskByID(id int64) (v *bridageModels.TimeTask, err error) {
	o := orm.NewOrm()
	v = &bridageModels.TimeTask{ID: id}
	if err = o.Read(v); err != nil {
		return nil, err
	}
	return v, nil
}

// GetAllTimeTask get all bots
func GetAllTimeTask(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.TimeTask))
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

	var l []bridageModels.TimeTask
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

// UpdateTimeTaskByID ...
func UpdateTimeTaskByID(m *bridageModels.TimeTask) (err error) {
	var v = bridageModels.TimeTask{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of Bot update in database:", num)
			var timetasker taskflow.TaskFactory
			if timetasker, err = taskflow.GetTaskIns(m.Type); err != nil {
				logs.Error(err.Error())
				return err
			}
			// 修改定时任务
			go timetasker.ModifyTimeTask(*m)
		}
	}
	return
}

// DeleteTimeTaskByID deletes Material by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTimeTaskByID(id int64) (err error) {
	o := orm.NewOrm()
	v := bridageModels.TimeTask{ID: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&bridageModels.TimeTask{ID: id}); err == nil {
			logs.Debug("Number of Bots deleted in database:", num)
			var timetasker taskflow.TaskFactory
			if timetasker, err = taskflow.GetTaskIns(v.Type); err != nil {
				logs.Error(err.Error())
				return err
			}
			// 删除定时任务
			go timetasker.ModifyTimeTask(v)
		}
	}
	return
}

// GetAllRecode ...
func GetAllRecode(manager, status string) (ret interface{}, count int64, err error) {
	o := orm.NewOrm()
	var l []*bridageModels.TimeTaskRecode
	if count, err = o.QueryTable(new(bridageModels.TimeTaskRecode)).Filter("Manager", manager).Filter("Status", status).All(&l); err != nil {
		logs.Error("GetAllRecode failed, err is ", err.Error())
		return nil, 0, err
	}
	return l, count, nil
}
