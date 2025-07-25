# Contributing to Protokit

First off, glad you're here and want to contribute! :heart:

## Getting Started

It's always good to start simple. Clone the repo and `make test` to make sure you're starting from a good place.

## Submitting a PR

Here are some general guidelines for making PRs for this repo.

1. [Fork this repo](https://github.com/Djarvur/protokit/fork)
1. Make a branch off of master (`git checkout -b <your_branch_name>`)
1. Make focused commits with descriptive messages
1. Add tests that fail without your code, and pass with it (`make test` is your friend)
1. GoFmt your code! (see <https://blog.golang.org/go-fmt-your-code> to setup your editor to do this for you)
1. **Ping someone on the PR** (Lots of people, including myself, won't get a notification unless pinged directly)

Every PR should have a well detailed summary of the changes being made and the reasoning behind them. I've added a
PR template that should help with this.

## Code Guidelines

I don't want to be too dogmatic about this, but here are some general things I try to keep in mind:

* GoFmt all the things!
* Imports are grouped into external, stdlib, internal groups in each file (see any go file in this repo for an example) - really just use `goimports` and be done with it.
* Test are defined in `<package>_test` packages to ensure only the public interface is tested.
* If you export something, make sure you add appropriate godoc comments and tests.

## Tagging a Release

* Ensure you're on a clean master
* Run `make release`
