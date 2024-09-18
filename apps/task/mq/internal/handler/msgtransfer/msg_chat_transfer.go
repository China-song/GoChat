package msgtransfer

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/ws/ws"
	"GoChat/apps/task/mq/internal/svc"
	"GoChat/apps/task/mq/mq"
	"GoChat/pkg/bitmap"
	"GoChat/pkg/constants"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svcCtx *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svcCtx),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Printf("key: %s, value: %s\n", key, value)

	var (
		data mq.MsgChatTransfer
		//ctx   = context.Background()
		//msgId = primitive.NewObjectID()
	)
	// 反序列化数据
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		m.Errorf("task mq - msgChatTransfer - Consume - json.Unmarshal err: %v", err)
		return err
	}
	// 记录数据
	msgId, err := m.addChatLog(ctx, &data)
	if err != nil {
		m.Errorf("task mq - msgChatTransfer - Consume - m.addChatLog err: %v", err)
		return err
	}
	// 转发数据
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MsgId:          msgId.Hex(),
		MsgType:        data.MsgType,
		ContentType:    constants.ContentChatMsgType,
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, data *mq.MsgChatTransfer) (msgId primitive.ObjectID, err error) {
	// 记录消息

	// todo 消息的存储和会话的更新 放到事务中
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		MsgFrom:        0,
		ChatType:       data.ChatType,
		MsgType:        data.MsgType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	readRecord := bitmap.NewBitmap(0)
	// 发送者自己发送的消息对其本人是已读的
	readRecord.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecord.Export()

	err = m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		m.Errorf("task mq - msgChatTransfer - Consume - m.addChatLog - ChatLogModel.Insert err: %v", err)
		return
	}
	err = m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
	if err != nil {
		m.Errorf("task mq - msgChatTransfer - Consume - m.addChatLog - ConversationModel.UpdateMsg err: %v", err)
		return
	}
	return chatLog.ID, nil
}
