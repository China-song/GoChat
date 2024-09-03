package friend

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"
	"github.com/jinzhu/copier"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	UID := ctxdata.GetUID(l.ctx)
	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &socialclient.FriendPutInListReq{
		UserId: UID,
	})
	if err != nil {
		fmt.Println("social api FriendPutInList: call social rpc FriendPutInList err: ", err)
		return nil, err
	}

	var respList []*types.FriendRequests
	copier.Copy(&respList, list.List)

	return &types.FriendPutInListResp{List: respList}, nil
}
