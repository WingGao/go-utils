package wmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MgCollection struct {
	*mongo.Collection
}

func NewMgCollection(c *mongo.Collection) (mc *MgCollection) {
	mc = &MgCollection{Collection: c}
	return
}
func newTrue() *bool {
	b := true
	return &b
}
func (c *MgCollection) UpsertId(id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpdateOne(context.Background(), bson.D{
		{"_id", ToObjectId(id)},
	}, BSet(update), &options.UpdateOptions{Upsert: newTrue()})
}
func (c *MgCollection) UpdateId(id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return c.UpdateOne(context.Background(), bson.D{
		{"_id", ToObjectId(id)},
	}, update)
}
func (c *MgCollection) FindId(id interface{}, out interface{}) (err error) {
	sr := c.FindOne(context.Background(), bson.D{
		{"_id", ToObjectId(id)},
	})
	return DecodeSingleRes(sr, out)
}
func (c *MgCollection) RemoveId(id interface{}, ) (*mongo.DeleteResult, error) {
	return c.DeleteOne(context.Background(), bson.D{
		{"_id", id},
	})
}
func (c *MgCollection) RemoveOne(q interface{}) (*mongo.DeleteResult, error) {
	return c.DeleteOne(context.Background(), q)
}

// 创建索引
func (c *MgCollection) CreateIndex(field, idxName string, asc, unique bool) (string, error) {
	indexes := c.Indexes()
	var mv int32 = 1
	if !asc {
		mv = -1
	}
	mod := mongo.IndexModel{
		Keys:    bsonx.Doc{{field, bsonx.Int32(mv)}},
		Options: options.Index().SetName(idxName).SetUnique(unique).SetBackground(true),
	}
	return indexes.CreateOne(context.Background(), mod)
}

// 创建索引
func (c *MgCollection) CreateIndexes(fields []string, idxName string, unique bool) (string, error) {
	indexes := c.Indexes()
	v := bsonx.Int32(1)
	keys := make([]bsonx.Elem, len(fields))
	for i, k := range fields {
		keys[i] = bsonx.Elem{k, v}
	}
	mod := mongo.IndexModel{
		Keys:    bsonx.Doc(keys),
		Options: options.Index().SetName(idxName).SetUnique(unique).SetBackground(true),
	}

	return indexes.CreateOne(context.Background(), mod)
}
