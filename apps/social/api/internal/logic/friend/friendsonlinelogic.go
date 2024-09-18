package friend

import (
	"GoChat/apps/social/rpc/social"
	"GoChat/pkg/constants"
	"GoChat/pkg/ctxdata"
	"context"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线情况
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendsOnlineResp, err error) {
	// 获取当前用户ID
	uid := ctxdata.GetUID(l.ctx)
	// 获取当前用户所有好友列表
	friendList, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})
	if err != nil {
		l.Errorf("api social - FriendsOnline - call rpc Social.FriendList err: %v", err)
		return &types.FriendsOnlineResp{}, err
	}
	if len(friendList.List) == 0 {
		return &types.FriendsOnlineResp{}, nil
	}
	// 查询缓存中在线的用户
	uids := make([]string, 0, len(friendList.List))
	for _, friend := range friendList.List {
		uids = append(uids, friend.UserId)
	}
	onlines, err := l.svcCtx.Redis.Hgetall(constants.RedisOnlineUser)
	if err != nil {
		l.Errorf("api social - FriendsOnline - Redis.Hgetall err: %v", err)
		return &types.FriendsOnlineResp{}, err
	}
	resOnlineList := make(map[string]bool, len(uids))
	for _, uid := range uids {
		if _, ok := onlines[uid]; ok {
			resOnlineList[uid] = true
		} else {
			resOnlineList[uid] = false
		}
	}

	return &types.FriendsOnlineResp{
		OnlineList: resOnlineList,
	}, nil
}
