package handler

import (
	"GoChat/apps/task/mq/internal/handler/msgtransfer"
	"GoChat/apps/task/mq/internal/svc"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct {
	svcCtx *svc.ServiceContext
}

func (l *Listen) Services() []service.Service {
	return []service.Service{
		// todo: 此处可以加载多个消费者
		kq.MustNewQueue(l.svcCtx.Config.MsgChatTransfer, msgtransfer.NewMsgChatTransfer(l.svcCtx)),
		kq.MustNewQueue(l.svcCtx.Config.MsgReadTransfer, msgtransfer.NewMsgReadTransfer(l.svcCtx)),
	}
}

func NewListen(svcCtx *svc.ServiceContext) *Listen {
	return &Listen{
		svcCtx: svcCtx,
	}
}
