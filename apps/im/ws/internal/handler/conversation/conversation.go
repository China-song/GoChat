package conversation

import (
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
	"GoChat/apps/im/ws/ws"
	"GoChat/apps/task/mq/mq"
	"GoChat/pkg/constants"
	"GoChat/pkg/wuid"
	"github.com/mitchellh/mapstructure"
	"time"
)

func Chat(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			err = srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("server send message error: %v", err)
			}
			return
		}

		// 如果传递了会话ID，直接发送，否则分类型创建会话ID
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
			case constants.GroupChatType:
				// 群聊的会话ID 是 群聊ID
				data.ConversationId = data.RecvId
			}
		}

		// 交给消息队列
		err := svcCtx.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,

			SendTime: time.Now().Unix(),
			MsgType:  data.Msg.MsgType,
			Content:  data.Msg.Content,
		})
		if err != nil {
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}
		srv.Info("发送成功")
	}
}

func MarkRead(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			err = srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("server send message error: %v", err)
			}
			return
		}

		// 交给消息队列
		err := svcCtx.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})
		if err != nil {
			err := srv.Send(websocket.NewErrMessage(err), conn)
			if err != nil {
				srv.Errorf("error message send error: %v", err)
			}
			return
		}
		srv.Info("发送成功")
	}
}
