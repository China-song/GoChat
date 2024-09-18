package msgtransfer

import (
	"GoChat/apps/im/ws/websocket"
	"GoChat/apps/im/ws/ws"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/apps/task/mq/internal/svc"
	"GoChat/pkg/constants"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type baseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svcCtx *svc.ServiceContext) *baseMsgTransfer {
	return &baseMsgTransfer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *baseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.SingleChatType:
		err = m.single(ctx, data)
	case constants.GroupChatType:
		err = m.group(ctx, data)
	}
	return err
}

func (m *baseMsgTransfer) single(ctx context.Context, data *ws.Push) error {
	// 推送消息
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SystemRootUid,
		Data:      data,
	})
}

func (m *baseMsgTransfer) group(ctx context.Context, data *ws.Push) error {
	// 先获取群成员ID
	// 查询群用户
	users, err := m.svcCtx.Social.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	// 获取待发送的群用户ID
	data.RecvIds = make([]string, 0, len(users.List))
	for _, user := range users.List {
		// 不包含发送者自己
		if user.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, user.UserId)
	}
	// 推送消息
	return m.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SystemRootUid,
		Data:      data,
	})
}
