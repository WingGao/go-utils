package wmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
