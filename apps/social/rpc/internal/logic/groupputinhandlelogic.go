package logic

import (
	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"
	"GoChat/apps/social/socialmodels"
	"GoChat/pkg/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrGroupReqPassed  = errors.New("群申请已通过")
	ErrGroupReqRefused = errors.New("群申请已拒绝")
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	// todo: 需验证handleUid 是否为 该群的管理员？
	// User信息？
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		fmt.Println("social rpc GroupPutInHandle: call GroupRequestsModel.FindOne err: ", err)
		return &social.GroupPutInHandleResp{}, err
	}
	// 根据处理结果进行判断
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandlerResult:
		// 若之前已经通过，则返回错误
		return &social.GroupPutInHandleResp{}, ErrGroupReqPassed
	case constants.RefuseHandlerResult:
		// 若之前已经拒绝，则返回错误
		return &social.GroupPutInHandleResp{}, ErrGroupReqRefused
	}

	// 未处理
	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.GroupRequestsModel.UpdateWithSession(l.ctx, session, groupReq)
		if err != nil {
			fmt.Println("social rpc GroupPutInHandle: call GroupRequestsModel.UpdateWithSession err: ", err)
			return err
		}
		// 若处理结果不是通过，则直接返回
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandlerResult {
			return nil
		}
		// 通过申请 需要insert group_members
		groupMember := &socialmodels.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int64(constants.MemberGroupRoleLevel),
			JoinTime:    sql.NullTime{Time: time.Now(), Valid: true},
			JoinSource:  groupReq.JoinSource,
			InviterUid:  groupReq.InviterUserId,
			OperatorUid: sql.NullString{String: in.HandleUid, Valid: true},
		}
		_, err = l.svcCtx.GroupMembersModel.InsertWithSession(l.ctx, session, groupMember)
		if err != nil {
			fmt.Println("social rpc GroupPutInHandle: call GroupMembersModel.InsertWithSession err: ", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("social rpc GroupPutInHandle: call GroupRequestsModel.Trans err: ", err)
		return &social.GroupPutInHandleResp{}, err
	}
	return &social.GroupPutInHandleResp{
		GroupId: groupReq.GroupId,
	}, nil
}
