package logic

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/rpc/im"
	"GoChat/apps/im/rpc/internal/svc"
	"GoChat/pkg/constants"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type PutConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPutConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PutConversationsLogic {
	return &PutConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话
func (l *PutConversationsLogic) PutConversations(in *im.PutConversationsReq) (*im.PutConversationsResp, error) {
	// 查询用户的会话列表
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("rpc im - PutConversations - ConversationsModel.FindByUserId err: %v", err)
		return &im.PutConversationsResp{}, err
	}

	// 存在会话，会话列表为空
	if data.ConversationList == nil {
		data.ConversationList = make(map[string]*immodels.Conversation)
	}

	for s, conversation := range in.ConversationList {
		// 获取用户原本读取的会话消息量
		var oldTotal int
		if data.ConversationList[s] != nil {
			oldTotal = data.ConversationList[s].Total
		}
		// 设置新结果
		data.ConversationList[s] = &immodels.Conversation{
			ConversationId: conversation.ConversationId,
			ChatType:       constants.ChatType(conversation.ChatType),
			IsShow:         conversation.IsShow,
			Total:          int(conversation.Read) + oldTotal, // 已读记录量 + 原本读取的会话消息量
			Seq:            conversation.Seq,
		}
	}

	_, err = l.svcCtx.ConversationsModel.Update(l.ctx, data)
	if err != nil {
		l.Errorf("rpc im - PutConversations - ConversationsModel.Update err: %v", err)
		return &im.PutConversationsResp{}, err
	}

	return &im.PutConversationsResp{}, nil
}
