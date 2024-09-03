package socialmodels

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ GroupsModel = (*customGroupsModel)(nil)

type (
	// GroupsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customGroupsModel.
	GroupsModel interface {
		groupsModel
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		InsertWithSession(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error)
		ListByGroupIds(ctx context.Context, ids []string) ([]*Groups, error)
	}

	customGroupsModel struct {
		*defaultGroupsModel
	}
)

// NewGroupsModel returns a model for the database table.
func NewGroupsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) GroupsModel {
	return &customGroupsModel{
		defaultGroupsModel: newGroupsModel(conn, c, opts...),
	}
}

func (m *customGroupsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

func (m *customGroupsModel) InsertWithSession(ctx context.Context, session sqlx.Session, data *Groups) (sql.Result, error) {
	groupsIdKey := fmt.Sprintf("%s%v", cacheGroupsIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, groupsRowsExpectAutoSet)
		return session.ExecCtx(ctx, query, data.Id, data.Name, data.Icon, data.Status, data.CreatorUid, data.GroupType, data.IsVerify, data.Notification, data.NotificationUid)
	}, groupsIdKey)
	return ret, err
}

func (m *customGroupsModel) ListByGroupIds(ctx context.Context, ids []string) ([]*Groups, error) {
	query := fmt.Sprintf("select %s from %s where `id` in ('%s')", groupsRows, m.table, strings.Join(ids, "','"))
	var resp []*Groups
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err

	}
}
