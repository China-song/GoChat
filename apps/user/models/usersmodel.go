package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
)

var _ UsersModel = (*customUsersModel)(nil)
var (
	cacheUsersPhonePrefix string = "cache:users:phone:"
)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		FindByPhone(ctx context.Context, phone string) (*Users, error)
		ListByName(ctx context.Context, name string) ([]*Users, error)
		ListByIds(ctx context.Context, ids []string) ([]*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}

func (m *customUsersModel) FindByPhone(ctx context.Context, phone string) (*Users, error) {
	//usersPhoneKey := fmt.Sprintf("%s%v", cacheUsersPhonePrefix, phone)
	var resp Users
	query := fmt.Sprintf("select %s from %s where `phone` = ? limit 1", usersRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, phone)
	switch {
	case err == nil:
		return &resp, nil
	case errors.Is(err, sqlc.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) ListByName(ctx context.Context, name string) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s where `nickname` like ?", usersRows, m.table)
	var resp []*Users
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, fmt.Sprint("%", name, "%"))
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}

func (m *customUsersModel) ListByIds(ctx context.Context, ids []string) ([]*Users, error) {
	query := fmt.Sprintf("select %s from %s where `id` in ('%s')", usersRows, m.table, strings.Join(ids, "','"))
	var resp []*Users
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return resp, nil
	default:
		return nil, err
	}
}
