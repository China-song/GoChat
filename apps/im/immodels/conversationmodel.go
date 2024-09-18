package immodels

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
)

var _ ConversationModel = (*customConversationModel)(nil)

type (
	// ConversationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customConversationModel.
	ConversationModel interface {
		conversationModel
		FindByConversationId(ctx context.Context, conversationId string) (*Conversation, error)
		ListByConversationIds(ctx context.Context, ids []string) ([]*Conversation, error)
		UpdateMsg(ctx context.Context, chatLog *ChatLog) error
	}

	customConversationModel struct {
		*defaultConversationModel
	}
)

// NewConversationModel returns a model for the mongo.
func NewConversationModel(url, db, collection string) ConversationModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customConversationModel{
		defaultConversationModel: newDefaultConversationModel(conn),
	}
}

func MustConversationModel(url, db string) ConversationModel {
	return NewConversationModel(url, db, "conversation")
}

func (m *customConversationModel) FindByConversationId(ctx context.Context, conversationId string) (*Conversation, error) {
	var data Conversation

	err := m.conn.FindOne(ctx, &data, bson.M{"conversationId": conversationId})
	switch {
	case err == nil:
		return &data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// ListByConversationIds query conversations WHERE conversationId IN (ids)
func (m *customConversationModel) ListByConversationIds(ctx context.Context, ids []string) ([]*Conversation, error) {
	var data []*Conversation

	err := m.conn.Find(ctx, &data, bson.M{
		"conversationId": bson.M{
			"$in": ids,
		},
	})
	switch {
	case err == nil:
		return data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// UpdateMsg Update conversations WHERE conversationId = chatLog.ConversationId
//
// total = total+1 msg = chatLog
func (m *customConversationModel) UpdateMsg(ctx context.Context, chatLog *ChatLog) error {
	_, err := m.conn.UpdateOne(ctx,
		bson.M{"conversationId": chatLog.ConversationId},
		bson.M{
			// 更新会话总消息数
			"$inc": bson.M{"total": 1},
			"$set": bson.M{"msg": chatLog},
		},
	)
	return err
}
