package rpcserver

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	logx.WithContext(ctx).Error("【RPC SRV ERR】 %v", err)

	return resp, err
}
