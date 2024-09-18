package main

import (
	"GoChat/apps/im/ws/internal/config"
	"GoChat/apps/im/ws/internal/handler"
	"GoChat/apps/im/ws/internal/svc"
	"GoChat/apps/im/ws/websocket"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"time"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	ctx := svc.NewServiceContext(c)

	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerMaxConnectionIdle(10*time.Minute), // 10 min
		websocket.WithServerAck(websocket.NoAck),
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Printf("Starting im server at %s...\n", c.ListenOn)
	srv.Start()
}
