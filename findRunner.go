package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/cli/go-gh"
)

func FindRunner() (filename string, url string, err error) {
	ghRest, err := gh.RESTClient(nil)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	latestRelease := struct{ Name string }{}
	err = ghRest.Get("repos/actions/runner/releases/latest", &latestRelease)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	fmt.Printf("Getting Latest Release, %s\n", latestRelease.Name)
	fileName := Build_Filename(strings.Split(latestRelease.Name, "v")[1])

	url = fmt.Sprintf("https://github.com/actions/runner/releases/download/%s/%s", latestRelease.Name, fileName)
	return fileName, url, nil
}

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
