package auth

import "fmt"

// Auth ...
type Auth interface {
	Auth(...string) (bool, error) //认证结果和错误与否信息
}

var authMap = make(map[string]Auth)

// Register ...
func Register(authType string, auth Auth) {
	if auth == nil {
		panic("auth: Register failed, auth is nil")
	}
	if _, ok := authMap[authType]; ok {
		panic("auth: authType has called twice" + authType)
	}
	authMap[authType] = auth
}

// GetAuthIns ...
func GetAuthIns(authType string) (Auth, error) {
	if authType == "" {
		panic("auth: authType is null")
	}
	if auth, ok := authMap[authType]; ok {
		return auth, nil
	}
	return nil, fmt.Errorf("GetAuthIns failed, authtype is %s ", authType)
}
