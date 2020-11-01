package controllers

import (
	"encoding/json"
	"weiXinBot/app/bridage/common/base"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/main/models"
)

// TimeTaskController ...
type TimeTaskController struct {
	base.BaseController
}

// Post ...
// @router / [post]
func (c *TimeTaskController) Post() {
	var v bridageModels.TimeTask
	var err error
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	models.AddTimeTask(&v)
}
