package models

import (
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// Bots ...
type Bots struct {
	ID             int64      `orm:"auto;column(id)"`
	WXID           string     `orm:"size(200);column(wx_id);unique"`                  // 原始微信ID
	BigHeadImage   string     `orm:"size(200);column(big_head_img_url)"`              //
	SmallHeadImage string     `orm:"size(200);column(small_head_img_url)"`            //
	NickName       string     `orm:"size(50);column(nick_name)"`                      //
	Country        string     `orm:"size(50);column(country)"`                        //
	City           string     `orm:"size(50);column(city)" `                          //
	Sex            bool       `orm:"column(sex)"`                                     //
	Signature      string     `orm:"size(50);column(signature)"`                      //
	Alias          string     `orm:"size(50);column(alias)"`                          // 自己设置的微信号
	Device         string     `orm:"size(50);column(device)"`                         //
	LoginStatus    int        `orm:"column(login_status)"`                            // 机器人登录状态
	Status         int        `orm:"column(status)"`                                  // 机器人状态
	LoginTime      time.Time  `orm:"auto_now;column(login_time);type(datetime)"`      //
	CreateTime     time.Time  `orm:"auto_now_add;column(create_time);type(datetime)"` //
	ExpireTime     time.Time  `orm:"type(datetime);column(expiretime);null"`          // 到期时间
	Token          string     `orm:"size(50);column(token)"`                          // Token
	IsDeleted      bool       `orm:"column(is_deleted); default(0)"`                  //逻辑删除字段
	Manager        *Manager   `orm:"rel(fk)"`                                         //
	Contacts       []*Contact `orm:"reverse(many)"`                                   //好友
	Groups         []*Group   `orm:"reverse(many)"`                                   //
}

// TableUnique ...
// 多字段唯一键
func (u *Bots) TableUnique() [][]string {
	return [][]string{
		{"WXID", "Manager"},
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

// UpdateBotByWXID ...
func UpdateBotByWXID(m *Bots) (err error) {
	o := orm.NewOrm()
	v := Bots{WXID: m.WXID}
	if err = o.Read(&v, "WXID"); err == nil {
		// update needs ID
		m.ID = v.ID
		var num int64
		if num, err = o.Update(m); err == nil {
			logs.Debug("Number of User update in database:", num)
		}
	}
	return
}

// IsManagerNewBot 判断该用户下是否存在过这个机器人
func IsManagerNewBot(m *Bots) (isExist bool) {
	o := orm.NewOrm()
	if !o.QueryTable(new(Bots)).Filter("WXID", m.WXID).Filter("Manager", m.Manager).Exist() {
		return false
	}
	return true
}

// GetBotByWXID ...
func GetBotByWXID(WXID string) (v *Bots, err error) {
	o := orm.NewOrm()
	var bot = &Bots{WXID: WXID}
	if err = o.Read(bot, "WXID"); err != nil {
		logs.Error("SendText: get bot info by WXID failed, err is ", err.Error())
		return nil, err
	}
	return bot, nil
}
