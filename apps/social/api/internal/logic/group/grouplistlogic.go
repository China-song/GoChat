package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户申群列表
func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupListLogic) GroupList(req *types.GroupListRep) (resp *types.GroupListResp, err error) {
	uid := ctxdata.GetUID(l.ctx)
	list, err := l.svcCtx.Social.GroupList(l.ctx, &socialclient.GroupListReq{
		UserId: uid,
	})
	if err != nil {
		fmt.Println("social api group GroupList: call rpc Social.GroupList err: ", err)
		return nil, err
	}

	var respList []*types.Groups
	// TODO: types.Groups add field CreatorUid
	err = copier.Copy(&respList, &list.List)
	if err != nil {
		fmt.Println("social api group GroupList: call copier.Copy err: ", err)
		return nil, err
	}

	return &types.GroupListResp{List: respList}, nil
}
