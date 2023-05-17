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

func ExtractToFolder(fileName string, folderName string) {
	// extract the file to directory
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
		defer fileStream.Close()
		err = untar.Untar(fileStream, folderName)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zip.Unzip(fileName, folderName)
	}
	fmt.Printf("Extracted file %s to actions-runner directory", fileName)

	// remove the compressed file
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Removed file %s", fileName)
}
