package logic

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/rpc/im"
	"GoChat/apps/im/rpc/internal/svc"
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetConversations 获取会话
//
// 首先根据用户的所有会话（历史会话）
// 根据会话表的当前信息 更新 用户的会话未读消息量
func (l *GetConversationsLogic) GetConversations(in *im.GetConversationsReq) (*im.GetConversationsResp, error) {
	// 查询用户的会话列表 map{conversationId: conversation}
	data, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, immodels.ErrNotFound) {
			l.Infof("用户(%s)无会话", in.UserId)
			return &im.GetConversationsResp{}, nil
		}
		l.Errorf("rpc im - GetConversations - ConversationsModel.FindByUserId err: %v", err)
		return &im.GetConversationsResp{}, err
	}
	var res im.GetConversationsResp
	err = copier.Copy(&res, &data)
	if err != nil {
		l.Errorf("rpc im - GetConversations - copier.Copy err: %v", err)
		return &im.GetConversationsResp{}, err
	}

	// 根据会话id查询具体会话
	ids := make([]string, 0, len(data.ConversationList))
	for _, conversation := range data.ConversationList {
		ids = append(ids, conversation.ConversationId)
	}
	conversations, err := l.svcCtx.ConversationModel.ListByConversationIds(l.ctx, ids)
	if err != nil {
		l.Errorf("rpc im - GetConversations - ConversationModel.ListByConversationIds err: %v", err)
		return &im.GetConversationsResp{}, err
	}

	// 计算是否存在未读消息
	for _, conversation := range conversations {
		if _, ok := res.ConversationList[conversation.ConversationId]; !ok {
			continue
		}
		// 用户读取的消息量
		total := res.ConversationList[conversation.ConversationId].Total
		if total < int32(conversation.Total) {
			// 有新的消息
			res.ConversationList[conversation.ConversationId].Total = int32(conversation.Total)
			// 有多少是未读
			res.ConversationList[conversation.ConversationId].ToRead = int32(conversation.Total) - total
			// 更改当前会话为显示状态
			res.ConversationList[conversation.ConversationId].IsShow = true
		}
	}

	return &res, nil
}
