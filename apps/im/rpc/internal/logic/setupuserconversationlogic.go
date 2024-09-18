package logic

import (
	"GoChat/apps/im/immodels"
	"GoChat/apps/im/rpc/im"
	"GoChat/apps/im/rpc/internal/svc"
	"GoChat/pkg/constants"
	"GoChat/pkg/wuid"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUpUserConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUpUserConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUpUserConversationLogic {
	return &SetUpUserConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUpUserConversation 建立会话: 群聊, 私聊
//
// 会话 user1 和 user2( 或 group2 )的会话
func (l *SetUpUserConversationLogic) SetUpUserConversation(in *im.SetUpUserConversationReq) (*im.SetUpUserConversationResp, error) {
	// 会话分为私聊 和 群聊会话
	switch constants.ChatType(in.ChatType) {
	case constants.SingleChatType:
		// 私聊会话 生成会话ID
		conversationId := wuid.CombineId(in.SendId, in.RecvId)
		// 判断会话是否已建立
		conversation, err := l.svcCtx.ConversationModel.FindByConversationId(l.ctx, conversationId)
		if err != nil {
			// 没有建立过会话，建立会话
			if errors.Is(err, immodels.ErrNotFound) {
				err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
					ConversationId: conversationId,
					ChatType:       constants.SingleChatType,
				})
				if err != nil {
					l.Errorf("rpc im - SetUpUserConversation - ConversationModel.Insert err: %v", err)
					return &im.SetUpUserConversationResp{}, err
				}
			} else {
				l.Errorf("rpc im - SetUpUserConversation - ConversationModel.FindOne err: %v", err)
				return &im.SetUpUserConversationResp{}, err
			}
		} else if conversation != nil {
			// 会话已经建立过，不需要重复建立
			l.Info("conversation has set up!")
			return &im.SetUpUserConversationResp{}, nil
		}
		// 建立用户会话记录:
		// sendUser: conversation
		// recvUser: conversation
		// 建立两者的会话
		err = l.setUpUserConversation(conversationId, in.SendId, in.RecvId, constants.SingleChatType, true)
		if err != nil {
			l.Errorf("rpc im - SetUpUserConversation - (Single)setUpUserConversation err: %v", err)
			return &im.SetUpUserConversationResp{}, err
		}
		// 接收者是被动与目标用户建立连接，因此理论上是不需要在会话列表里展示
		err = l.setUpUserConversation(conversationId, in.RecvId, in.SendId, constants.SingleChatType, false)
		if err != nil {
			l.Errorf("rpc im - SetUpUserConversation - (Single)setUpUserConversation err: %v", err)
			return &im.SetUpUserConversationResp{}, err
		}
	case constants.GroupChatType:
		// 接收者ID就是群会话ID
		err := l.setUpUserConversation(in.RecvId, in.SendId, in.RecvId, constants.GroupChatType, true)
		if err != nil {
			l.Errorf("rpc im - SetUpUserConversation - (Group)setUpUserConversation err: %v", err)
			return &im.SetUpUserConversationResp{}, err
		}
	}
	return &im.SetUpUserConversationResp{}, nil
}

func (l *SetUpUserConversationLogic) setUpUserConversation(conversationId, userId, recvId string, chatType constants.ChatType, isShow bool) error {
	// 用户的会话列表
	conversations, err := l.svcCtx.ConversationsModel.FindByUserId(l.ctx, userId)
	if err != nil {
		if errors.Is(err, immodels.ErrNotFound) {
			// 为空，创建新会话列表
			conversations = &immodels.Conversations{
				UserId:           userId,
				ConversationList: make(map[string]*immodels.Conversation),
			}
		} else {
			l.Errorf("rpc im - SetUpUserConversation - ConversationsModel.FindByUserId err: %v", err)
			return err
		}
	}
	// 根据会话ID判断是否有过会话
	if _, ok := conversations.ConversationList[conversationId]; ok {
		return nil
	}
	// 添加会话记录
	conversations.ConversationList[conversationId] = &immodels.Conversation{
		ConversationId: conversationId,
		ChatType:       chatType,
		IsShow:         isShow,
	}
	// 更新
	_, err = l.svcCtx.ConversationsModel.UpdateOrInsert(l.ctx, conversations)
	if err != nil {
		l.Errorf("rpc im - SetUpUserConversation - ConversationsModel.UpdateOrInsert err: %v", err)
		return err
	}
	return nil
}
