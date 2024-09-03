package logic

import (
	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendList 获取好友列表
func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	friendsList, err := l.svcCtx.FriendsModel.LiseByUserId(l.ctx, in.UserId)
	if err != nil {
		return &social.FriendListResp{}, err
	}

	var respList []*social.Friends
	err = copier.Copy(&respList, &friendsList)
	if err != nil {
		fmt.Println("social rpc GroupList: copier.Copy err: ", err)
		return &social.FriendListResp{}, err
	}

	return &social.FriendListResp{
		List: respList,
	}, nil
}
