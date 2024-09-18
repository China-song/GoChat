package logic

import (
	"GoChat/apps/im/rpc/im"
	"GoChat/apps/social/rpc/socialclient"
	"GoChat/apps/user/rpc/user"
	"GoChat/pkg/bitmap"
	"GoChat/pkg/constants"
	"context"

	"GoChat/apps/im/api/internal/svc"
	"GoChat/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	chatLogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		l.Errorf("api im - GetChatLogReadRecords - call rpc Im.GetChatLog err: %v", err)
		return nil, err
	}

	var (
		chatlog = chatLogs.List[0]
		reads   = []string{chatlog.SendId}
		unreads []string
		ids     []string
	)
	// 分别设置已读未读
	switch constants.ChatType(chatlog.ChatType) {
	case constants.SingleChatType:
		// todo 私聊检测接收者是否已读
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId}
		} else {
			reads = append(reads, chatlog.RecvId)
		}
		ids = []string{chatlog.RecvId, chatlog.SendId}
	case constants.GroupChatType:
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatlog.RecvId,
		})
		if err != nil {
			l.Errorf("api im - GetChatLogReadRecords - call rpc Social.GroupUsers err: %v", err)
			return nil, err
		}

		aBitmap := bitmap.Load(chatlog.ReadRecords)
		for _, member := range groupUsers.List {
			ids = append(ids, member.UserId)
			if member.UserId == chatlog.SendId {
				continue
			}
			if aBitmap.IsSet(member.UserId) {
				reads = append(reads, member.UserId)
			} else {
				unreads = append(unreads, member.UserId)
			}
		}
	}

	userEntities, err := l.svcCtx.User.FindUser(l.ctx, &user.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		l.Errorf("api im - GetChatLogReadRecords - call rpc User.FindUser err: %v", err)
		return nil, err
	}
	userEntitySet := make(map[string]*user.UserEntity, len(userEntities.User))
	for i, userEntity := range userEntities.User {
		userEntitySet[userEntity.Id] = userEntities.User[i]
	}

	// 设置手机号码
	for i, read := range reads {
		if u := userEntitySet[read]; u != nil {
			reads[i] = u.Phone
		}
	}
	for i, unread := range unreads {
		if u := userEntitySet[unread]; u != nil {
			unreads[i] = u.Phone
		}
	}

	return &types.GetChatLogReadRecordsResp{
		Reads:   reads,
		UnReads: unreads,
	}, nil
}
