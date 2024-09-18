package msgtransfer

import (
	"GoChat/apps/im/ws/ws"
	"GoChat/apps/task/mq/internal/svc"
	"GoChat/apps/task/mq/mq"
	"GoChat/pkg/bitmap"
	"GoChat/pkg/constants"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime  = time.Second
	GroupMsgReadRecordDelayCount = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*baseMsgTransfer

	mu sync.Mutex

	groupMsgs map[string]*groupMsgRead
	push      chan *ws.Push
}

func NewMsgReadTransfer(svcCtx *svc.ServiceContext) *MsgReadTransfer {
	m := &MsgReadTransfer{
		baseMsgTransfer: NewBaseMsgTransfer(svcCtx),
		groupMsgs:       make(map[string]*groupMsgRead, 1),
		push:            make(chan *ws.Push, 1),
	}
	// 如果开启
	if svcCtx.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		// 最大计数
		if svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayCount > 0 {
			// 设置值
			GroupMsgReadRecordDelayCount = svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayCount
		}
		// 超时时间
		if svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svcCtx.Config.MsgReadHandler.GroupMsgReadRecordDelayTime) * time.Second
		}
	}

	go m.transfer()

	return m
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Printf("key: %s, value: %s\n", key, value)

	var (
		data mq.MsgMarkRead
	)
	// 反序列化数据
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		m.Errorf("task mq - msgReadTransfer - Consume - json.Unmarshal err: %v", err)
		return err
	}
	// 更新消息的已读情况
	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}
	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkReadType,
		ReadRecords:    readRecords,
	}
	// 判断消息类型
	switch data.ChatType {
	case constants.SingleChatType:
		// 直接推送
		m.push <- push
	case constants.GroupChatType:
		// 判断是否采用合并发送
		// 若不开启
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.push <- push
			break
		}
		// 开启
		m.mu.Lock()
		defer m.mu.Unlock()
		push.SendId = "" // 群聊不需要发送者ID，统一清空
		// 已经有消息
		if _, ok := m.groupMsgs[push.ConversationId]; ok {
			// 和并请求
			m.Infof("merge push: %v", push.ConversationId)
			m.groupMsgs[push.ConversationId].mergePush(push)
		} else {
			// 没有记录，创建
			m.Infof("create merge push %v", push.ConversationId)
			m.groupMsgs[push.ConversationId] = newGroupMsgRead(push, m.push)
		}
	}
	return nil
}

// UpdateChatLogRead return map[msgId]: readRecords
func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	result := make(map[string]string)
	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		m.Errorf("task mq - msgReadTransfer - Consume - UpdateChatLogRead - ChatLogModel.ListByMsgIds err: %v", err)
		return nil, err
	}
	// 处理已读消息
	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			// todo 私聊和群聊的消息已读更新应该不一致？
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			// 设置当前发送者用户为已读状态
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		result[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)

		err = m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			m.Errorf("task mq - msgReadTransfer - Consume - UpdateChatLogRead - ChatLogModel.UpdateMarkRead err: %v", err)
			return nil, err
		}
	}
	return result, nil
}

// 异步处理消息发送
func (m *MsgReadTransfer) transfer() {
	for push := range m.push {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("transfer err: %s", err.Error())
			}
		}
		if push.ChatType == constants.SingleChatType {
			continue
		}
		// 不采用合并推送
		if m.svcCtx.Config.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		// 清空数据
		m.mu.Lock()
		if _, ok := m.groupMsgs[push.ConversationId]; ok && m.groupMsgs[push.ConversationId].IsIdle() {
			m.groupMsgs[push.ConversationId].Clear()
			delete(m.groupMsgs, push.ConversationId)
		}
		m.mu.Unlock()
	}
}
