# Jira Cloud Github Action

## Introduction

A Github action that does things in Jira Cloud when triggered from a Github Actions workflow.

Currently supported:

* Adding a comment on a push or a pull request

## Usage

Use a workflow like this:

```yaml
on:
    - push
    - pull_request

jobs:
    jiracomment:
        name: "Comment on Jira"
        runs-on: ubuntu-latest
        steps:
            -   name: Checkout repository
                uses: actions/checkout@v4
                with:
                    fetch-depth: "0"
            -   uses: dodevops/jira-cloud-github-action@feature/dpr/first-version
                with:
                    url: "https://mycompany.atlassian.net"
                    username: "serviceuser@company.com"
                    token: ${{ secrets.JIRA_SERVICEUSER_TOKEN }}
```

This will automatically send comments to Jira tickets identified by commits that contain a specific format.

For example, these three commits in a pull request created by "user4 <user4@company.com>":

* `(SUP-1234): Fixed problem` by user1 <user1@company.com>
* `(SUP-1234): Added feature` by user2 <user2@company.com>
* `(SUP-2345): Removed feature` by user3 <user3@company.com>

will result in these two comments (the markdown format is converted to Jira markdown):

Comment on SUP-1234
```
*Pull Request*: [user4](mailto:user4@company.com) sent 2 commit(s) to [Pull request title]({{ https://linktopullrequest }})

* (SUP-1234): Fixed problem ([user1](mailto:user1@company.com))
* (SUP-1234): Added feature ([user2](mailto:user2@company.com))
```

Comment on SUP-2345
```
*Pull Request*: [user4](mailto:user4@company.com) sent 1 commit(s) to [Pull request title]({{ https://linktopullrequest }})

* (SUP-2345): Removed feature ([user3](mailto:user3@company.com))
```

## Options

* *url*: The URL to your Jira Cloud instance (required)
* *username*: The user name for API requests (required)
* *token*: The user's token for API requests (required)
* *loglevel*: The loglevel to use [INFO]
* *commitParseRegExp*: The Regexp string to use to parse Jira issue ids from commit messages. Requires a named group called id that catches the whole issue id and a group project that catches only the project key [\\((?P<id>(?P<project>[A-Za-z]+)-[0-9]+)\\)"]
* *commentTemplate*: The go template for the issue comment. Takes the CommitsInfo struct as an input [default see below]
  ```yaml
    *{{ .CommitsInfo.Type }}*: [{{ .CommitsInfo.AuthorName }}](mailto:{{ .CommitsInfo.AuthorEmail }}) sent {{ .Commits | len }} commit(s) to [{{ .CommitsInfo.Target }}]({{ .CommitsInfo.Target }})
    
    {{ range .Commits }}
    * {{ .Message }} ([{{ .AuthorName }}](mailto:{{ .AuthorEmail }}))
    {{ end }}
    ```
* *onlyProjects*: If specified, only issues with these project keys are processed
