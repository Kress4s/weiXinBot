package models

import (
	"errors"
	"reflect"
	"strconv"
	"weiXinBot/app/bridage/common"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddQuestion ...
func AddQuestion(v *bridageModels.Question) (id int64, err error) {
	o := orm.NewOrm()
	defer func() {
		if err != nil {
			// mysql data roolback
			o.Rollback()
			return
		}
		o.Commit()
		return
	}()
	if err = o.Begin(); err != nil {
		return
	}
	if id, err = o.Insert(v); err != nil {
		logs.Error("insert question failed, err is ", err.Error())
		return
	}
	// 模糊的插入
	for _, _f := range v.FuzzWords {
		_f.Question = v
		if _, err = o.Insert(_f); err != nil {
			logs.Error("insert FuzzWords failed, err is ", err.Error())
			return
		}
	}
	// 精准的插入
	for _, _e := range v.ExactWords {
		_e.Question = v
		if _, err = o.Insert(_e); err != nil {
			logs.Error("insert ExactWords failed, err is ", err.Error())
			return
		}
	}
	return
}

// GetQuestionByID ...
func GetQuestionByID(id int64) (v *bridageModels.Question, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Question{ID: id}
	if err = o.QueryTable(new(bridageModels.Question)).Filter("ID", id).One(v); err != nil {
		return nil, err
	}
	o.LoadRelated(v, "ExactWords")
	o.LoadRelated(v, "FuzzWords")
	return v, nil
}

// GetAllQuestion get all bots
func GetAllQuestion(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.Question))
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

	var l []bridageModels.Question
	qs = qs.OrderBy(sortFields...).RelatedSel()
	// add IsDelete filter
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				o.LoadRelated(&v, "ExactWords")
				o.LoadRelated(&v, "FuzzWords")
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

// UpdateQuestionByID ...
func UpdateQuestionByID(m *bridageModels.Question) (err error) {
	o := orm.NewOrm()
	var created bool
	defer func() {
		if err != nil {
			// mysql data roolback
			o.Rollback()
			return
		}
		o.Commit()
		return
	}()
	if err = o.Begin(); err != nil {
		return
	}
	// 删除修改时被删除的关键词
	var v = bridageModels.Question{ID: m.ID}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	for _, _e := range m.ExactWords {
		_e.Question = &v
		if created, _, err = o.ReadOrCreate(_e, "Word"); err != nil {
			return
		}
		if created {
			// 已创建 更新
			if err = o.ReadForUpdate(_e); err != nil {
				return
			}
		}
	}
	for _, _f := range m.FuzzWords {
		_f.Question = &v
		if created, _, err = o.ReadOrCreate(_f, "Word"); err != nil {
			return
		}
		if created {
			// 已创建 更新
			if err = o.ReadForUpdate(_f); err != nil {
				return
			}
		}
	}
	return
}

// DeleteQuestionByID deletes bot by Id and returns error if
// the record to be deleted doesn't exist
func DeleteQuestionByID(id int64) (err error) {
	o := orm.NewOrm()
	v := bridageModels.Question{ID: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&bridageModels.Question{ID: id}); err == nil {
			logs.Debug("Number of Bots deleted in database:", num)
		}
	}
	return
}

// MultiDeleteQuestionByIDs multi delete Bot
func MultiDeleteQuestionByIDs(ids []interface{}) (err error) {
	o := orm.NewOrm()
	var num int64
	if num, err = o.QueryTable(new(bridageModels.Question)).Filter("ID__in", ids...).Delete(); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}
