gbch
=======

[![Build Status](https://travis-ci.org/vvatanabe/gbch.svg?branch=master)](https://travis-ci.org/vvatanabe/gbch)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](http://www.opensource.org/licenses/mit-license.php)

## Description

Generate changelog from git history, tags and merged pull requests for Backlog

## Installation

    % go get github.com/vvatanabe/gbch/cmd/gbch/

## Synopsis

    % gbch -r /path/to/repo [--format markdown]

## Options

```
Options
-r, --repo=         git repository path (default: .)
-g, --git=          git path (default: git)
-f, --from=         git commit revision range start from
-t, --to=           git commit revision range end to
    --latest        output changes between latest two semantic versioned tags
    --apikey=       backlog api key
    --remote=       default remote name (default: origin)
-F, --format=       json or markdown
-A, --all           output all changes
-N, --next-version=
-w                  write result to file
    --show-uid      show the unique id on nulab account (Only user integrated with nulab account)
```

## Backlog API Key

Backlog's api key is required because this CLI depends on [Backlog API v2](https://developer.nulab.com/docs/backlog/), it is used in the following order of priority.

- command line option `--apikey`
- enviroment variable `BACKLOG_API_KEY`

## Requirements

git 1.8.5 or newer is required.

## Examples

### display all changes

    % gbch --format=markdown --next-version=v0.30.3 --all
    ...

### display changes between specified two revisions

    % gbch --from v0.9.0 --to v0.9.1
    ...

## Acknowledgments

gbch's origin is [github.com/Songmu/ghch](https://github.com/Songmu/ghch). I appreciate it a lot.

## Bugs and Feedback

For bugs, questions and discussions please use the GitHub Issues.

## Author

[vvatanabe](https://github.com/vvatanabe)
