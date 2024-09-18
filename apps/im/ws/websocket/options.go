package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	pattern string

	maxConnectionIdle time.Duration

	ack          AckType
	ackTimeout   time.Duration
	sendErrCount int

	concurrency int // 群消息并发量级
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication: new(authentication),
		pattern:        "/ws",

		maxConnectionIdle: defaultMaxConnectionIdle,

		ackTimeout:   defaultAckTimeout,
		sendErrCount: defaultSendErrCount,

		concurrency: defaultConcurrency,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPattern(pattern string) ServerOptions {
	return func(opt *serverOption) {
		opt.pattern = pattern
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithServerSendErrCount(sendErrCount int) ServerOptions {
	return func(opt *serverOption) {
		opt.sendErrCount = sendErrCount
	}
}
