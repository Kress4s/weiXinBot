package index

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	"weiXinBot/app/bridage/constant"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

// IndexController ...
type IndexController struct {
	base.BaseController
}

//GetQrCode ...
/*
1. header中找微信号(有直接orm设备id)
2. header中无微信号，前端header中给设备id
3. 目前无加密措施，后期安全性考虑加密
*/
func (c *IndexController) GetQrCode() {
	var token string
	var restBody common.StandardRestResult
	var resp *http.Response
	var err error
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: restBody}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if token = c.GetString(constant.P_TOKEN); token == "" {
		err = fmt.Errorf("token cant be null")
		return
	}
	token = fmt.Sprintf("Bearer %s", token)
	wxID := c.GetString(constant.H_WXID)
	if len(wxID) != 0 {
		// 登录过
		if resp, err = httplib.Get(constant.LOGIN_QRCODE_URL).Param(constant.P_WXID, wxID).Header(constant.H_AUTHORIZATION, token).DoRequest(); err != nil {
			logs.Error("get URL[%s] failed, err is ", err.Error())
			return
		}
	} else {
		// 新微信号
		if resp, err = httplib.Get(constant.LOGIN_QRCODE_URL).Header(constant.H_AUTHORIZATION, token).DoRequest(); err != nil {
			logs.Error("get URL[%s] failed, err is ", err.Error())
			return
		}
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("get URL[%s] body failed, err is ", constant.LOGIN_QRCODE_URL, err.Error())
		return
	}
	if err = json.Unmarshal(body, &restBody); err != nil {
		logs.Error("json.Unmarshal qrcode failed, err is ", err.Error())
		return
	}
}

// Check ...
func (c *IndexController) Check() {
	// var resp *http.Response
	var timeout int
	type restQRcode struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
		Data    struct {
			Alias      string `json:"alias"`
			HeadImgURL string `json:"head_image_url"`
			NickName   string `json:"nick_name"`
			Token      string `json:"token"`
			WXID       string `json:"wx_id"`
			Status     string `json:"status"`
		} `json:"data"`
	}
	var restBody *restQRcode
	var err error
	var token string
	if token = c.GetString(constant.P_TOKEN); token == "" {
		err = fmt.Errorf("token is null")
		return
	}
	token = fmt.Sprintf("Bearer %s", token)
	qrFlag := c.GetString("flag")
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: restBody}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	for {
		var resp *http.Response
		var verr error
		if resp, verr = httplib.Get(constant.LOGIN_CHECK_URL).Header(constant.H_AUTHORIZATION, token).DoRequest(); verr != nil {
			logs.Error("get response[%s] failed, err is ", constant.LOGIN_CHECK_URL, verr.Error())
			return
		}
		var body []byte
		if body, verr = ioutil.ReadAll(resp.Body); verr != nil {
			logs.Error("get URL[%s] body failed, err is ", constant.LOGIN_CHECK_URL, verr.Error())
			return
		}
		if err = json.Unmarshal(body, &restBody); verr != nil {
			logs.Error("json Unmarshal failed, err is ", verr.Error())
			return
		}
		// 正常
		if restBody.Code == 0 && qrFlag == "first" && restBody.Data.Status == "Scanned" {
			break
		} else if restBody.Code == 0 && qrFlag == "second" && restBody.Data.Status == "Confirmed" {
			break
		} else if qrFlag == "cancel" {
			restBody = nil
			break
		}
		// 异常
		time.Sleep(1 * time.Second)
		if timeout == 30 {
			break
		}
		timeout++
		continue
	}
}
