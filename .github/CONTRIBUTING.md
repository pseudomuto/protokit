# Contributing to Protokit

First off, glad you're here and want to contribute! :heart:

## Getting Started

This project uses [mise](https://mise.jdx.dev) to manage its toolchain (Go, golangci-lint, goreleaser, and Task) and [Task](https://taskfile.dev) as the command runner. Once mise is installed, it reads the pinned versions from `mise.toml` so everyone builds against the same tools.

1. [Install mise](https://mise.jdx.dev/getting-started.html)
1. Clone the repo
1. Run `mise install` to install the pinned tool versions
1. Run `task test` to make sure you're starting from a good place

Run `task --list` to see everything that's available.

## Submitting a PR

Here are some general guidelines for making PRs for this repo.

1. [Fork this repo](https://github.com/pseudomuto/protokit/fork)
1. Make a branch off of master (`git checkout -b <your_branch_name>`)
1. Make focused commits with descriptive messages
1. Add tests that fail without your code, and pass with it (`task test` is your friend)
1. GoFmt your code! (see <https://blog.golang.org/go-fmt-your-code> to setup your editor to do this for you)
1. Run `task lint` before pushing (`task lint:fix` fixes what it can automatically)
1. **Ping someone on the PR** (Lots of people, including myself, won't get a notification unless pinged directly)

Every PR should have a well detailed summary of the changes being made and the reasoning behind them. I've added a
PR template that should help with this.

## Code Guidelines

I don't want to be too dogmatic about this, but here are some general things I try to keep in mind:

* GoFmt all the things!
* Imports are grouped into external, stdlib, internal groups in each file (see any go file in this repo for an example) - really just use `goimports` and be done with it.
* Test are defined in `<package>_test` packages to ensure only the public interface is tested.
* If you export something, make sure you add appropriate godoc comments and tests.

## Cutting a Release

Releases are cut by manually running the Release workflow; there's no need to create a tag locally.

* Go to **Actions -> Release -> Run workflow**
* Choose the bump type (`patch`, `minor`, or `major`)

The workflow computes the next version from the latest `vX.Y.Z` tag, creates and pushes the tag, then runs goreleaser. Only maintainers with write access can trigger it.
