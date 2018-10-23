package main

import (
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
)

func main() {

	http.HandleFunc("/", handle)
	http.HandleFunc("/update", Update)
	appengine.Main()
	// router := mux.NewRouter()
	// log.Fatal(http.ListenAndServe(":8080", router))
	// router.HandleFunc("/update", Update).Methods("GET")
	// router.HandleFunc("/", info)

	fmt.Println("should be here")
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Topify backend")
}

func Update(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.WithContext(r.Context(), r)
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Print("info", err.Error())
	}

	bkt := client.Bucket("topify-data")
	obj := bkt.Object("topify-list")
	writer := obj.NewWriter(ctx)
	writer.ContentType = "text/csv"

	fetcher := NewFetcher()
	repo := NewRepo(client, bkt, obj, writer)

	err = fetcher.FetchData()
	if err != nil {
		log.Print("info", err.Error())
	}

	err = repo.Upload()
	if err != nil {
		log.Print("info", err.Error())
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Updated"))
}
