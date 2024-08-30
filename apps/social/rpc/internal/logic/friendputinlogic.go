package logic

import (
	"GoChat/apps/social/socialmodels"
	"GoChat/pkg/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendPutIn 好友申请
func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	fmt.Println("enter FriendPutIn logic")
	// 检测申请者是否已与申请对象为好友关系
	friends, err := l.svcCtx.FriendsModel.FindByUidAndFid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		fmt.Println("FriendsModel.FindByUidAndFid err: ", err)
		return nil, err
	}
	// 是好友关系
	if friends != nil {
		return &social.FriendPutInResp{}, nil
	}

	// 不是好友关系，检测申请者是否已向申请对象发送过好友申请请求
	friendRequests, err := l.svcCtx.FriendRequestsModel.FindByReqUidAndUserId(l.ctx, in.ReqUid, in.UserId)
	if err != nil && !errors.Is(err, socialmodels.ErrNotFound) {
		return nil, err
	}

	// 已发送过好友请求
	if friendRequests != nil {
		return &social.FriendPutInResp{}, err
	}

	// 未发送过好友申请请求，将这条请求保存到好友申请表中
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &socialmodels.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: sql.NullString{
			Valid:  true,
			String: in.ReqMsg,
		},
		ReqTime: time.Unix(in.ReqTime, 0),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})

	if err != nil {
		return nil, err
	}

	return &social.FriendPutInResp{}, nil
}
