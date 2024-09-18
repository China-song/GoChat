package immodels

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ ConversationsModel = (*customConversationsModel)(nil)

type (
	// ConversationsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationsModel.
	ConversationsModel interface {
		conversationsModel
		FindByUserId(ctx context.Context, uid string) (*Conversations, error)
		UpdateOrInsert(ctx context.Context, data *Conversations) (*mongo.UpdateResult, error)
	}

	customConversationsModel struct {
		*defaultConversationsModel
	}
)

// NewConversationsModel returns a model for the mongo.
func NewConversationsModel(url, db, collection string) ConversationsModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationsModel{
		defaultConversationsModel: newDefaultConversationsModel(conn),
	}
}

func MustConversationsModel(url, db string) ConversationsModel {
	return NewConversationsModel(url, db, "conversations")
}

// FindByUserId query conversations WHERE userId = uid
func (m *customConversationsModel) FindByUserId(ctx context.Context, uid string) (*Conversations, error) {
	var data Conversations

	err := m.conn.FindOne(ctx, &data, bson.M{"userId": uid})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customConversationsModel) UpdateOrInsert(ctx context.Context, data *Conversations) (*mongo.UpdateResult, error) {
	if data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	} else {
		data.UpdateAt = time.Now()
	}

	res, err := m.conn.UpdateOne(ctx, bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	}, options.Update().SetUpsert(true))
	return res, err
}
