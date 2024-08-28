package logic

import (
	"GoChat/apps/user/models"
	"GoChat/pkg/ctxdata"
	"GoChat/pkg/encrypt"
	"context"
	"errors"
	"time"

	"GoChat/apps/user/rpc/internal/svc"
	"GoChat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegistered = errors.New("该手机号未注册")
	ErrUserPwdError       = errors.New("用户密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// todo: add your logic here and delete this line

	// 1. 查找用户是否注册
	userEntity, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrPhoneNotRegistered
		}
		return nil, err
	}

	// 2. 判断密码是否正确
	if !encrypt.ValidatePasswordHash(in.Password, userEntity.Password.String) {
		return nil, ErrUserPwdError
	}

	// 3. 生成token
	now := time.Now().Unix()
	exp := l.svcCtx.Config.Jwt.AccessExpire
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, exp, userEntity.Id)
	if err != nil {
		return nil, err
	}
	return &user.LoginResp{
		Token:  token,
		Expire: now + exp,
		Id:     userEntity.Id,
		User: &user.UserEntity{
			Id:       userEntity.Id,
			Avatar:   userEntity.Avatar,
			Nickname: userEntity.Nickname,
			Phone:    userEntity.Phone,
			Status:   int32(userEntity.Status.Int64),
			Sex:      int32(userEntity.Sex.Int64),
		},
	}, nil
}
