package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	routes   map[string]HandlerFunc
	addr     string
	pattern  string
	opt      *serverOption
	upgrader websocket.Upgrader
	logx.Logger

	connToUser map[*Conn]string
	userToConn map[string]*Conn
	sync.RWMutex

	authentication Authentication

	*threading.TaskRunner
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		routes:   make(map[string]HandlerFunc),
		addr:     addr,
		pattern:  opt.pattern,
		opt:      &opt,
		upgrader: websocket.Upgrader{},
		Logger:   logx.WithContext(context.Background()),

		connToUser: make(map[*Conn]string),
		userToConn: make(map[string]*Conn),

		authentication: opt.Authentication,

		TaskRunner: threading.NewTaskRunner(opt.concurrency),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("server handler ws recover err: %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}

	// 对连接鉴权
	if !s.authentication.Authenticate(w, r) {
		s.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不具备访问权限")}, conn)
		conn.Close()
		return
	}

	// 记录连接
	s.addConn(conn, r)

	// 处理连接
	go s.handleConn(conn)
}

func (s *Server) handleConn(conn *Conn) {
	go s.handleWrite(conn)
	if s.isAck(nil) {
		go s.readAck(conn)
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket conn read message err: %v", err)
			s.Close(conn)
			return
		}

		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("json unmarshal err: %v, msg: %s", err, string(msg))
			s.Close(conn)
			return
		}

		// 需要ACK
		// 将消息放到消息队列中
		if s.isAck(&message) {
			s.Infof("save msg to readMessage queue. msg: %v", message)
			conn.appendMsgMq(&message)
		} else {
			// 不需要ACK
			// 直接进行message的处理
			conn.message <- &message
		}
	}
}

// 进行消息的处理
func (s *Server) handleWrite(conn *Conn) {
	for {
		select {
		// todo ??? 是否阻塞在此
		case <-conn.done:
			// 连接关闭
			return
		case message := <-conn.message:
			// 有消息 处理
			switch message.FrameType {
			case FramePing:
				s.Send(&Message{
					FrameType: FramePing,
				}, conn)
			case FrameData:
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					err := s.Send(&Message{
						FrameType: FrameData,
						Data:      fmt.Sprintf("不存在执行的方法 %v 请检查", message.Method),
					}, conn)
					if err != nil {
						s.Errorf("websocket conn write message err: %v", err)
						s.Close(conn)
						return
					}
				}
			}

			// todo 清除消息确认
			// 清除消息确认
			if s.isAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

// 判断websocket server 是否开启 ACK模式
// 以及消息是否需要ACK
func (s *Server) isAck(message *Message) bool {
	return s.opt.ack != NoAck && (message == nil || message.FrameType != FrameNoAck)
}

// 对消息队列中的消息返回ACK
// ACK消息:
/*
	Message {
		FrameType: FrameAck,
		Id:        msg.Id,
		AckSeq:    msg.AckSeq + 1,
	}

*/
func (s *Server) readAck(conn *Conn) {
	// send ACK message to client
	send := func(msg *Message, conn *Conn) error {
		err := s.Send(msg, conn)
		if err != nil {
			s.Errorf("send ACK message err: %v, message: %v", err, msg)
			return err
		}
		return nil
	}
	for {
		select {
		case <-conn.done:
			s.Infof("conn close when readAck. (conn of %v)", conn.Uid)
			return
		default:
		}
		// 从队列中读取消息 返回ACK
		conn.messageMu.Lock()
		// 当前队列中没有消息
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			time.Sleep(100 * time.Microsecond)
			continue
		}

		msg := conn.readMessage[0]
		switch s.opt.ack {
		case OnlyAck:
			// 仅需确认一次
			if err := send(&Message{
				FrameType: FrameAck,
				Id:        msg.Id,
				AckSeq:    msg.AckSeq + 1,
				//AckTime:   time.Time{},
			}, conn); err != nil {
				conn.messageMu.Unlock()
				continue
			}
			// 发送完ACK 处理消息
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- msg
			s.Infof("message ack (OnlyAck) send success, msg ID: %v", msg.Id)
		case RigorAck:
			// 两次确认
			// 需判断当前是否已完成第一次确认
			if msg.AckSeq == 0 {
				// 还未向客户端发送ACK
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].AckTime = time.Now()
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        msg.Id,
					AckSeq:    msg.AckSeq,
				}, conn); err != nil {
					// todo: 发送失败 是否应撤回对消息的修改
					conn.messageMu.Unlock()
					continue
				}
				conn.messageMu.Unlock()
				s.Infof("message ack RigorAck send mid: %v, seq: %v , time: %v", msg.Id, msg.AckSeq, msg.AckTime)
				continue
			}
			// 验证
			// 1. 客户端返回结果，再一次确认
			// 获取之前记录的Ack信息，得到客户端的序号
			msgSeq := conn.readMessageSeq[msg.Id]
			// todo ??? 同一个消息
			s.Infof("msg == msgSeq? %v", msg == msgSeq)
			if msgSeq.AckSeq > msg.AckSeq {
				// 客户端进行了确认
				// 删除消息
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- msg
				s.Infof("message ack RigorAck success mid: %v", msg.Id)
				continue
			}

			// 2. 客户端没有确认，考虑是否超过了ack的确认时间
			val := s.opt.ackTimeout - time.Since(msg.AckTime)
			if !msg.AckTime.IsZero() && val <= 0 {
				// 2.1 超过结束确认
				// todo ??? 超时 为什么要删除
				s.Infof("client ack timeout! message: %v, last ack time: %v - current time: %v", msg.Id, msg.AckTime, time.Now())
				// 删除消息序号
				delete(conn.readMessageSeq, msg.Id)
				// 删除消息
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}

			// 2.2 未超时，重新发送
			conn.messageMu.Unlock()
			if val > 0 && val > 300*time.Microsecond {
				if err := send(&Message{
					FrameType: FrameAck,
					Id:        msg.Id,
					AckSeq:    msg.AckSeq,
				}, conn); err != nil {
					continue
				}
			}
			// 睡眠一定的时间
			//time.Sleep(300 * time.Microsecond)
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.pattern, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	logx.Close()
	fmt.Println("停止服务")
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

// addConn 存储连接对象
func (s *Server) addConn(conn *Conn, req *http.Request) {
	uid := s.authentication.UserId(req)
	conn.Uid = uid

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if c := s.userToConn[uid]; c != nil {
		c.Close()
	}

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// GetConn 根据uid获取conn连接对象
func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUsers(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取部分
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) Close(conn *Conn) {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		return
	}

	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
	conn.Close()
}

func (s *Server) Send(msg any, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) SendByUserIds(msg any, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}
