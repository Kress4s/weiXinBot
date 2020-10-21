package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"weiXinBot/app/bridage/constant"
	pb "weiXinBot/app/bridage/grpc/proto"
	bridageModels "weiXinBot/app/bridage/models"

	"github.com/astaxie/beego/logs"
	"google.golang.org/grpc"
)

// BotWorker ... 后期优化
type BotWorker struct {
	// Conn  *grpc.ClientConn
	// Lock  sync.Mutex
	Token string
	BotID string
}

// Message ...
type Message struct {
	FromUserName struct {
		Str string `json:"str"`
	} `json:"from_user_name"` //
	ToUserName struct {
		Str string `json:"str"`
	} `json:"to_user_name"` //
	MsgType string `json:"msg_type"` // 消息类型
	Content struct {
		Str string `json:"str"`
	} `json:"content"` // 内容(我发：{"str":"程序监控你"}；别人发：{"str":"aaaa520jj:\nG吐总冠军"})
	Status      int    `json:"status"`       //貌似群的消息都是
	CreateTime  int    `json:"create_time"`  //消息时间戳
	MsgSource   string `json:"msg_source"`   // ?
	PushContent string `json:"push_content"` //提示消息(聊天输入框提示) (别人发有这个字段，我发没有这个字段)
}

const (
	// Address grpc连接地址
	Address string = constant.GRPC_BASE_URL
)

var conn *grpc.ClientConn // 一个连接
var lock sync.Mutex

// GetConnInstance 获取连接
func GetConnInstance() (*grpc.ClientConn, error) {
	var err error
	lock.Lock()
	defer lock.Unlock()
	if conn == nil {
		if conn, err = grpc.Dial(Address, grpc.WithInsecure()); err != nil {
			logs.Error("create grpc conn failed, err is ", err.Error())
			return nil, err
		}
	}
	return conn, nil
}

// NewBotWorker ...
func NewBotWorker() *BotWorker {
	return new(BotWorker)
}

// PrepareParams 预置参数
func (c *BotWorker) PrepareParams(token, botID string) {
	c.BotID = botID
	c.Token = botID
}

// Run 开始监听
func (c *BotWorker) Run() {
	var message Message
	defer func() {
		runtime.Goexit()
		/*
			TODO:
			异常退出
			1. 退出当前goroutine
			2. 更改数据库机器人的状态
			3. 记录日志(微信号、掉线时间)
			4. 通过websoket方式通知web端掉线的微信号
		*/

	}()
	grpcClient := pb.NewRockRpcServerClient(conn)
	req := pb.StreamRequest{
		Token: &c.Token,
	}
	res, verr := grpcClient.Sync(context.Background(), &req)
	if verr != nil {
		log.Fatalf("Call Route err: %v", verr)
	}
	var isNeedServer bool
	var err error
	for {
		response, _ := res.Recv()
		json.Unmarshal([]byte(*response.Payload), &message)
		if isNeedServer, err = bridageModels.GroupIsNeedServer(message.FromUserName.Str, message.ToUserName.Str); err != nil {
			logs.Error("GroupIsNeedServer failed, err is", err.Error())
		} else if isNeedServer && err == nil {
			bridageModels.GroupService(message.FromUserName.Str, message.ToUserName.Str, message.PushContent)
		}
		/*
			TODO:处理信息内容
			1. 群号，WXID查功能表，不存在直接pass
			2. 存在会是个list，遍历根据type区分功能，查功能表
			3. 分析结果
			4. 根据机器人的微信号做出相应的动作
		*/
	}

}

// CloseConn 关闭
func (c *BotWorker) CloseConn() {
	conn.Close()
}

func init() {
	var err error
	conn, err = GetConnInstance()
	if err != nil {
		panic(fmt.Errorf("get gRPC intance faield,err is %s", err.Error()))
	}
}
