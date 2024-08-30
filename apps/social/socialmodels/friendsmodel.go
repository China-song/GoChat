package socialmodels

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface {
		friendsModel
		FindByUidAndFid(ctx context.Context, uid, fid string) (*Friends, error)
		InsertsWithSession(ctx context.Context, session sqlx.Session, data ...*Friends) (sql.Result, error)
		LiseByUserId(ctx context.Context, userId string) ([]*Friends, error)
	}

	customFriendsModel struct {
		*defaultFriendsModel
	}
)

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c, opts...),
	}
}

func (m *customFriendsModel) FindByUidAndFid(ctx context.Context, uid, fid string) (*Friends, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `friend_uid` = ?", friendsRows, m.table)
	var resp Friends
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, uid, fid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customFriendsModel) InsertsWithSession(ctx context.Context, session sqlx.Session, data ...*Friends) (sql.Result, error) {
	var (
		sql  strings.Builder
		args []any
	)

	if len(data) == 0 {
		return nil, nil
	}

	// insert into table (field1, field2, ...) values (数据), (数据)
	s := fmt.Sprintf("insert into %s (%s) values ", m.table, friendsRowsExpectAutoSet)
	fmt.Println(s)
	sql.WriteString(s)

	for i, v := range data {
		sql.WriteString("(?, ?, ?, ?)")
		args = append(args, v.UserId, v.FriendUid, v.Remark, v.AddSource)
		if i == len(data)-1 {
			break
		}

		sql.WriteString(",")
	}

	return session.ExecCtx(ctx, sql.String(), args...)
}

func (m *customFriendsModel) LiseByUserId(ctx context.Context, userId string) ([]*Friends, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? ", friendsRows, m.table)
	var resp []*Friends
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
