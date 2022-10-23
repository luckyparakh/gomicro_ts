package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var mongodb *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting logging service")
	client, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}
	// go app.serve()
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Fail to start logging service")
	}
}
func (app *Config) serve() {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Panic("Fail to start logging service")
	}
}
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("unable to connect to db")
		return nil, err
	}
	return c, nil
}
