package logic

import (
	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"
	"GoChat/apps/social/socialmodels"
	"GoChat/pkg/constants"
	"GoChat/pkg/wuid"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GroupCreateLogic) GroupCreate(in *social.GroupCreateReq) (*social.GroupCreateResp, error) {
	// 群模型实例
	groups := &socialmodels.Groups{
		Id:         wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Name:       in.Name,
		Icon:       in.Icon,
		CreatorUid: in.CreatorUid,
		IsVerify:   false,
	}

	// 在groups表中insert
	// 在group_members中insert
	err := l.svcCtx.GroupsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 在groups表中insert Groups实例
		_, err := l.svcCtx.GroupsModel.InsertWithSession(l.ctx, session, groups)
		if err != nil {
			fmt.Println("social rpc GroupCreate: social GroupsModel.InsertWithSession err: ", err)
			return err
		}

		// 在group_members中insert group_id - creator_id
		_, err = l.svcCtx.GroupMembersModel.InsertWithSession(l.ctx, session, &socialmodels.GroupMembers{
			GroupId:   groups.Id,
			UserId:    in.CreatorUid,
			RoleLevel: int64(constants.CreatorGroupRoleLevel),
		})
		if err != nil {
			fmt.Println("social rpc GroupCreate: social GroupMembersModel.InsertWithSession err: ", err)
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("social rpc GroupCreate: social GroupsModel.Trans err: ", err)
		return &social.GroupCreateResp{}, nil
	}
	return &social.GroupCreateResp{
		Id: groups.Id,
	}, nil
}
