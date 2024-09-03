package group

import (
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/ctxdata"
	"context"
	"fmt"

	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群
func NewGroupPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInLogic {
	return &GroupPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInLogic) GroupPutIn(req *types.GroupPutInRep) (resp *types.GroupPutInResp, err error) {
	// todo: types.GroupPutInRep add field 邀请者 被邀请者
	// 若是邀请 ReqId = 被邀请者   InviterUid = uid
	// 不是邀请 ReqId = uid
	uid := ctxdata.GetUID(l.ctx)
	_, err = l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})
	if err != nil {
		fmt.Println("social api group GroupPutIn: call rpc Social.GroupPutin err: ", err)
		return nil, err
	}
	return
}
