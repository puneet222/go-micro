package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/puneet222/go-micro/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Models data.Models
}

const (
	webPort  = "80"
	mongoURL = "mongodb://admin:password@localhost:27017/?maxPoolSize=20&w=majority"
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

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("error while starting logger server", err)
	}

	log.Printf("logging server started at %v", webPort)

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

	log.Println("Connected to mongodb")

	return client
}
