package user

import (
	"GoChat/apps/user/rpc/user"
	"GoChat/pkg/constants"
	"context"
	"github.com/jinzhu/copier"

	"GoChat/apps/user/api/internal/svc"
	"GoChat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登入
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.User.Login(l.ctx, &user.LoginReq{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var res types.LoginResp
	err = copier.Copy(&res, loginResp)
	if err != nil {
		return nil, err
	}

	// 将用户ID和在线状态"1"存储到Redis的hash中，标记用户为在线。
	// 这里使用Redis来管理在线用户，是因为Redis的高并发读写性能和键值对存储特性适合此类场景。
	err = l.svcCtx.Redis.HsetCtx(l.ctx, constants.RedisOnlineUser, loginResp.Id, "1")
	if err != nil {
		// 如果设置Redis中用户在线状态失败，返回错误。
		return nil, err
	}

	return &res, nil
}
