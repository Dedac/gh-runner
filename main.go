package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cli/go-gh"
)

func main() {
	ghRest, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	latestRelease := struct{ Name string }{}
	err = ghRest.Get("repos/actions/runner/releases/latest", &latestRelease)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Getting Latest Release, %s\n", latestRelease.Name)
	fileName := Build_Filename(strings.Split(latestRelease.Name, "v")[1])

	url := fmt.Sprintf("https://github.com/actions/runner/releases/download/%s/%s", latestRelease.Name, fileName)

	if url == "" {
		fmt.Println("Failed to get latest release")
		return
	}
	fmt.Printf("url is %s\n", url)

	//download the file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	httpClient := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d", fileName, size)

	//extract the file to an actions-runner directory
	//err = Untar(fileName, "actions-runner")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Extracted file %s to actions-runner directory", fileName)

}

// Build file name
func Build_Filename(version string) string {
	ext := "tar.gz"
	os := ""
	switch runtime.GOOS {
	case "darwin":
		os = "osx"
	case "windows":
		os = "win"
		ext = "zip"
	default:
		os = "linux"
	}

	arch := ""
	switch runtime.GOARCH {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	default:
		fmt.Printf("Unsupported Architecture: %s, exiting", runtime.GOARCH)
	}

	// combine OS, Arch, and name to build the file name
	fileName := fmt.Sprintf("actions-runner-%s-%s-%s.%s", os, arch, version, ext)

	return fileName
}
