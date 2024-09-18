package logic

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
	"GoChat/apps/im/ws/ws"
	"GoChat/pkg/wuid"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ConversationLogic struct {
	ctx    context.Context
	srv    *websocket.Server
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConversationLogic(ctx context.Context, srv *websocket.Server, svcCtx *svc.ServiceContext) *ConversationLogic {
	return &ConversationLogic{
		ctx:    ctx,
		srv:    srv,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConversationLogic) SingleChat(data *ws.Chat, userId string) error {
	if data.ConversationId == "" {
		data.ConversationId = wuid.CombineId(userId, data.RecvId)
	}

	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         userId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MsgType,
		MsgContent:     data.Content,
		SendTime:       time.Now().Unix(),
	}
	err := l.svcCtx.ChatLogModel.Insert(l.ctx, &chatLog)

	return err
}
