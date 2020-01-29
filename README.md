# git-code

Simple cli tool to assist in searching and cloning Organizational repositories.
This quick tool simply fills a need to be able to quickly search and clone
company repositories.  Making it so I don't have to leave the shell...

## Setup

Currently the only authentication method being used is a Git Personal Access token.
An access token can be created via Settings / Developer Settings / Personal access tokens.

Save the token in ~/.gittoken and chmod the file 0600.

## Search / List Repositories

The following will show any repositories within the Organization that contains
the name "reponame".  The "show" command returns data in JSON.

```bash
git-code show "reponame"
```

## Clone Repositories

The "clone" command will match repositories that contain the "reponame".  If there
are multiple repositories that match, the names will be displayed and you will be
asked to make the search more specific.

```bash
# This will clone the repository into a directory name that matches the repository
# name
git-code clone "reponame"
# This will clone the repository into a directory named "JIRA-2112"
git-code clone "reponame" --directory "JIRA-2112"
```
