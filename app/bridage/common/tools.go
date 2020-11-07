package common

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"strings"
	"weiXinBot/app/bridage/constant"

	"github.com/astaxie/beego/httplib"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

// Base64Encode ...
func Base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

// Base64Decode ...
func Base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}

// EncodeMD5 ...
// encodeMD5 md5加密
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// PraseXMLString ...
// 解析XML的内容 content->  21592794431@chatroom:xml
func PraseXMLString(content string) (wxsysmsg *WxSysMsg, err error) {
	defer func() {
		if verr := recover(); verr != nil {
			fmt.Println(verr)
		}
	}()
	var conSlice []string
	if conSlice = strings.SplitN(content, ":", 2); len(conSlice) < 2 {
		err = fmt.Errorf("PraseXMLString contentfromproto[%s] is not right, please cheack it", content)
		return nil, err
	}
	var sysmsg WxSysMsg
	if err = xml.Unmarshal([]byte(conSlice[1]), &sysmsg); err != nil {
		// logs.Error("xml.Unmarshal failed, err is ", err.Error())
		return nil, err
	}
	return &sysmsg, nil
}

// ListenGrpcStatus ...监听grpc状态通报
func ListenGrpcStatus() {
	httplib.Get(strings.ReplaceAll(constant.BASE_URL[:7]+constant.GRPC_BASE_URL, "8081", "10052")).DoRequest()
}
