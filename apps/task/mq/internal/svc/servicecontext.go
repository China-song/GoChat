package svc

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/ws/websocket"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/apps/task/mq/internal/config"
	"GoChat/pkg/constants"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"net/http"
)

type ServiceContext struct {
	Config config.Config

	WsClient websocket.Client
	*redis.Redis

	immodels.ChatLogModel
	immodels.ConversationModel

	socialclient.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.Redisx),

		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
	// 获取token
	token, err := svcCtx.GetSystemToken()
	if err != nil {
		panic(err)
	}
	// 设置token
	header := http.Header{}
	header.Set("Authorization", token)
	// 创建Websocket客户端
	svcCtx.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))
	return svcCtx
}

func (svcCtx *ServiceContext) GetSystemToken() (string, error) {
	return svcCtx.Redis.Get(constants.RedisSystemRootToken)
}
