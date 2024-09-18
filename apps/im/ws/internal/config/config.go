package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	service.ServiceConf
	ListenOn string

	JwtAuth struct {
		AccessSecret string
	}

	Mongo struct {
		Url string
		Db  string
	}

	MsgChatTransfer struct {
		Addrs []string
		Topic string
	}
	MsgReadTransfer struct {
		Addrs []string
		Topic string
	}
}
