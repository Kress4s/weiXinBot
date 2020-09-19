package common

//RestResult Rest接口返回值
type RestResult struct {
	Code    int         // 0 表示成功，其他失败
	Message string      // 错误信息
	Data    interface{} // 数据体
}

// StandardRestResult 标准的 rest 返回接口，字符小写化
type StandardRestResult struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
