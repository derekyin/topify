package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	logger "github.com/apsdehal/go-logger"
)

type country struct {
	Name string
	Code string
	URI  string
}

type Fetcher interface {
	FetchData() ([]country, error)
}

type fetcher struct {
	log      *logger.Logger
	fbClient *firestore.Client
}

func NewFetcher(fb *firestore.Client, logger *logger.Logger) fetcher {
	return fetcher{
		log:      logger,
		fbClient: fb,
	}
}

func (f *fetcher) FetchData() ([]country, error) {

	csvFile, err := os.Open("./spotify_countries.csv")
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
			outputPath := fmt.Sprintf("./temp_data/%s.csv", c.Code)
			fmt.Println(c.Code + " downloaded")
			err = DownloadFile(outputPath, fileURL)
			if err != nil {
				panic(err)
			}
		}
	}

	// topify, _ := os.Create("./topify-list.csv")

	// defer topify.Close()

	countriesData := []country{}

	for _, d := range countries {
		countryFile, _ := os.Open(fmt.Sprintf("./temp_data/%s.csv", d.Code))
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

		// fmt.Fprintf(topify, "%s,%s\n", strings.ToUpper(d.Code), strings.Split(line[4], "/")[4])
		// a := []byte(fmt.Sprintf("%s,%s\n", strings.ToUpper(d.Code), strings.Split(line[4], "/")[4]))
		// w.Write(a)
		d.URI = strings.Split(line[4], "/")[4]
		countriesData = append(countriesData, d)
		fmt.Printf("%s saved with song: %s by %s \n", strings.ToUpper(d.Code), line[1], line[2])
	}

	return countriesData, nil
}

func DownloadFile(filepath string, url string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
