package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/grpc"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// AddBot ...
func AddBot(bot *bridageModels.Bots) (id int64, err error) {
	if len(bot.WXID) == 0 {
		return 0, errors.New("添加机器人的wxid不能为空")
	}
	o := orm.NewOrm()
	grpcToken := bot.Token
	bot.Token = fmt.Sprintf("Bearer %s", bot.Token)
	// 判断wxid不能为空；判断是否存在(是否请求)；存在更新；不存在新增；先是后端处理
	if bridageModels.IsManagerNewBot(bot) {
		if err = StartListenGRPC(grpcToken, bot.WXID); err != nil {
			return
		}
		// update bot info
		if err = bridageModels.UpdateBotByWXID(bot); err != nil {
			logs.Error("when add bot interface update bot accured error, err is ", err.Error())
			return 0, err
		}
		// 这里存在id返回0的不准确的问题，暂时搁置，不影响
	} else {
		if err = StartListenGRPC(grpcToken, bot.WXID); err != nil {
			return
		}
		if id, err = o.Insert(bot); err != nil {
			return 0, err
		}
		return
	}
	return
}

// StartListenGRPC ...
func StartListenGRPC(grpcToken, WXID string) (err error) {
	var isNeed bool
	if isNeed, err = IsNeedRestart(WXID); err != nil {
		logs.Error("grpc %s failed, err is %s", WXID, err.Error())
		return err
	}
	if isNeed {
		// 开启监听此微信号
		botWork := grpc.NewBotWorker()
		fmt.Printf("bot token is %s\n", grpcToken)
		botWork.PrepareParams(grpcToken, WXID)
		// goroutine 监听
		go botWork.Run()
	}
	return
}

// IsNeedRestart ...
func IsNeedRestart(WXID string) (bool, error) {
	o := orm.NewOrm()
	var err error
	var bot = bridageModels.Bots{WXID: WXID}
	if err = o.Read(&bot, "WXID"); err == nil {
		if bot.LoginStatus == 0 {
			return true, nil
		}
	} else if err == orm.ErrNoRows {
		return true, nil
	}
	return false, err
}

// GetBotByID ...
func GetBotByID(id int64) (v *bridageModels.Bots, err error) {
	o := orm.NewOrm()
	v = &bridageModels.Bots{ID: id}
	if err = o.Read(v); err != nil {
		return nil, err
	}
	return v, nil
}

// GetAllBots get all bots
func GetAllBots(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(bridageModels.Bots))
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

	var l []bridageModels.Bots
	qs = qs.OrderBy(sortFields...).RelatedSel()
	// add IsDelete filter
	if _, err = qs.Limit(limit, offset-1).Filter("IsDeleted", false).All(&l, fields...); err == nil {
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

// UpdateBotByID ...
func UpdateBotByID(m *bridageModels.Bots) (err error) {
	var v = bridageModels.Bots{ID: m.ID}
	o := orm.NewOrm()
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of Bot update in database:", num)
		}
	}
	return
}

// DeleteBotByID deletes bot by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBotByID(id int64) (err error) {
	o := orm.NewOrm()
	v := bridageModels.Bots{ID: id}
	if err = o.Read(&v); err == nil {
		var num int64
		// 不做物理删除
		// if num, err = o.Delete(&bridageModels.Bots{ID: id}); err == nil {
		// 	logs.Debug("Number of Bots deleted in database:", num)
		// }

		// 逻辑删除
		v.IsDeleted = true
		if num, err = o.Update(&v); err == nil {
			logs.Debug("Number of Bots deleted in database:", num)
		}
	}
	return
}

// MultiDeleteBotsByIDs multi delete Bot
func MultiDeleteBotsByIDs(ids []interface{}) (err error) {
	o := orm.NewOrm()
	var num int64
	// 物理删除
	// if num, err = o.QueryTable(new(bridageModels.Bots)).Filter("ID__in", ids...).Delete(); err == nil {
	// 	logs.Debug("Number of Bots deleted in database:", num)
	// 	return nil
	// }

	//逻辑删除
	if num, err = o.QueryTable(new(bridageModels.Bots)).Filter("ID__in", ids...).Update(
		orm.Params{
			"IsDeleted": true,
		}); err == nil {
		logs.Debug("Number of Bots deleted in database:", num)
		return nil
	}
	return err
}
