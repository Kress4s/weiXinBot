package controllers

import (
	"encoding/json"
	"strconv"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/main/models"
)

// ConfigurationController ...
type ConfigurationController struct {
	base.BaseController
}

// Post ...
// @Title Post
// @Description create Configuration
// @Param	body		body 	models.Configuration	true		"body for Configuration content"
// @Success 201 {int} models.Configuration
// @Failure 403 body is empty
// @router / [post]
func (c *ConfigurationController) Post() {
	var v *bridageModels.Configuration
	var err error
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	_, err = models.AddConfiguration(v)
}

// Put ...
// @router /:id [put]
func (c *ConfigurationController) Put() {
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	var v = bridageModels.Configuration{ID: id}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	err = models.UpdateConfigurationByID(&v)
}

// MultiUpdateAdd ...
// @router /multiupdate [post]
func (c *ConfigurationController) MultiUpdateAdd() {
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	var v bridageModels.MultiDealConfig
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	err = models.UpdateOrAddConfig(v)
}
