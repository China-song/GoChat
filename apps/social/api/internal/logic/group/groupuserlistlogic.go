package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/apps/user/rpc/userclient"
	"context"
	"fmt"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群成员列表
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	// 获取群成员
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil {
		fmt.Println("social api group GroupUserList: call rpc Social.GroupUsers err: ", err)
		return nil, err
	}
	// 还需要获取用户的信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: uids,
	})
	if err != nil {
		fmt.Println("social api group GroupUserList: call rpc User.FindUser err: ", err)
		return nil, err
	}

	// 构造返回结果
	userRecords := make(map[string]*userclient.UserEntity, len(userList.User))
	for i := range userList.User {
		user := userList.User[i]
		userRecords[user.Id] = user
	}
	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {

		member := &types.GroupMembers{
			Id:        int64(v.Id),
			GroupId:   v.GroupId,
			UserId:    v.UserId,
			RoleLevel: int(v.RoleLevel),
		}
		if u, ok := userRecords[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}
		respList = append(respList, member)
	}

	return &types.GroupUserListResp{List: respList}, err
}
