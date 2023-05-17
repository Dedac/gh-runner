package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
)

func main() {
	repoOverride := flag.String("repo", "", "Repository to add the runner to. Defaults to the current repository")
	orgOverride := flag.String("org", "", "Organization to add the runner to")
	entOverride := flag.String("enterprise", "", "Enterprise to add the runner to")
	labels := flag.String("labels", "", "Comma-separated list of labels to add to the runner")
	name := flag.String("name", "actions-runner", "Name of the runner, creates a folder and a runner with this name, defualts to 'actions-runner'")
	group := flag.String("group", "", "Runner group to add the runner to, defaults to 'Default'")
	remove := flag.Bool("remove", false, "Remove the existing configured runner")
	skipDownload := flag.Bool("skip-download", false, "Skip downloading the runner binary, because you already have one extracted")
	flag.Parse()

	var repo repository.Repository
	folderName := *name
	var org string
	var ent string
	var URL string
	var err error

	if *repoOverride == "" {
		repo, err = gh.CurrentRepository()
		URL = fmt.Sprintf("https://%s/%s/%s", repo.Host(), repo.Owner(), repo.Name())
	} else {
		repo, err = repository.Parse(*repoOverride)
	}
	if err != nil {
		fmt.Println("could not determine repository to query: %w", err)
		return
	}

	if *orgOverride != "" {
		org = *orgOverride
		URL = fmt.Sprintf("https://%s/%s", repo.Host(), org)
	}

	if *entOverride != "" {
		ent = *entOverride
		URL = fmt.Sprintf("https://%s/enterprises/%s", repo.Host(), ent)
	}

	if !*remove && !*skipDownload {
		// Get the correct runner for the current platform
		fileName, url := FindRunner()
		Download(fileName, url)
		ExtractToFolder(fileName, folderName)
	}

	err = os.Chdir(folderName)
	if err != nil {
		log.Fatal(err)
	}

	token := GetToken(repo, org, "", *remove)

	var args []string
	if *remove {
		args = []string{"remove", "--token", token}
	} else {
		args = []string{"--url", URL, "--token", token, "--name", *name, "--unattended"}
		if *labels != "" {
			args = append(args, "--labels", *labels)
		}
		if *group != "" {
			args = append(args, "--runnergroup", *group)
		}
	}
	cmd := exec.Command("./config.sh", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	runcmd := exec.Command("./run.sh")
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	err = runcmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
