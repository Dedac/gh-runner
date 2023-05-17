package main

import (
	"fmt"
	"os"

	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

func _main() error {
	var repo repository.Repository
	var URL string
	var name *string
	rootCmd := &cobra.Command{
		Use:   "runner <subcommand> [flags]",
		Short: "gh runner",
	}
	repoOverride := rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository to use in OWNER/REPO format")
	name = rootCmd.PersistentFlags().StringP("name", "N", "actions-runner", "Name of the runner, creates a folder and a runner with this name, defualts to 'actions-runner'")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
		if *repoOverride != "" {
			repo, err = repository.Parse(*repoOverride)
		} else {
			repo, err = repository.Current()
		}
		URL = fmt.Sprintf("https://%s/%s/%s", repo.Host, repo.Owner, repo.Name)
		return
	}

	createCmd := &cobra.Command{
		Use:   "create [<options>]",
		Short: "Create a new runner with the given options",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runCreate(cmd, repo, false, *name, URL)
			return
		},
	}
	createCmd.Flags().StringP("org", "o", "", "Add the runner at the Organization level with the organization's name")
	createCmd.Flags().StringP("enterprise", "e", "", "Add the runner at the Enterprise level with the enterprise's name")
	createCmd.Flags().StringP("labels", "l", "", "Comma-separated list of labels to add to the runner")
	createCmd.Flags().StringP("group", "g", "", "Runner group to add the runner to, defaults to 'Default'")
	createCmd.Flags().BoolP("skip-download", "s", false, "Skip downloading the runner binary, because you already have one extracted in the current directory")

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start an already configured runner",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runStart(*name)
			return
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a runner that is currently running",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runStop(*name)
			return
		},
	}

	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "create a service (and start it) on this machine to keep the runner running",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			createService(*name)
			return
		},
	}

	serviceStopCmd := &cobra.Command{
		Use:   "serviceStop",
		Short: "Stop the runner configured on this machine",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			stopService(*name)
			return
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove configuration and runner",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return runCreate(cmd, repo, true, *name, URL)
		},
	}

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(serviceStopCmd)

	return rootCmd.Execute()
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "X %s \n", err.Error())
	}
}
