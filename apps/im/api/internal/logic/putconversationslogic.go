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

type PutConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新会话
func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PutConversationsLogic) PutConversations(req *types.PutConversationsReq) (resp *types.PutConversationsResp, err error) {
	uid := ctxdata.GetUID(l.ctx)
	var conversationList map[string]*imclient.Conversation
	// todo: copy不完整
	err = copier.Copy(&conversationList, req.ConversationList)
	if err != nil {
		l.Errorf("api im - PutConversations - copier.Copy err: %v", err)
		return nil, err
	}

	_, err = l.svcCtx.Im.PutConversations(l.ctx, &imclient.PutConversationsReq{
		UserId:           uid,
		ConversationList: conversationList,
	})
	if err != nil {
		l.Errorf("api im - PutConversations - call rpc Im.PutConversations err: %v", err)
		return nil, err
	}
	return
}
