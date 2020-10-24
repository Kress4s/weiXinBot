package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"weiXinBot/app/bridage/common"
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
	c.Token = token
}

// Run 开始监听
func (c *BotWorker) Run() {
	var message common.ProtoMessage
	var err error
	ctx, cancle := context.WithCancel(context.Background())
	defer func() {
		/*
			TODO:
			异常退出
			1. 退出当前goroutine
			2. 更改数据库机器人的状态
			3. 记录日志(微信号、掉线时间)
			4. 通过websoket方式通知web端掉线的微信号
		*/
		cancle() //通知所有的goroutine退出
		if err = bridageModels.UpdateBotLoginStatusByWXID(c.BotID); err == nil {
			logs.Info("%s has offlined, please check it to relogin", c.BotID)
		}
		// wetsocket 通知前端

	}()
	grpcClient := pb.NewRockRpcServerClient(conn)
	req := pb.StreamRequest{
		Token: &c.Token,
	}
	res, verr := grpcClient.Sync(context.Background(), &req)
	if verr != nil {
		log.Fatalf("Call Route err: %v", verr)
	}
	for {
		response, _ := res.Recv()
		if err = json.Unmarshal([]byte(*response.Payload), &message); err == nil {
			// 开始执行监控操作(后期对message进行类型的解析，避免开启多余的goroutine资源)
			// 目前只需要处理
			/*
				MsgType = 1    (群或者联系人发送文本消息)
				MsgType = 10002(新人入群和踢人出群)
			*/
			// 目前只处理 MsgType = 1 或者 10002消息类型
			if message.MsgType == 1 || message.MsgType == 10002 {
				logs.Info("确认开始服务...")
				go BeginServer(ctx, message)
			} else {
				logs.Info("message type is not need to deal with")
			}
		} else {
			logs.Error("json Unmarshal meaasge failed, err is ", err.Error())
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

// BeginServer By Message
func BeginServer(ctx context.Context, message common.ProtoMessage) {
	var isNeedServer bool
	var err error
	select {
	case <-ctx.Done():
		logs.Debug("one Bot quit...")
		return
	default:
		if isNeedServer, err = bridageModels.GroupIsNeedServer(message); err != nil {
			// if isNeedServer, err = bridageModels.GroupIsNeedServer(message.FromUserName.Str, message.ToUserName.Str); err != nil {
			logs.Error("GroupIsNeedServer failed, err is", err.Error())
		} else if isNeedServer && err == nil {
			logs.Info("确认需要群服务...")
			bridageModels.GroupService(message)
		}
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
