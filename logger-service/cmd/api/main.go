package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/puneet222/go-micro/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Models data.Models
}

const (
	webPort  = "80"
	mongoURL = "mongodb://admin:password@mongo:27017/?maxPoolSize=20&w=majority"
)

func main() {
	client := connectToMongo()

	var app = Config{
		Models: data.New(client),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("starting logging service at %v", webPort)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("error while starting logger server", err)
	}

}

func connectToMongo() *mongo.Client {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to mongodb")

	return client
}
