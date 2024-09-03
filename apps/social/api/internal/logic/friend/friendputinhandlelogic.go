package friend

import (
	"GoChat/apps/social/rpc/social"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请处理
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// 先获取自己的id
	uid := ctxdata.GetUID(l.ctx)
	_, err = l.svcCtx.Social.FriendPutInHandle(l.ctx, &social.FriendPutInHandleReq{
		FriendReqId:  req.FriendReqId,
		UserId:       uid,
		HandleResult: req.HandleResult,
	})

	if err != nil {
		fmt.Println("social api FriendPutInHandle err: ", err)
		return nil, err
	}

	return
}
