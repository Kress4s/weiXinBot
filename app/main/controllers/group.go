package controllers

import (
	"encoding/json"
	"fmt"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"
	"weiXinBot/app/main/models"
)

// GroupController ...
type GroupController struct {
	base.BaseController
}

// Post ...
// @Title Post
// @Description create Group
// @Param	body		body 	models.Group	true		"body for Group content"
// @Success 201 {int} models.Group
// @Failure 403 body is empty
// @router / [post]
func (c *GroupController) Post() {
	var v *bridageModels.Group
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
	_, err = models.AddGroup(v)
}

// MultiPost ...
// @router /multi [post]
func (c *GroupController) MultiPost() {
	type Groups struct {
		Data []*bridageModels.Group `json:"data"`
	}
	var v Groups
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
	err = models.MultiAddGroup(v.Data)
}

// GetOne ...
// @Title Get One
// @Description get Group by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Group
// @Failure 403 :id is empty
// @router /:gid [get]
func (c *GroupController) GetOne() {
	var v *bridageModels.Group
	var err error
	gid := c.Ctx.Input.Param(":gid")
	defer func() {
		if err == nil {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		} else {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		}
		c.ServeJSON()
	}()
	fmt.Println(gid)
	v, err = models.GetGroupByGID(gid)
}

// GetGroupFromProto ...
// @router /groupfromproto [get]
func (c *GroupController) GetGroupFromProto() {
	var v interface{}
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: v}
		}
		c.ServeJSON()
	}()
	Authorization := c.Ctx.Input.Header(constant.H_TOKEN_KEY)
	v, err = bridageModels.ProtoGiveGroup(Authorization)
}

// Put ...
// @router /:gid [put]
func (c *GroupController) Put() {
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	gid := c.Ctx.Input.Param(":gid")
	var v = bridageModels.Group{GID: gid}
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	err = models.UpdateGrouByID(&v)
}

// MultiPut ...
// @router /updatemulti [put]
func (c *GroupController) MultiPut() {
	var moveOutGroups string
	var err error
	defer func() {
		if err != nil {
			c.Data["json"] = common.RestResult{Code: -1, Message: err.Error()}
		} else {
			c.Data["json"] = common.RestResult{Code: 0, Message: "ok"}
		}
		c.ServeJSON()
	}()
	type Groups struct {
		Data []*bridageModels.Group `json:"Data"`
	}
	// delGroupIDS:   GID:BOTID,GID2:BOTID2...
	if moveOutGroups = c.GetString("delGroupIDS"); len(moveOutGroups) == 0 {
		err = fmt.Errorf("moveOutGroups cant be null")
		return
	}
	var v Groups
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		return
	}
	err = models.MultiUpdateGrouByID(v.Data, moveOutGroups)
}
