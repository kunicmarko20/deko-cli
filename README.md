Deko CLI
========

```                                                ░▒▓ ✔    at 16:21:55  ▓▒░
Usage:
  deko-cli [command]

Available Commands:
  config      Provides sub commands around config for this cli tool
  help        Help about any command
  release     Creates release branch and a PR in current git repository

Flags:
  -h, --help      help for deko-cli
  -v, --version   version for deko-cli

Use "deko-cli [command] --help" for more information about a command.
```

## Config

By running `deko-cli config list`, you are able to see list of current configs,
you will see something like:

```yaml
github_access_token: 1234567890
```

For now only config you need to set before using the tool is `github_access_token`.
Github token is required to open a pull request. When creating a new token
be sure to give it `repo` permissions.

## Commands

* config
    * [list](#list)
    * [set](#set)
* [release](#release)


### config

#### list

```
Lists all saved config

Usage:
  deko-cli config list [flags]

Flags:
  -h, --help   help for list
```

#### set

```          
Creates or Updates a config.

Usage:
  deko-cli config set <key> <value> [flags]

Flags:
  -h, --help   help for set
```

### release

This command should be run from the local directory of the project you are doing
a release for. 

What does this command do in sequential order:
* Reads your git config to determine upstream repository
* Checkouts out to `master` branch
* Pulls new changes
* Creates new branch `release-yyyyMMdd` from `master`
* Merges new changes from "origin/staging" into `release-yyyyMMdd`
* Pushes new branch to remote
* Creates pull request on the remote repository
* Updates body of the pull request with all detected tickets from commits (If commit
has a prefix like: "[PROJECT-123]" it will be detected)
* Opens created PR in your default browser

```
Creates release branch and a PR on current git repository

Usage:
  deko-cli release [flags]

Aliases:
  release, r

Flags:
  -h, --help   help for release
```
