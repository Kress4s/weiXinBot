package index

import (
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
	var psa bool
	var err error
	var token string
	var manager *bridageModels.Manager
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: struct {
				Token string
				ID    int64
			}{
				Token: token,
				ID:    manager.ID,
			}}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	type LoginParams struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	var loginPams LoginParams
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &loginPams); err != nil {
		err = fmt.Errorf("login body is error")
		return
	}
	var _auth auth.Auth
	if _auth, err = auth.GetAuthIns(c.Ctx.Input.Param(":authtype")); err != nil {
		return
	}
	if psa, err = _auth.Auth([]string{loginPams.Account, loginPams.Password}...); err != nil || psa == false {
		return
	}
	if manager, err = bridageModels.GetManagerByAccount(loginPams.Account); err != nil {
		return
	}
	jwt := common.NewJWT()
	myClaim := common.CreatCliamIns(loginPams.Account, loginPams.Password)
	token, err = jwt.CreateToken(myClaim)
	// c.Ctx.Input.CruSession.Set(constant.S_ACCOUNT, loginPams.Account)
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
	if isExist := bridageModels.FindManagerByTel(newMgr.Tel); isExist {
		err = fmt.Errorf("手机号已经被注册")
		return
	}
	newMgr.PassWord = string(common.Base64Encode([]byte(newMgr.PassWord)))
	_, err = bridageModels.AddManager(&newMgr)
}

// GetMyInfo 获取登录用户详情
func (c *MgrIndexController) GetMyInfo() {
	var v interface{}
	var account string
	var err error
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if account = c.GetString(constant.S_ACCOUNT); account == "" {
		err = fmt.Errorf("account is nil")
		return
	}
	v, err = bridageModels.GetManagerByAccount(account)
}
