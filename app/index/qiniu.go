package index

import (
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/common/base"

	"github.com/qiniu/api.v7/v7/auth"
	"github.com/qiniu/api.v7/v7/storage"
)

// QiniuController ...
type QiniuController struct {
	base.BaseController
}

var (
	accessKey = "dzHmJvsrtk283rHfAcxHZvlEtCt6Jx-J703W1SGV"
	secretKey = "ek8esIHQ3CmoSRyIuIKJ1azl-oC2NGc5a4lSM4U1"
	bucket    = "cuncui"
)

// GetUploadToken ...
// @router / [get]
func (c *QiniuController) GetUploadToken() {
	// 简单上传凭证
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := auth.New(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	c.Data["json"] = common.RestResult{Code: 0, Message: "ok", Data: upToken}
	c.ServeJSON()
}
