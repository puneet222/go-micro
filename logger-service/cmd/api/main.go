package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	mongoURL = "mongodb://localhost:27017"
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
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("pinging mongodb client")

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	log.Println("Connected to mongodb")

	return client
}
