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
	repoOverride := rootCmd.PersistentFlags().StringP("repo", "R", "", "Repository to use in OWNER/REPO format, defaults to the current repository")
	name = rootCmd.PersistentFlags().StringP("name", "N", "actions-runner", "Name of the runner, creates a folder and a runner with this name, defualts to 'actions-runner' \nWhen you set a name, you will need to use that name for all subsequent commands commands to that runner")

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
		Short: "Start an already configured runner as a process in your current context",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runStart(*name)
			return
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop runner processes that are currently running in a local process",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runStop(*name)
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

	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the runner with a machine-level service",
	}

	serviceCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "create a service (and start it) on this machine to keep the runner running",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			createService(*name)
			return
		},
	}

	serviceStartCmd := &cobra.Command{
		Use:   "start",
		Short: "create a service (and start it) on this machine to keep the runner running",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			runService(*name)
			return
		},
	}

	serviceStopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the runner configured on this machine",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			stopService(*name)
			return
		},
	}

	serviceRemoveCmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove the service configured on this machine, ",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			removeService(*name)
			return
		},
	}

	removeCmd.Flags().StringP("org", "o", "", "Remove the runner at the Organization level with the organization's name")
	removeCmd.Flags().StringP("enterprise", "e", "", "Remove the runner at the Enterprise level with the enterprise's name")
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceCreateCmd)
	serviceCmd.AddCommand(serviceStartCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceRemoveCmd)

	return rootCmd.Execute()
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "X %s \n", err.Error())
	}
}
