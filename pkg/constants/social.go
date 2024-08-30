package constants

// HandlerResult 申请处理结果 1. 未处理，2. 通过， 3. 拒绝
type HandlerResult int

const (
	NoHandlerResult     HandlerResult = iota + 1 // 未处理
	PassHandlerResult                            // 通过
	RefuseHandlerResult                          // 拒绝
	CancelHandlerResult                          // 撤销
)
