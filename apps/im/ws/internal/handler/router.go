package handler

import (
	"GoChat/apps/im/ws/internal/handler/conversation"
	"GoChat/apps/im/ws/internal/handler/push"
	"GoChat/apps/im/ws/internal/handler/user"
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svcCtx *svc.ServiceContext) {
	srv.AddRoutes(
		[]websocket.Route{
			{
				Method:  "user.online",
				Handler: user.OnLine(svcCtx),
			},
			{
				Method:  "conversation.chat",
				Handler: conversation.Chat(svcCtx),
			},
			{
				Method:  "conversation.markChat",
				Handler: conversation.MarkRead(svcCtx),
			},
			{
				Method:  "push",
				Handler: push.Push(svcCtx),
			},
		})
}
