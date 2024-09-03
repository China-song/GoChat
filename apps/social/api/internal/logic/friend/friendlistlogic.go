package friend

import (
	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"
	"GoChat/apps/social/rpc/social"
	"GoChat/apps/user/rpc/userclient"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// 先获取自己的id
	uid := ctxdata.GetUID(l.ctx)
	friends, err := l.svcCtx.Social.FriendList(l.ctx, &social.FriendListReq{
		UserId: uid,
	})

	if err != nil {
		fmt.Println("social api FriendList err: ", err)
		return nil, err
	}

	if len(friends.List) == 0 {
		fmt.Println("social api FriendList: no friends")
		return &types.FriendListResp{}, nil
	}

	// 如果有好友，根据好友 id 获取好友信息
	friendUIDs := make([]string, 0, len(friends.List))
	for _, friend := range friends.List {
		friendUIDs = append(friendUIDs, friend.FriendUid)
	}

	// 根据 uids 查询用户信息
	users, err := l.svcCtx.User.FindUser(l.ctx, &userclient.FindUserReq{
		Ids: friendUIDs,
	})
	if err != nil {
		fmt.Println("social api FriendList: call user rpc FindUser err: ", err)
		return nil, err
	}

	// UID - UserEntity
	userRecords := make(map[string]*userclient.UserEntity, len(users.User))
	for i := range users.User {
		user := users.User[i]
		userRecords[user.Id] = user
	}

	respList := make([]*types.Friends, 0, len(friends.List))
	for _, v := range friends.List {
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: v.FriendUid,
		}
		if user, ok := userRecords[v.FriendUid]; ok {
			friend.Nickname = user.Nickname
			friend.Avatar = user.Avatar
		}
		respList = append(respList, friend)
	}
	return &types.FriendListResp{
		List: respList,
	}, nil
}
