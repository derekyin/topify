package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
)

type country struct {
	Name string
	Code string
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()

	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		// TODO: Handle error.
	}

	bkt := client.Bucket("topify")
	obj := bkt.Object("data")
	w := obj.NewWriter(ctx)
	w.ContentType = "text/csv"

	csvFile, err := os.Open("./backend/spotify_countries.csv")

	if err != nil {
		fmt.Print(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	var countries []country

	for {
		line, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				errors.New("failed")
			}
		}

		countries = append(countries, country{
			Name: line[0],
			Code: line[1],
		})
	}

	for _, c := range countries {
		if c.Code != "" {
			fileURL := fmt.Sprintf("https://spotifycharts.com/regional/%s/daily/latest/download", c.Code)
			outputPath := fmt.Sprintf("./backend/countries_top_songs/%s.csv", c.Code)
			fmt.Println(c.Code + " downloaded")
			err = DownloadFile(outputPath, fileURL)
			if err != nil {
				panic(err)
			}
		}
	}

	topify, _ := os.Create("./backend/topify-list.csv")

	defer topify.Close()

	for _, d := range countries {
		countryFile, _ := os.Open(fmt.Sprintf("./backend/countries_top_songs/%s.csv", d.Name))
		reader := csv.NewReader(bufio.NewReader(countryFile))

		_, error := reader.Read()
		if error != nil {
			errors.New("failed")
		}
		_, error = reader.Read()
		if error != nil {
			errors.New("failed")
		}
		line, error := reader.Read()
		if error != nil {
			errors.New("failed")
		}

		fmt.Fprintf(topify, "%s,%s\n", strings.ToUpper(d.Code), strings.Split(line[4], "/")[4])
		a := []byte(fmt.Sprintf("%s,%s\n", strings.ToUpper(d.Code), strings.Split(line[4], "/")[4]))
		w.Write(a)
		fmt.Printf("%s saved with song: %s by %s \n", strings.ToUpper(d.Code), line[1], line[2])
	}

	if err := w.Close(); err != nil {
		errors.New("failed to close GCP storage bucket")
	}
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Back end for Topify")
}
