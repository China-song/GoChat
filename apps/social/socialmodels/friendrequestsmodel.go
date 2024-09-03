package socialmodels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error)
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *FriendRequests) error
		ListNoHandler(ctx context.Context, uid string) ([]*FriendRequests, error)
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}

func (m *customFriendRequestsModel) FindByReqUidAndUserId(ctx context.Context, rid, uid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `req_uid` = ? and `user_id` = ?", friendRequestsRows, m.table)
	var resp FriendRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, rid, uid)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customFriendRequestsModel) UpdateWithSession(ctx context.Context, session sqlx.Session, data *FriendRequests) error {
	friendRequestsIdKey := fmt.Sprintf("%s%v", cacheFriendRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, friendRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt, data.Id)
	}, friendRequestsIdKey)
	return err
}

// ListNoHandler returns user(uid)'s friend request lists which are not be handled
func (m *customFriendRequestsModel) ListNoHandler(ctx context.Context, uid string) ([]*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where `handle_result` = 1 and `req_uid` = ?", friendRequestsRows, m.table)
	var resp []*FriendRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, uid)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
