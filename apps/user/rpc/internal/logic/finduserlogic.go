package logic

import (
	"GoChat/apps/user/models"
	"context"
	"github.com/jinzhu/copier"

	"GoChat/apps/user/rpc/internal/svc"
	"GoChat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	// todo: add your logic here and delete this line
	var (
		userEntities []*models.Users
		err          error
	)

	// 根据不同的请求参数进行查询
	if in.Phone != "" {
		// 根据手机号查询用户
		userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
		if err == nil {
			userEntities = append(userEntities, userEntity)
		}
	} else if in.Name != "" {
		// 根据用户名查询用户列表
		userEntities, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Name)
	} else if len(in.Ids) > 0 {
		// 根据用户ID列表查询用户列表
		userEntities, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var resp []*user.UserEntity
	err = copier.Copy(&resp, &userEntities)
	if err != nil {
		return nil, err
	}

	// 返回查询结果
	return &user.FindUserResp{
		User: resp,
	}, nil
}
