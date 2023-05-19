package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/evilsocket/islazy/zip"
)

func ExtractToFolder(fileName string, folderName string) {
	// extract the file to directory
	if runtime.GOOS != "windows" {
		folderCommand := exec.Command("mkdir", "-p", folderName)
		err := folderCommand.Run()
		if err != nil {
			log.Fatal(err)
		}
		extractcmd := exec.Command("tar", "xzf", fileName, "-C", folderName)
		err = extractcmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zip.Unzip(fileName, folderName)
	}
	fmt.Printf("Extracted file %s to actions-runner directory \n", fileName)

	//remove the compressed file
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Removed file %s \n", fileName)
}
