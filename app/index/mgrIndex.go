package index

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/index/auth"
)

// MgrIndexController ...
type MgrIndexController struct {
	base.BaseController
}

// Login ...
func (c *MgrIndexController) Login() {
	var account, password string
	var psa bool
	var err error
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if account = c.GetString("account"); account == "" {
		err = fmt.Errorf("account cant nil")
		return
	}
	if password = c.GetString("password"); password == "" {
		err = fmt.Errorf("password cant nil")
		return
	}
	var _auth auth.Auth
	if _auth, err = auth.GetAuthIns(c.Ctx.Input.Param(":authtype")); err != nil {
		return
	}
	if psa, err = _auth.Auth([]string{account, password}...); err != nil || psa == false {
		return
	}
	c.Ctx.Input.CruSession.Set(constant.S_ACCOUNT, account)
}

// Register ...
func (c *MgrIndexController) Register() {
	var newMgr bridageModels.Manager
	var err error
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()

	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &newMgr); err != nil {

		return
	}

	//查询是否存在手机号注册记录
	bool := bridageModels.FindManagerByTel(newMgr.Tel)

	if bool == true {
		err = fmt.Errorf("手机号已经被注册")
		return
	}

	newMgr.PassWord = encodeMD5(newMgr.PassWord)

	_, err = bridageModels.AddManager(&newMgr)
}

/**
md5加密
*/
func encodeMD5(value string) string {

	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}
