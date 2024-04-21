package common

import (
	"context"
	"errors"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDB(col string) *mongo.Collection {
	return db.Collection(col)
}

func InitDB() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		return errors.New("MongoDB URI empty...")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err !=nil {
		return err
	}
	db = client.Database("book_inventory")
	return nil
}

func CloseDB() error {
	return db.Client().Disconnect(context.Background())
}