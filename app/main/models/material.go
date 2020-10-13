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

// AddMaterial ...
func AddMaterial(materials []*bridageModels.Material) (id int64, err error) {
	o := orm.NewOrm()
	for _, material := range materials {
		if id, err = o.Insert(material); err != nil {
			return 0, err
		}
	}
	return
}

// GetMaterialByID ...
func GetMaterialByID(id int64) (v *bridageModels.Material, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Material{ID: id}
	if err = o.Read(v); err != nil {
		return nil, err
	}
	return v, nil
}

// GetAllMaterial get all bots
func GetAllMaterial(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.Material))
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

	var l []bridageModels.Material
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

// UpdateMaterialByID ...
func UpdateMaterialByID(m *bridageModels.Material) (err error) {
	var v = bridageModels.Material{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}

// MultiUpdateMaterial ...
func MultiUpdateMaterial(ms []*bridageModels.Material) (err error) {
	for _, m := range ms {
		var v = bridageModels.Material{ID: m.ID}
		o := orm.NewOrm()
		if err = o.Read(&v); err == nil {
			var num int64
			if num, err = o.Update(m); err == nil {
				logs.Debug("Number of Bot update in database:", num)
			}
		}
	}
	return
}

// DeleteMaterialByID deletes Material by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMaterialByID(id int64) (err error) {
	o := orm.NewOrm()
	v := bridageModels.Material{ID: id}
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&bridageModels.Material{ID: id}); err == nil {
			logs.Debug("Number of Bots deleted in database:", num)
		}
	}
	return
}

// MultiDeleteMaterialByIDs multi delete Bot
func MultiDeleteMaterialByIDs(ids []interface{}) (err error) {
	o := orm.NewOrm()
	var num int64
	if num, err = o.QueryTable(new(bridageModels.Material)).Filter("ID__in", ids...).Delete(); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}
