package logic

import (
	"GoChat/apps/im/rpc/imclient"
	"context"
	"github.com/jinzhu/copier"

	"GoChat/apps/im/api/internal/svc"
	"GoChat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 根据用户获取聊天记录
func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogLogic) GetChatLog(req *types.ChatLogReq) (resp *types.ChatLogResp, err error) {
	// todo: 权限未检测 如该用户是否有相应的会话ID
	data, err := l.svcCtx.Im.GetChatLog(l.ctx, &imclient.GetChatLogReq{
		ConversationId: req.ConversationId,
		StartSendTime:  req.StartSendTime,
		EndSendTime:    req.EndSendTime,
		Count:          req.Count,
	})
	if err != nil {
		l.Errorf("api im - GetChatLog - call rpc Im.GetChatLog err: %v", err)
		return nil, err
	}

	var res types.ChatLogResp
	err = copier.Copy(&res, &data)
	if err != nil {
		l.Errorf("api im - GetChatLog - copier.Copy err: %v", err)
		return nil, err
	}

	return &res, err
}
