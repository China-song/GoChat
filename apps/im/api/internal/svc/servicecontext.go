package svc

import (
	"GoChat/apps/im/api/internal/config"
	"GoChat/apps/im/rpc/imclient"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/apps/user/rpc/userclient"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
