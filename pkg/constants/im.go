package constants

type MsgType int

const (
	TextMType MsgType = iota
	// todo 视频 图片 语音 文件
)

type ChatType int

const (
	GroupChatType ChatType = iota + 1
	SingleChatType
)

type ContentType int

const (
	ContentChatMsgType ContentType = iota
	ContentMarkReadType
)
