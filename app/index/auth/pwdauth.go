package auth

import (
	"fmt"
	"weiXinBot/app/bridage/common"
	"weiXinBot/app/bridage/constant"
	bridageModels "weiXinBot/app/bridage/models"
)

// PwdAuth password auth
type PwdAuth struct{}

func init() {
	Register(constant.AUTH_PWD, &PwdAuth{})
}

// Auth 账号密码认证
// authParams = []string{account, password}
func (c *PwdAuth) Auth(authParams ...string) (psa bool, err error) {
	switch len(authParams) {
	case 0:
		err = fmt.Errorf("account or password is nil")
		return false, err
	case 1:
		err = fmt.Errorf("auth params is short[%v]", authParams)
		return false, err
	case 2:
		var mgr *bridageModels.Manager
		if mgr, err = bridageModels.GetManagerByAccount(authParams[0]); err == nil {
			if pwd := string(common.Base64Encode([]byte(authParams[1]))); pwd != mgr.PassWord {
				err := fmt.Errorf("账号或密码错误")
				return false, err
			}
			return true, nil
		}
		return false, err
	default:
		err = fmt.Errorf("参数错误")
		return false, err
	}
}
