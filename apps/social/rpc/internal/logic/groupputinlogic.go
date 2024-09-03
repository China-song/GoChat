package logic

import (
	"GoChat/apps/social/socialmodels"
	"GoChat/pkg/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {

	// 首先查询申请者是否已是该群的成员
	// select from group_member where group_id = and user_id =
	userGroupMember, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		fmt.Println("social rpc GroupPutin: call GroupMembersModel.FindByGroupIdAndUserId err: ", err)
		return &social.GroupPutinResp{}, nil
	}
	if userGroupMember != nil {
		// 已经是成员
		fmt.Println("social rpc GroupPutin: 已是群成员")
		return &social.GroupPutinResp{GroupId: in.GroupId}, nil
	}

	// 查询该用户是否已发过群申请
	// TODO: 发送过的群申请若被拒绝 则本次请求可以接受？
	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		fmt.Println("social rpc GroupPutin: call GroupRequestsModel.FindByGroupIdAndReqId err: ", err)
		return &social.GroupPutinResp{}, nil
	}
	if groupReq != nil {
		// 该用户已申请过
		return &social.GroupPutinResp{}, nil
	}

	// 未申请过

	// insert group_requests
	groupReq = &socialmodels.GroupRequests{
		ReqId:      in.ReqId,
		GroupId:    in.GroupId,
		ReqMsg:     sql.NullString{String: in.ReqMsg, Valid: true},
		ReqTime:    sql.NullTime{Time: time.Unix(in.ReqTime, 0), Valid: true},
		JoinSource: sql.NullInt64{Int64: int64(in.JoinSource), Valid: true},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	}

	groupMem := &socialmodels.GroupMembers{
		GroupId:   in.GroupId,
		UserId:    in.ReqId,
		RoleLevel: int64(constants.MemberGroupRoleLevel),
		JoinTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		JoinSource: sql.NullInt64{Int64: int64(in.JoinSource), Valid: true},
		InviterUid: sql.NullString{String: in.InviterUid, Valid: true},
	}

	// 查看群申请是否需要验证
	groupInfo, err := l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		fmt.Println("social rpc GroupPutin: call GroupsModel.FindOne err: ", err)
		return &social.GroupPutinResp{}, nil
	}
	// 若该群不需要被验证 则直接通过请求
	// insert group_members
	// insert group_req handle_result = 2
	if !groupInfo.IsVerify {
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		return l.passGroupPutin(groupReq, groupMem)
	}

	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		return l.createGroupReq(groupReq)
	}

	// 邀请加入群 查询邀请者是否为群的管理者
	inviterGroupMember, err := l.svcCtx.GroupMembersModel.FindByGroupIdAndUserId(l.ctx, in.GroupId, in.InviterUid)
	if err != nil {
		fmt.Println("social rpc GroupPutin: call GroupMembersModel.FindByGroupIdAndUserId err: ", err)
		return &social.GroupPutinResp{}, nil
	}
	// 若是管理者邀请 则直接通过
	inviterGroupRoleLevel := constants.GroupRoleLevel(inviterGroupMember.RoleLevel)
	if inviterGroupRoleLevel == constants.CreatorGroupRoleLevel || inviterGroupRoleLevel == constants.ManagerGroupRoleLevel {
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandlerResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}

		groupMem.OperatorUid = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.passGroupPutin(groupReq, groupMem)
	}
	// 普通成员邀请
	return l.createGroupReq(groupReq)
}

// 直接通过申请 需要insert group_members 和 group_requests
func (l *GroupPutinLogic) passGroupPutin(groupReq *socialmodels.GroupRequests, groupMember *socialmodels.GroupMembers) (*social.GroupPutinResp, error) {
	err := l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.GroupRequestsModel.InsertWithSession(l.ctx, session, groupReq)
		if err != nil {
			fmt.Println("social rpc GroupPutin: call GroupRequestsModel.InsertWithSession err: ", err)
			return err
		}
		_, err = l.svcCtx.GroupMembersModel.InsertWithSession(l.ctx, session, groupMember)
		if err != nil {
			fmt.Println("social rpc GroupPutin: call GroupMembersModel.InsertWithSession err: ", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("social rpc GroupPutin: call passGroupPutin err: ", err)
		return &social.GroupPutinResp{}, nil
	}
	return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
}

func (l *GroupPutinLogic) createGroupReq(groupReq *socialmodels.GroupRequests) (*social.GroupPutinResp, error) {
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		fmt.Println("social rpc GroupPutin: call GroupRequestsModel.Insert err: ", err)
		return &social.GroupPutinResp{}, err
	}

	return &social.GroupPutinResp{}, nil
}
