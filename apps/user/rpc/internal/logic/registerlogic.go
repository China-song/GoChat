package logic

import (
	"GoChat/apps/user/models"
	"GoChat/pkg/ctxdata"
	"GoChat/pkg/encrypt"
	"GoChat/pkg/wuid"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GoChat/apps/user/rpc/internal/svc"
	"GoChat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegistered = errors.New("该手机号已被注册过")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// todo: add your logic here and delete this line

	// 1. 判断该手机号是否已被注册
	_, err := l.svcCtx.UsersModel.FindByPhone(l.ctx, in.Phone)
	fmt.Println("phone = ", in.Phone)
	fmt.Println("err = ", err)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		// 该手机号已注册
		return nil, ErrPhoneIsRegistered
	}

	// 2. 可注册 创建用户 给密码加密
	userEntity := &models.Users{
		Id:       wuid.GenUid(l.svcCtx.Config.Mysql.DataSource),
		Avatar:   in.Avatar,
		Nickname: in.Nickname,
		Phone:    in.Phone,
		Sex: sql.NullInt64{
			Int64: int64(in.Sex),
			Valid: true,
		},
	}

	if len(in.Password) > 0 {
		hashedPassword, err := encrypt.GenPasswordHash([]byte(in.Password))
		if err != nil {
			return nil, err
		}
		userEntity.Password = sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		}
	}

	// 3. 新增用户
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, err
	}

	// 4. 生成token
	now := time.Now().Unix()
	exp := l.svcCtx.Config.Jwt.AccessExpire
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, exp, userEntity.Id)
	if err != nil {
		return nil, err
	}

	return &user.RegisterResp{
		Token:  token,
		Expire: now + exp,
	}, nil
}
