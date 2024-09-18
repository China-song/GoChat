package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	idleMu sync.Mutex

	Uid string

	*websocket.Conn
	s *Server

	idle              time.Time
	maxConnectionIdle time.Duration

	done chan struct{}

	messageMu      sync.Mutex
	readMessage    []*Message // 消息队列
	readMessageSeq map[string]*Message
	message        chan *Message
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("Error upgrading connection: %s", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(),
		maxConnectionIdle: s.opt.maxConnectionIdle,
		done:              make(chan struct{}),
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 1),
	}
	go conn.keepalive()
	return conn
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() { // 连接非空闲状态。
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 已超时 应关闭连接
				c.s.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{}
	return
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	err := c.Conn.WriteMessage(messageType, data)
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Now()
	return err
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	// 读队列中，判断之前是否在队列中存过消息
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 该消息已经有Ack的确认过程
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}
		// 要求 msg.AckSeq > m.AckSeq
		if msg.AckSeq <= m.AckSeq {
			// 没有进行ack的确认, 重复
			return
		}
		// 更新最新的消息
		c.readMessageSeq[msg.Id] = msg // 指向了新消息
		return
	}
	// 还没有进行Ack确认，避免客户端重复发送多余的Ack消息
	if msg.FrameType == FrameAck {
		return
	}
	// 记录消息
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
}
