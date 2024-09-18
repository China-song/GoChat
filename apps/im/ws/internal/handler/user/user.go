package user

import (
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
)

// OnLine
// todo: ???
func OnLine(svcCtx *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocket.NewMessage(u[0], uids), conn)
		srv.Info("err: ", err)
	}
}
