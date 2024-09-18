package group

import (
	"GoChat/apps/im/rpc/imclient"
	"GoChat/apps/social/api/internal/svc"
	"GoChat/apps/social/api/internal/types"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/pkg/constants"
	"GoChat/pkg/ctxdata"
	"context"

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
	// 若是邀请 ReqId = 被邀请者   InviterUid = uid  JoinSource = InviteGroupJoinSource
	// 不是邀请 ReqId = uid JoinSource = PutInGroupJoinSource
	uid := ctxdata.GetUID(l.ctx)
	res, err := l.svcCtx.Social.GroupPutin(l.ctx, &socialclient.GroupPutinReq{
		GroupId:    req.GroupId,
		ReqId:      uid,
		ReqMsg:     req.ReqMsg,
		ReqTime:    req.ReqTime,
		JoinSource: int32(req.JoinSource),
	})
	if err != nil {
		l.Errorf("social api - (group) GroupPutIn - rpc Social.GroupPutin err: ", err)
		return nil, err
	}

	if res.GroupId == "" {
		return nil, err
	}
	// 建立会话
	_, err = l.svcCtx.Im.SetUpUserConversation(l.ctx, &imclient.SetUpUserConversationReq{
		SendId:   uid,
		RecvId:   res.GroupId,
		ChatType: int32(constants.GroupChatType),
	})
	if err != nil {
		l.Errorf("social api - (group) GroupPutIn - rpc Im.SetUpUserConversation err: ", err)
		return nil, err
	}
	return
}
