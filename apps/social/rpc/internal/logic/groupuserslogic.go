package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUsersLogic {
	return &GroupUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupUsers returns members of group(GroupId)
func (l *GroupUsersLogic) GroupUsers(in *social.GroupUsersReq) (*social.GroupUsersResp, error) {
	// query group_members model by group_id
	groupMembers, err := l.svcCtx.GroupMembersModel.ListByGroupId(l.ctx, in.GroupId)
	if err != nil {
		fmt.Println("social rpc GroupUsers: call GroupMembersModel.ListByGroupId err: ", err)
		return &social.GroupUsersResp{}, nil
	}

	var resp []*social.GroupMembers
	err = copier.Copy(&resp, &groupMembers)
	if err != nil {
		fmt.Println("social rpc GroupUsers: call copier.Copy err: ", err)
		return &social.GroupUsersResp{}, nil
	}
	return &social.GroupUsersResp{
		List: resp,
	}, nil
}
