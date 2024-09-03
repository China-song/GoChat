package logic

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInListLogic) FriendPutInList(in *social.FriendPutInListReq) (*social.FriendPutInListResp, error) {
	friendReqList, err := l.svcCtx.FriendRequestsModel.ListNoHandler(l.ctx, in.UserId)
	if err != nil {
		fmt.Println("social rpc FriendPutInList err: ", err)
		return &social.FriendPutInListResp{}, err
	}

	var resp []*social.FriendRequests
	err = copier.Copy(&resp, &friendReqList)
	if err != nil {
		fmt.Println("social rpc GroupList: copier.Copy err: ", err)
		return &social.FriendPutInListResp{}, err
	}

	return &social.FriendPutInListResp{
		List: resp,
	}, nil
}
