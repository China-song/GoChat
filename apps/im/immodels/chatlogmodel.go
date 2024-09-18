package immodels

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var _ ChatLogModel = (*customChatLogModel)(nil)

type (
	// ChatLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customChatLogModel.
	ChatLogModel interface {
		chatLogModel
		ListBySendTime(ctx context.Context, conversationId string, startSendTime, endSendTime, limit int64) ([]*ChatLog, error)
		ListByMsgIds(ctx context.Context, msgIds []string) ([]*ChatLog, error)
		UpdateMarkRead(ctx context.Context, msgId primitive.ObjectID, readRecords []byte) error
	}

	customChatLogModel struct {
		*defaultChatLogModel
	}
)

// NewChatLogModel returns a model for the mongo.
func NewChatLogModel(url, db, collection string) ChatLogModel {
	conn := mon.MustNewModel(url, db, collection)
	return &customChatLogModel{
		defaultChatLogModel: newDefaultChatLogModel(conn),
	}
}

func MustChatLogModel(url, db string) ChatLogModel {
	return NewChatLogModel(url, db, "chat_log")
}

// ListBySendTime 返回startSendTime之前(end< t < start)的消息
func (m *customChatLogModel) ListBySendTime(ctx context.Context, conversationId string, startSendTime, endSendTime, limit int64) ([]*ChatLog, error) {
	var data []*ChatLog

	opt := options.FindOptions{
		Limit: &DefaultChatLogLimit,
		Sort: bson.M{
			"sendTime": -1,
		},
	}
	if limit > 0 {
		opt.Limit = &limit
	}

	filter := bson.M{
		"conversationId": conversationId,
	}

	if endSendTime > 0 {
		filter["sendTime"] = bson.M{
			"$gt":  endSendTime,
			"$lte": startSendTime,
		}
	} else {
		filter["sendTime"] = bson.M{
			"$lt": startSendTime,
		}
	}
	err := m.conn.Find(ctx, &data, filter, &opt)
	switch {
	case err == nil:
		return data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customChatLogModel) ListByMsgIds(ctx context.Context, msgIds []string) ([]*ChatLog, error) {
	var data []*ChatLog
	ids := make([]primitive.ObjectID, 0, len(msgIds))
	for _, id := range msgIds {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, ErrInvalidObjectId
		}
		ids = append(ids, oid)
	}
	filter := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	err := m.conn.Find(ctx, &data, filter)
	switch {
	case err == nil:
		return data, nil
	case errors.Is(err, mon.ErrNotFound):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customChatLogModel) UpdateMarkRead(ctx context.Context, msgId primitive.ObjectID, readRecords []byte) error {
	_, err := m.conn.UpdateOne(ctx, bson.M{"_id": msgId}, bson.M{"$set": bson.M{
		"readRecords": readRecords,
		"updateAt":    time.Now(),
	}})
	return err
}
