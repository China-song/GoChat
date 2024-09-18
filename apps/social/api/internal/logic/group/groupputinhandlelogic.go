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

type GroupPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请进群处理
func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(req *types.GroupPutInHandleRep) (resp *types.GroupPutInHandleResp, err error) {
	uid := ctxdata.GetUID(l.ctx)
	res, err := l.svcCtx.Social.GroupPutInHandle(l.ctx, &socialclient.GroupPutInHandleReq{
		GroupReqId:   req.GroupReqId,
		GroupId:      req.GroupId,
		HandleUid:    uid,
		HandleResult: req.HandleResult,
	})
	if err != nil {
		l.Errorf("social api - (group) GroupPutInHandle - rpc Social.GroupPutInHandle err: ", err)
		return nil, err
	}

	if constants.HandlerResult(req.HandleResult) != constants.PassHandlerResult {
		return
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
