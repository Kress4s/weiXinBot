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
	bridage "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
)

// IndexController ...
type IndexController struct {
	base.BaseController
}

// Index ...
func (c *IndexController) Index() {
	c.TplName = "index.html"
}

//GetQrCode ...
/*
1. header中找微信号(有直接orm设备id)
2. header中无微信号，前端header中给设备id
3. 目前无加密措施，后期安全性考虑加密
*/
func (c *IndexController) GetQrCode() {
	var deviceID string
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
	wxID := c.GetString(constant.H_WXID)
	if len(wxID) != 0 {
		// 登录过
		if deviceID, err = bridage.GetDeviceIDByWxID(wxID); err != nil {
			logs.Error("GetDeviceIDByWxID failed, err is ", err.Error())
			return
		}
	} else {
		deviceID = c.GetString(constant.H_DEVID)
		if len(deviceID) == 0 {
			err = fmt.Errorf("header[%s] cant be null", constant.H_DEVID)
			return
		}
	}
	if resp, err = httplib.Get(constant.LOGIN_QRCODE_URL).Param(constant.P_DEVICE_ID, deviceID).DoRequest(); err != nil {
		logs.Error("get URL[%s] failed, err is ", err.Error())
		return
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
	// session记录uuid和device的关系(更新user) conf配置session
	// if err = c.Ctx.Input.CruSession.Set(constant.UUID, deviceID); err != nil {
	// 	logs.Error("Set session[%s] failed, err is ", constant.UUID, err.Error())
	// 	return
	// }
	// if isNeedAdd {
	// 	// 把deviceID存进User表
	// 	_, err = bridage.AddUserByDeviceID(deviceID)
	// }
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
		} `json:"data"`
	}
	var restBody *restQRcode
	var err error
	UUID := c.GetString(constant.H_UUID)
	qrFlag := c.GetString("flag")
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: restBody}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if len(UUID) == 0 {
		err = fmt.Errorf("uuid is null")
		return
	}
	for {
		var resp *http.Response
		var verr error
		if resp, verr = httplib.Get(constant.LOGIN_CHECK_URL).Param(constant.P_UUID, UUID).DoRequest(); verr != nil {
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
		if restBody.Code == 0 {
			if qrFlag == "first" && restBody.Data.Token == "" {
				break
			} else if qrFlag == "second" && restBody.Data.Token != "" {
				break
			}
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
