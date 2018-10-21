package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	logger "github.com/apsdehal/go-logger"
	"google.golang.org/api/option"
)

func main() {

	logger, err := logger.New("topify_data_fetcher", 1, os.Stdout)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	opt := option.WithCredentialsFile("./topify-217503-firebase-adminsdk-vfh6n-9cdf695562.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	fetcher := NewFetcher(client, logger)
	repo := NewRepo(client, logger)

	countryData, err := fetcher.FetchData()
	if err != nil {
		logger.Log("info", err.Error())
	}

	err = repo.Update(countryData)
	if err != nil {
		logger.Log("info", err.Error())
	}

}
