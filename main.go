package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/artyom/untar"
	"github.com/evilsocket/islazy/zip"
)

func main() {
	fileName, url, err := FindRunner()

	if err != nil {
		log.Fatal(err)
	}
	if url == "" {
		fmt.Println("Failed to get latest release")
		return
	}
	fmt.Printf("url is %s\n", url)

	err = Download(fileName, url)
	if err != nil {
		log.Fatal(err)
	}

	//extract the file to an actions-runner directory
	if runtime.GOOS != "windows" {
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		fileStream, err := gzip.NewReader(file)
		if err != nil {
			log.Fatal(err)
		}
		err = untar.Untar(fileStream, "actions-runner")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zip.Unzip(fileName, "actions-runner")
	}
	fmt.Printf("Extracted file %s to actions-runner directory", fileName)

	//remove the compressed file
	err = os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}

	//configure the runner

}
