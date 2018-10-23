package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type country struct {
	Name string
	Code string
	URI  string
}

type Fetcher interface {
	FetchData() error
	CreateCSV() error
}

type fetcher struct {
	countrySlice []country
}

func NewFetcher() fetcher {
	return fetcher{}
}

func (f *fetcher) FetchData() error {

	csvFile, err := os.Open("./spotify_countries.csv")
	if err != nil {
		fmt.Print(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			} else {
				errors.New("failed")
			}
		}

		f.countrySlice = append(f.countrySlice, country{
			Name: line[0],
			Code: line[1],
		})
	}

	for _, c := range f.countrySlice {
		if c.Code != "" {
			fileURL := fmt.Sprintf("https://spotifycharts.com/regional/%s/daily/latest/download", c.Code)
			outputPath := fmt.Sprintf("./temp_data/%s.csv", c.Code)
			log.Print(c.Code + " downloaded")
			err = DownloadFile(outputPath, fileURL)
			if err != nil {
				panic(err)
			}
		}
	}

	err = f.CreateCSV()
	if err != nil {
		return err
	}

	return nil
}

func (f *fetcher) CreateCSV() error {

	topify, err := os.Create("./topify-list.csv")
	if err != nil {
		return err
	}

	defer topify.Close()

	for _, d := range f.countrySlice {
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

		fmt.Fprintf(topify, "%s,%s\n", strings.ToUpper(d.Code), strings.Split(line[4], "/")[4])
		log.Print(fmt.Sprintf("%s saved with song: %s by %s \n", strings.ToUpper(d.Code), line[1], line[2]))
	}

	return nil
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
