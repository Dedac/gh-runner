package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Download(fileName string, url string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	httpClient := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return err
		},
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
}
