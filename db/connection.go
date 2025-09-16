package db

import (
	"context"
	"gin/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func Init() {
	DB = connectDB()
}
func connectDB() *mongo.Client {
	log.Println("Connection to MongoDB...")
	log.Println("Config", config.DB_URL, config.DB_NAME)
	clientOptions := options.Client().ApplyURI(config.DB_URL)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB is not reachable: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	return client
}

func OpenCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		log.Fatal("MongoDB client is not initialized. Please call ConnectDB first.")
	}
	return DB.Database(config.DB_NAME).Collection(collectionName)
}
