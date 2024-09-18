package push

import (
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
	"GoChat/apps/im/ws/ws"
	"GoChat/pkg/constants"
	"github.com/mitchellh/mapstructure"
)

// Push msg from mq to client
func Push(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		// 解析来自mq的msg
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("send err: %v", err)
			}
			return
		}

		// 发送给client
		switch data.ChatType {
		case constants.SingleChatType:
			err := single(srv, &data, data.RecvId)
			if err != nil {
				srv.Errorf("single err: %v", err)
				return
			}
		case constants.GroupChatType:
			group(srv, &data)
		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	// 获取发送的目标
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		// todo: 目标离线
		srv.Info("目标离线")
		return nil
	}
	// 发送消息
	srv.Infof("push msg `%v` to `%s`", data, recvId)
	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MsgId:       data.MsgId,
			MsgType:     data.MsgType,
			Content:     data.Content,
			ReadRecords: data.ReadRecords,
		},
	}), rconn)
}

func group(srv *websocket.Server, data *ws.Push) {
	for _, id := range data.RecvIds {
		func(recvId string) {
			srv.Schedule(func() {
				err := single(srv, data, recvId)
				if err != nil {
					srv.Errorf("push err: %v", err)
					return
				}
			})
		}(id)
	}
	return
}
