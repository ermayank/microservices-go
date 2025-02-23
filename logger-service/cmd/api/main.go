package main

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

const (
	port     = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
}

func main() {
	// Connect to Mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//Close Connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {
	// Set client options with authentication
	clientOptions := options.Client().ApplyURI(mongoURL).
		SetAuth(options.Credential{
			Username: "admin",
			Password: "password",
		})

	// Connect to MongoDB
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	// Ping the database to ensure connection is successful
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Println("Error pinging MongoDB:", err)
		return nil, err
	}

	return client, nil
}
