package logic

import (
	"GoChat/apps/im/immodels"
	"GoChat/pkg/constants"
	"context"
	"errors"

	"GoChat/apps/im/rpc/im"
	"GoChat/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupConversationLogic {
	return &CreateGroupConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateGroupConversation 创建群聊
func (l *CreateGroupConversationLogic) CreateGroupConversation(in *im.CreateGroupConversationReq) (*im.CreateGroupConversationResp, error) {
	resp := &im.CreateGroupConversationResp{}

	_, err := l.svcCtx.ConversationModel.FindByConversationId(l.ctx, in.GroupId)
	// 存在群ID
	if err == nil {
		return resp, nil
	}
	// 其他错误
	if !errors.Is(err, immodels.ErrNotFound) {
		l.Errorf("rpc im - CreateGroupConversation - ConversationModel.FindByConversationId err: %v", err)
		return resp, err
	}
	// 不存在，创建
	err = l.svcCtx.ConversationModel.Insert(l.ctx, &immodels.Conversation{
		ConversationId: in.GroupId,
		ChatType:       constants.GroupChatType,
	})
	if err != nil {
		l.Errorf("rpc im - CreateGroupConversation - ConversationModel.Insert err: %v", err)
		return resp, err
	}
	// 设置创建者用户会话列表
	_, err = NewSetUpUserConversationLogic(l.ctx, l.svcCtx).SetUpUserConversation(&im.SetUpUserConversationReq{
		SendId:   in.CreateId,
		RecvId:   in.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	if err != nil {
		l.Errorf("rpc im - CreateGroupConversation - SetUpUserConversation err: %v", err)
		return resp, err
	}

	return resp, nil
}
