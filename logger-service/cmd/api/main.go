package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"log-service/cmd/data"
	"net"
	"net/http"
	"net/rpc"

	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// Connect to Mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	//Close Connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	//Register RPC Server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	go app.gRPCListen()

	//Start Web Server
	log.Println("Starting mongo server on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// Start RPC Server
func (app *Config) rpcListen() error {
	log.Println("Starting rpc server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Panic(err)
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
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
