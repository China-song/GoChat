package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

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
	// TODO: types.GroupCreateReq add field IsVerify
	uid := ctxdata.GetUID(l.ctx)
	_, err = l.svcCtx.Social.GroupCreate(l.ctx, &socialclient.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil {
		fmt.Println("social api group CreateGroup: call rpc Social.GroupCreate err: ", err)
		return nil, err
	}
	return
}
