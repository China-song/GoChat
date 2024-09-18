package logic

import (
	"GoChat/apps/im/rpc/imclient"
	"GoChat/pkg/ctxdata"
	"context"
	"github.com/jinzhu/copier"

	"GoChat/apps/im/api/internal/svc"
	"GoChat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会话
func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsLogic) GetConversations(req *types.GetConversationsReq) (resp *types.GetConversationsResp, err error) {
	uid := ctxdata.GetUID(l.ctx)
	data, err := l.svcCtx.Im.GetConversations(l.ctx, &imclient.GetConversationsReq{
		UserId: uid,
	})
	if err != nil {
		l.Errorf("api im - GetConversations - call rpc Im.GetConversations err: %v", err)
		return nil, err
	}

	var res types.GetConversationsResp
	// todo: 未copy完备
	err = copier.Copy(&res, &data)
	if err != nil {
		l.Errorf("api im - GetConversations - copier.Copy err: %v", err)
		return nil, err
	}

	return &res, err
}
