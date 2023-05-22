# gh-runner
A GitHub CLI extension to create a self hosted runner on your current local machine.
The runner agent will be downloaded and installed in the current directory.
By default the runner will be created for the current repository, but you can specify a different repository with the `--repository` flag.
You can also specify an organization or enterprise to create the runner for with the `--org` or `--enterprise` flags.
You can start and stop the runner as a process, or use 'service' to install it as a service.

## Installation
`gh extension install dedac/gh-runner`

## Usage
`gh runner create`

## Limitations
- You may get 403 errors if your gh cli doesn't have the correct permissions
- If you are developing in a codespace you need to `unset GITHUB_TOKEN` and log in manually with `gh auth login` to get the correct permissions - This may cause trouble with your ability to commit code via git.  Whenever the codespace stops and restarts, the GITHUB_TOKEN is refreshed, and you will have to auth gh again.

```
Usage:
  runner create [<options>] [flags]

Flags:
  -e, --enterprise string   Add the runner at the Enterprise level with the enterprise's name
  -g, --group string        Runner group to add the runner to, defaults to 'Default'
  -h, --help                help for create
  -l, --labels string       Comma-separated list of labels to add to the runner
  -o, --org string          Add the runner at the Organization level with the organization's name
  -s, --skip-download       Skip downloading the runner binary, because you already have one extracted in the current directory

Global Flags:
  -N, --name string   Name of the runner, creates a folder and a runner with this name, defualts to 'actions-runner'
                      When you set a name, you will need to use that name for all subsequent commands commands to that runner (default "actions-runner")
  -R, --repo string   Repository to use in OWNER/REPO format, defaults to the current repository
```
