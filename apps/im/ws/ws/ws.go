package ws

import "GoChat/pkg/constants"

type (
	Msg struct {
		MsgId             string `mapstructure:"msgId"`
		constants.MsgType `mapstructure:"msgType"`
		Content           string            `mapstructure:"content"`
		ReadRecords       map[string]string `mapstructure:"readRecords"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string `mapstructure:"sendId"`
		RecvId             string `mapstructure:"recvId"`
		SendTime           int64  `mapstructure:"sendTime"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId     string `mapstructure:"conversationId"`
		constants.ChatType `mapstructure:"chatType"`
		SendId             string   `mapstructure:"sendId"`
		RecvId             string   `mapstructure:"recvId"`
		RecvIds            []string `mapstructure:"recvIds"`
		SendTime           int64    `mapstructure:"sendTime"`

		MsgId             string `mapstructure:"msgId"`
		constants.MsgType `mapstructure:"msgType"`
		ContentType       constants.ContentType `mapstructure:"contentType"`
		Content           string                `mapstructure:"content"`
		ReadRecords       map[string]string     `mapstructure:"readRecords"`
	}

	// MarkRead 处理已读消息
	MarkRead struct {
		constants.ChatType `mapstructure:"chatType"`
		RecvId             string   `mapstructure:"recvId"` // 已读结果推送给谁
		ConversationId     string   `mapstructure:"conversationId"`
		MsgIds             []string `mapstructure:"msgIds"`
	}
)
