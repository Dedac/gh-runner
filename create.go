package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

func runCreate(cmd *cobra.Command, repo repository.Repository, remove bool, folderName string, URL string) (err error) {
	orgOverride, _ := cmd.Flags().GetString("org")
	entOverride, _ := cmd.Flags().GetString("enterprise")
	labels, _ := cmd.Flags().GetString("labels")
	group, _ := cmd.Flags().GetString("group")
	skipDownload, _ := cmd.Flags().GetBool("skip-download")

	var org string
	var ent string

	if orgOverride != "" {
		org = orgOverride
		URL = fmt.Sprintf("https://%s/%s", repo.Host, org)
	}

	if entOverride != "" {
		ent = entOverride
		URL = fmt.Sprintf("https://%s/enterprises/%s", repo.Host, ent)
	}

	if !remove && !skipDownload {
		// Get the correct runner for the current platform
		fileName, url := FindRunner()
		Download(fileName, url)
		ExtractToFolder(fileName, folderName)
	}

	err = os.Chdir(folderName)
	if err != nil {
		log.Fatal(err)
	}

	token := GetToken(repo, org, ent, remove)

	var args []string
	if remove {
		args = []string{"remove", "--token", token}
	} else {
		args = []string{"--url", URL, "--token", token, "--name", folderName, "--unattended"}
		if labels != "" {
			args = append(args, "--labels", labels)
		}
		if group != "" {
			args = append(args, "--runnergroup", group)
		}
	}
	configCmd := exec.Command("./config.sh", args...)

	configCmd.Stdout = os.Stdout
	configCmd.Stderr = os.Stderr
	err = configCmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
