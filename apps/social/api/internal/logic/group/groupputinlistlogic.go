package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群列表
func NewGroupPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInListLogic {
	return &GroupPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInListLogic) GroupPutInList(req *types.GroupPutInListRep) (resp *types.GroupPutInListResp, err error) {
	list, err := l.svcCtx.Social.GroupPutinList(l.ctx, &socialclient.GroupPutinListReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		fmt.Println("social api group GroupPutInList: call rpc Social.GroupPutinList err: ", err)
		return nil, err
	}

	var respList []*types.GroupRequests
	err = copier.Copy(&respList, &list.List)
	if err != nil {
		fmt.Println("social api group GroupPutInList: call copier.Copy err: ", err)
		return nil, err
	}

	return &types.GroupPutInListResp{List: respList}, nil
}
