package svc

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	immodels.ConversationModel
	immodels.ConversationsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		ChatLogModel:       immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.Db),
		ConversationModel:  immodels.MustConversationModel(c.Mongo.Url, c.Mongo.Db),
		ConversationsModel: immodels.MustConversationsModel(c.Mongo.Url, c.Mongo.Db),
	}
}
