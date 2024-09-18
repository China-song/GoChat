package group

import (
	"GoChat/apps/im/rpc/imclient"
	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/ctxdata"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	// 群名 群图标
	// TODO: types.GroupCreateReq add field IsVerify 创建群时指定是否要群验证
	uid := ctxdata.GetUID(l.ctx)
	res, err := l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil {
		l.Errorf("social api - (group) CreateGroup - rpc Social.GroupCreate err: ", err)
		return nil, err
	}

	if res.Id == "" {
		l.Errorf("social api - (group) CreateGroup - rpc Social.GroupCreate err: ", err)
		return nil, err
	}
	// 建立会话
	_, err = l.svcCtx.Im.CreateGroupConversation(l.ctx, &imclient.CreateGroupConversationReq{
		GroupId:  res.Id,
		CreateId: uid,
	})
	if err != nil {
		l.Errorf("social api - (group) CreateGroup - rpc Im.CreateGroupConversation err: ", err)
		return nil, err
	}
	return
}
