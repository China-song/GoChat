package websocket

type AckType int

const (
	NoAck AckType = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case NoAck:
		return "NoAck"
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	default:
		panic("Unknown AckType")
	}
}
