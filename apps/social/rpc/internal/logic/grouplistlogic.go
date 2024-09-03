package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupListLogic) GroupList(in *social.GroupListReq) (*social.GroupListResp, error) {
	// 根据用户id查询其所在的所有群
	userGroups, err := l.svcCtx.GroupMembersModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		fmt.Println("social rpc GroupList: GroupMembersModel.ListByUserId err: ", err)
		return &social.GroupListResp{}, err
	}

	// 根据group ID查询groups表
	groupIds := make([]string, 0, len(userGroups))
	for _, v := range userGroups {
		groupIds = append(groupIds, v.GroupId)
	}

	// 根据groupIds查询groups表 得到相应的groups
	groups, err := l.svcCtx.GroupsModel.ListByGroupIds(l.ctx, groupIds)
	if err != nil {
		fmt.Println("social rpc GroupList: GroupsModel.ListByGroupIds err: ", err)
		return &social.GroupListResp{}, err
	}

	var resp []*social.Groups
	err = copier.Copy(&resp, &groups)
	if err != nil {
		fmt.Println("social rpc GroupList: copier.Copy err: ", err)
		return &social.GroupListResp{}, err
	}
	return &social.GroupListResp{
		List: resp,
	}, nil
}
