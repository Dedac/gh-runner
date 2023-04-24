package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/artyom/untar"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/evilsocket/islazy/zip"
)

func main() {
	repoOverride := flag.String("repo", "", "Repository to add the runner to. Defaults to the current repository")
	orgOverride := flag.String("org", "", "Organization to add the runner to. Defaults to the org of the current repository")
	entOverride := flag.String("ent", "", "Enterprise to add the runner to. Defaults to the enterprise of the current repository")
	flag.Parse()

	var repo repository.Repository
	var org string
	var ent string
	var URL string
	var err error

	if *repoOverride == "" {
		repo, err = gh.CurrentRepository()
		URL = fmt.Sprintf("https://%s/%s", repo.Host(), repo.Owner())
	} else {
		repo, err = repository.Parse(*repoOverride)
	}
	if err != nil {
		fmt.Println("could not determine repository to query: %w", err)
		return
	}

	if *orgOverride == "" {
		org = repo.Owner()
	} else if *orgOverride != "" {
		org = *orgOverride
	}

	if *entOverride == "" {
		ent = repo.Host()
	} else {
		ent = *entOverride
	}
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
		defer fileStream.Close()
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

	//execute config.sh to configure the runner with the repo or the org or the enterprise
	err = os.Chdir("actions-runner")
	if err != nil {
		log.Fatal(err)
	}
	var cmd *exec.Cmd
	if *repoOverride != "" {
		cmd = exec.Command("./config.sh", "--url", repo.Owner(), "--token", os.Getenv("GITHUB_TOKEN"), "--name", "actions-runner", "--labels", "actions,runner", "--unattended", "--repo", repo.Name())
	} else if *orgOverride != "" {
		cmd = exec.Command("./config.sh", "--url", URL, "--token", os.Getenv("GITHUB_TOKEN"), "--name", "actions-runner", "--labels", "actions,runner", "--unattended", "--org", org)
	} else if *entOverride != "" {
		cmd = exec.Command("./config.sh", "--url", URL, "--token", os.Getenv("GITHUB_TOKEN"), "--name", "actions-runner", "--labels", "actions,runner", "--unattended", "--enterprise", ent)
	} else {
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
