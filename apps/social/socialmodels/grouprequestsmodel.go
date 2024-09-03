package socialmodels

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ GroupRequestsModel = (*customGroupRequestsModel)(nil)

type (
	// GroupRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupRequestsModel.
	GroupRequestsModel interface {
		groupRequestsModel
		FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error)
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) (sql.Result, error)
		ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error)
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) error
	}

	customGroupRequestsModel struct {
		*defaultGroupRequestsModel
	}
)

// NewGroupRequestsModel returns a model for the database table.
func NewGroupRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupRequestsModel {
	return &customGroupRequestsModel{
		defaultGroupRequestsModel: newGroupRequestsModel(conn, c, opts...),
	}
}

func (m *customGroupRequestsModel) FindByGroupIdAndReqId(ctx context.Context, groupId, reqId string) (*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `req_id` = ?", groupRequestsRows, m.table)
	var resp GroupRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, groupId, reqId)
	switch err {
	case nil:
		return &resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customGroupRequestsModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) (sql.Result, error) {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupRequestsRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource, data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult)
	}, groupRequestsIdKey)
	return ret, err
}

func (m *customGroupRequestsModel) ListNoHandler(ctx context.Context, groupId string) ([]*GroupRequests, error) {
	query := fmt.Sprintf("select %s from %s where `group_id` = ? and `handle_result` = 1 ", groupRequestsRows, m.table)
	var resp []*GroupRequests
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, groupId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customGroupRequestsModel) UpdateWithSession(ctx context.Context, session sqlx.Session, data *GroupRequests) error {
	groupRequestsIdKey := fmt.Sprintf("%s%v", cacheGroupRequestsIdPrefix, data.Id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, groupRequestsRowsWithPlaceHolder)
		return session.ExecCtx(ctx, query, data.ReqId, data.GroupId, data.ReqMsg, data.ReqTime, data.JoinSource, data.InviterUserId, data.HandleUserId, data.HandleTime, data.HandleResult, data.Id)
	}, groupRequestsIdKey)
	return err
}
