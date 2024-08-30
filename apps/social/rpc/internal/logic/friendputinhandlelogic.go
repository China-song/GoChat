package logic

import (
	"GoChat/apps/social/socialmodels"
	"GoChat/pkg/constants"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"GoChat/apps/social/rpc/internal/svc"
	"GoChat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrFriendReqPassed  = errors.New("该好友申请已经通过")
	ErrFriendReqRefused = errors.New("该好友申请已被拒绝")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// 首先判断要处理的好友申请记录是否存在，若不存在则返回
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, uint64(in.FriendReqId))
	if err != nil {
		fmt.Println("FriendRequestsModel.FindOne err: ", err)
		return nil, err
	}

	// 好友申请记录存在，判断是否已处理
	switch constants.HandlerResult(friendRequest.HandleResult.Int64) {
	case constants.PassHandlerResult:
		return nil, ErrFriendReqPassed
	case constants.RefuseHandlerResult:
		return nil, ErrFriendReqRefused
	}

	fmt.Println("friend putin record exist, need to be handled!")
	// 修改处理状态 若为通过 添加好友表
	friendRequest.HandleResult.Int64 = int64(in.HandleResult)
	// 修改申请结果 -> 通过【建立两条好友关系记录】 -> 事务
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err := l.svcCtx.FriendRequestsModel.UpdateWithSession(l.ctx, session, friendRequest)
		if err != nil {
			return err
		}

		// TODO: 处理后是否应该删掉这条好友申请

		if constants.HandlerResult(in.HandleResult) != constants.PassHandlerResult {
			return nil
		}

		friends := []*socialmodels.Friends{
			{
				UserId:    friendRequest.UserId,
				FriendUid: friendRequest.ReqUid,
			}, {
				UserId:    friendRequest.ReqUid,
				FriendUid: friendRequest.UserId,
			},
		}

		_, err = l.svcCtx.FriendsModel.InsertsWithSession(l.ctx, session, friends...)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("handle trans err: ", err)
		return nil, err
	}
	return &social.FriendPutInHandleResp{}, nil
}
