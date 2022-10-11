package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(c *mongo.Client) Models {
	client = c
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}
type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	collection := client.Database("logs").Collection("logs")
	opt := options.Find()
	opt.SetSort(bson.D{{"created_at", -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opt)
	if err != nil {
		log.Println("Error finding docs")
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding docs")
			return nil, err
		}
		logs = append(logs, &item)
	}
	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		return nil, err
	}
	result, err := collection.UpdateByID(ctx, bson.M{"_id": docId}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
