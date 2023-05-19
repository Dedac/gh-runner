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
- Currently only supports Linux and MacOS (Windows is in progress)
- You may get 403 errors if your gh cli doesn't have the correct permissions
- If you are developing in a codespace you need to `unset GITHUB_TOKEN` and log in manually with `gh auth login` to get the correct permissions
