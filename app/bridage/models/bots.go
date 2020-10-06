package models

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Bots ...
type Bots struct {
	ID             int64      `orm:"auto;column(id)"`
	WXID           string     `orm:"size(200);column(wx_id);unique"`
	BigHeadImage   string     `orm:"size(200);column(big_head_img_url)"`
	SmallHeadImage string     `orm:"size(200);column(small_head_img_url)"`
	NickName       string     `orm:"size(50);column(nick_name)"`
	Country        string     `orm:"size(50);column(country)"`
	City           string     `orm:"size(50);column(city)" `
	Sex            bool       `orm:"column(sex)"`
	Signature      string     `orm:"size(50);column(signature)"`
	Alias          string     `orm:"size(50);column(alias)"`
	Device         string     `orm:"size(50);column(device)"`
	LoginStatus    int        `orm:"column(login_status)"`                            // 机器人登录状态
	Status         int        `orm:"column(status)"`                                  // 机器人状态
	LoginTime      time.Time  `orm:"auto_now;column(login_time);type(datetime)"`      //
	CreateTime     time.Time  `orm:"auto_now_add;column(create_time);type(datetime)"` //
	ExpireTime     time.Time  `orm:"type(datetime);column(expiretime);null"`          // 到期时间
	Token          string     `orm:"size(50);column(token)"`                          // Token
	IsDeleted      bool       `orm:"column(is_deleted); default(0)"`                  //逻辑删除字段
	Manager        *Manager   `orm:"rel(fk)"`
	Contacts       []*Contact `orm:"reverse(many)"` //好友
	Groups         []*Group   `orm:"reverse(many)"`
}

// TableIndex ...
// 多字段索引
func (u *Bots) TableIndex() [][]string {
	return [][]string{
		[]string{"Id", "WXID"},
	}
}

func init() {
	orm.RegisterModel(new(Bots))
}

// GetDeviceIDByWxID ...
func GetDeviceIDByWxID(wxid string) (deviceID string, err error) {
	o := orm.NewOrm()
	var user = Bots{WXID: wxid}
	if err = o.Read(&user, "WXID"); err != nil {
		logs.Error("get user by WXID failed, err is ", err.Error())
		return "", err
	}
	return user.Device, nil
}

// AddUserByDeviceID ...
func AddUserByDeviceID(deviceID string) (id int64, err error) {
	o := orm.NewOrm()
	user := Bots{Device: deviceID}
	id, err = o.Insert(&user)
	return
}

// UpdateUserByLoginCheckFunc ...
// 特殊方法(不规范，后续修改)
func UpdateUserByLoginCheckFunc(m *Bots) (err error) {
	o := orm.NewOrm()
	v := Bots{WXID: m.WXID}
	if err = o.Read(&v, "WXID"); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of User update in database:", num)
		}
	}
	return
}
