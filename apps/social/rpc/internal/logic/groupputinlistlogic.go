package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinListLogic {
	return &GroupPutinListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutinList returns groupRequests which are not be handled
func (l *GroupPutinListLogic) GroupPutinList(in *social.GroupPutinListReq) (*social.GroupPutinListResp, error) {
	groupReqs, err := l.svcCtx.GroupRequestsModel.ListNoHandler(l.ctx, in.GroupId)
	if err != nil {
		fmt.Println("social rpc GroupPutinList: call GroupRequestsModel.ListNoHandler err: ", err)
		return &social.GroupPutinListResp{}, err
	}
	var resp []*social.GroupRequests
	err = copier.Copy(&resp, &groupReqs)
	if err != nil {
		fmt.Println("social rpc GroupPutinList: call copier.Copy err: ", err)
		return &social.GroupPutinListResp{}, err
	}
	return &social.GroupPutinListResp{
		List: resp,
	}, nil
}
