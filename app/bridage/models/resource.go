package models

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Resource send source
// 素材库(单个组)
type Resource struct {
	ID       int64       `orm:"auto;column(id)"`              //
	Title    string      `orm:"size(30);column(title)"`       //素材标题
	IsPublic int         `orm:"column(is_public);default(0)"` //资源是否是公有（区分素材公共库展示和增加方案功能的具体回复内容）
	Material []*Material `orm:"reverse(many)"`                //
	Manager  string      `orm:"size(30);column(manager)"`     //属于哪个用户的素材库(用户Tel)
}

// Material 单个素材类型
type Material struct {
	ID       int64     `orm:"auto;column(id)"`          //
	Type     int       `orm:"column(type)"`             //消息内容类型(1:文本 2:图片....)
	Data     string    `orm:"size(1000); column(data)"` //类型的具体内容(后期可拆)
	Resource *Resource `orm:"rel(fk)"`                  //所属素材组
}

func init() {
	orm.RegisterModel(new(Resource), new(Material))
}

// GetResourceByIds 根据resource的多个id或许资源信息
func GetResourceByIds(resourceIds string) (ret []*Resource, err error) {
	if len(resourceIds) == 0 {
		err := fmt.Errorf("GetResourceByIds ids cant be null")
		return nil, err
	}
	o := orm.NewOrm()
	var resources []*Resource
	if _, err = o.Raw("select * from resource where ? like concat(id, ',%')"+"or ? like concat('%,', id)"+
		"or ? like concat('%,', id, ',%')"+"or ? = id", resourceIds, resourceIds, resourceIds, resourceIds).QueryRows(&resources); err != nil {
		logs.Error("Get Resource By Ids failed, err is ", err.Error())
		return nil, err
	}
	for _, r := range resources {
		o.LoadRelated(r, "Material")
	}
	return resources, nil
}
