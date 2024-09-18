package mqclient

import (
	"GoChat/apps/task/mq/mq"
	"context"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
)

type MsgReadTransferClient interface {
	Push(msg *mq.MsgMarkRead) error
}

type msgReadTransferClient struct {
	pusher *kq.Pusher
}

func NewMsgReadTransferClient(addr []string, topic string, opts ...kq.PushOption) MsgReadTransferClient {
	return &msgReadTransferClient{
		pusher: kq.NewPusher(addr, topic, opts...),
	}
}

func (c *msgReadTransferClient) Push(msg *mq.MsgMarkRead) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.pusher.Push(context.Background(), string(body))
}
