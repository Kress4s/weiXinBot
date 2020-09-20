package controllers

import (
	"encoding/json"
	"strconv"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/main/models"
)

// BotsController ...
type BotsController struct {
	base.BaseController
}

// Post ...
// @Title Post
// @Description create Bots
// @Param	body		body 	models.Bots	true		"body for Bots content"
// @Success 201 {int} models.Bots
// @Failure 403 body is empty
// @router / [post]
func (c *BotsController) Post() {
	var v *bridageModels.Bots
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
	_, err = models.AddBot(v)
}

// GetOne ...
// @Title Get One
// @Description get Bot by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Bot
// @Failure 403 :id is empty
// @router /:id [get]
func (c *BotsController) GetOne() {
	var v *bridageModels.Bots
	var err error
	idstr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idstr, 0, 64)
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	v, err = models.GetBotByID(id)
}
