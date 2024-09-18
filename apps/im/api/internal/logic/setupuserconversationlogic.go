package logic

import (
	"GoChat/apps/im/rpc/im"
	"GoChat/pkg/ctxdata"
	"context"

	"GoChat/apps/im/api/internal/svc"
	"GoChat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 建立会话
func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetUpUserConversationLogic) SetUpUserConversation(req *types.SetUpUserConversationReq) (resp *types.SetUpUserConversationResp, err error) {
	uid := ctxdata.GetUID(l.ctx)
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &im.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   req.RecvId,
		ChatType: req.ChatType,
	})
	if err != nil {
		l.Errorf("api im - SetUpUserConversation - call rpc Im.SetUpUserConversation err: %v", err)
		return nil, err
	}

	return
}
