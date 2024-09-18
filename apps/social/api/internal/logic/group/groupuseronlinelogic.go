package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/constants"
	"context"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群在线用户
func NewGroupUserOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserOnlineLogic {
	return &GroupUserOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserOnlineLogic) GroupUserOnline(req *types.GroupUserOnlineReq) (resp *types.GroupUserOnlineResp, err error) {
	// 获取当前群的所有成员
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
		GroupId: req.GroupId,
	})
	if err != nil || len(groupUsers.List) == 0 {
		l.Errorf("api social - GroupUserOnline - call rpc Social.GroupUsers err: %v", err)
		return &types.GroupUserOnlineResp{}, err
	}
	// 查询，缓存中查询在线的用户
	uids := make([]string, 0, len(groupUsers.List))
	for _, groupUser := range groupUsers.List {
		uids = append(uids, groupUser.UserId)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(constants.RedisOnlineUser)
	if err != nil {
		l.Errorf("api social - GroupUserOnline - Redis.Hgetall err: %v", err)
		return &types.GroupUserOnlineResp{}, err
	}
	resOnLineList := make(map[string]bool, len(uids))
	for _, s := range uids {
		if _, ok := onlines[s]; ok {
			resOnLineList[s] = true
		} else {
			resOnLineList[s] = false
		}
	}

	return &types.GroupUserOnlineResp{
		OnlineList: resOnLineList,
	}, nil
}
