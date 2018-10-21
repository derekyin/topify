package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	logger "github.com/apsdehal/go-logger"
)

type Repo interface {
	Update(c []country) error
}

type repo struct {
	log      *logger.Logger
	fbClient *firestore.Client
}

func NewRepo(fb *firestore.Client, logger *logger.Logger) repo {
	return repo{
		log:      logger,
		fbClient: fb,
	}
}

func (r *repo) Update(c []country) error {

	states := r.fbClient.Collection("top_songs")

	for _, a := range c {

		_, err := states.Doc(a.Code).Set(context.Background(), map[string]interface{}{
			"Name": a.Name,
			"Code": a.Code,
			"URI":  a.URI,
		})
		if err != nil {
			r.log.Log("error", fmt.Sprintf("error %s ", err.Error()))
		}
		r.log.Log("info", fmt.Sprintf("country %s updated", a.Name))
	}

	return nil
}
