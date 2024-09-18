package websocket

import "time"

type FrameType uint8

const (
	FrameData FrameType = 0x0 // 数据帧
	FramePing FrameType = 0x1 // Ping 帧

	FrameAck   FrameType = 0x2 // Ack 帧
	FrameNoAck FrameType = 0x3 // 无 Ack 帧

	FrameErr FrameType = 0x9
)

type Message struct {
	FrameType `json:"frameType"`
	Id        string    `json:"id"`       // 消息 ID
	AckSeq    int       `json:"ackSeq"`   // Ack序列号
	AckTime   time.Time `json:"ackTime"`  // 确认时间
	ErrCount  int       `json:"errCount"` // 错误计数
	Method    string    `json:"method"`
	FromId    string    `json:"fromId"` // 消息来源
	Data      any       `json:"data"`
}

func NewMessage(fromId string, data any) *Message {
	return &Message{
		FrameType: FrameData,

		FromId: fromId,
		Data:   data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
