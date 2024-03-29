package models

import "github.com/astaxie/beego/orm"

// GroupPlan ...
type GroupPlan struct {
	ID      int64    `orm:"auto;column(id)"`       //
	Name    string   `orm:"size(50);column(name)"` // 名称
	Manager *Manager `orm:"rel(fk)"`               // 属于哪个用户创建的群管方案
	Groups  []*Group `orm:"reverse(many)"`         //
}

// TableUnique ...
// 多字段唯一键
func (u *GroupPlan) TableUnique() [][]string {
	return [][]string{
		{"Name", "Manager"},
	}
}

func init() {
	orm.RegisterModel(new(GroupPlan))
}

// GetGrouPlanFuncSwitch 看方案下配置的开关状态
func GetGrouPlanFuncSwitch(grouplanID int64) (ret map[int]int, err error) {
	var wel Welcome
	var keyword KeyWords
	// var autokick AutoKickx
	var whitelist WhiteList
	m := make(map[int]int)
	o := orm.NewOrm()
	if err = o.QueryTable(new(Welcome)).Filter("GroupPlan__ID", grouplanID).One(&wel); err != nil {
		if err == orm.ErrNoRows {
			// 为配置功能 状态为-1
			m[0] = -1
		} else {
			return nil, err
		}
	} else {
		if wel.Switch == true {
			m[0] = 1
		} else {
			m[0] = 0
		}
	}
	if err = o.QueryTable(new(KeyWords)).Filter("GroupPlan__ID", grouplanID).One(&keyword); err != nil {
		if err == orm.ErrNoRows {
			// 为配置功能 状态为-1
			m[1] = -1
		} else {
			return nil, err
		}
	} else {
		if keyword.Switch == true {
			m[1] = 1
		} else {
			m[1] = 0
		}
	}
	m[2] = -1 //假的
	// if err = o.QueryTable(new(AutoKick)).Filter("GroupPlan__ID", groupID).One(&autokick); err != nil {
	// 	if err == orm.ErrNoRows {
	// 		// 为配置功能 状态为-1
	// 		m[2] = -1
	// 	} else {
	// 		return nil, err
	// 	}
	// } else {
	// 	if wel.Switch == true {
	// 		m[2] = 1
	// 	} else {
	// 		m[2] = 0
	// 	}
	// }
	if err = o.QueryTable(new(WhiteList)).Filter("GroupPlan__ID", grouplanID).One(&whitelist); err != nil {
		if err == orm.ErrNoRows {
			// 为配置功能 状态为-1
			m[3] = -1
		} else {
			return nil, err
		}
	} else {
		if whitelist.Switch == true {
			m[3] = 1
		} else {
			m[3] = 0
		}
	}
	return m, nil
}

// GetGrouplainFuncInfo ...获取grouplan下面所有配置类型和信息
func GetGrouplainFuncInfo(grouplanID int64) (ret map[string]int64, err error) {
	var wel Welcome
	var keyword KeyWords
	// var autokick AutoKickx
	// var whitelist WhiteList
	m := make(map[string]int64)
	o := orm.NewOrm()
	if err = o.QueryTable(new(Welcome)).Filter("GroupPlan__ID", grouplanID).One(&wel); err == nil {
		m["1"] = wel.ID
	}

	if err = o.QueryTable(new(KeyWords)).Filter("GroupPlan__ID", grouplanID).One(&keyword); err == nil {
		m["2"] = keyword.ID
	}
	// 自动踢人
	// m["3"] = -1

	// 白名单
	// m["4"] = -1
	return m, nil
}
