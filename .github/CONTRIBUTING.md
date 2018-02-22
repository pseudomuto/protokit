# Contributing to Protokit

First off, glad you're here and want to contribute! :heart:

## Getting Started

There are a number of dev tools required to work with Protokit. To make this a little cleaner and less invasive on your
global $GOPATH, we've elected to use [retool]. This also has the advantage of pinning things like `dep`, or
`protoc-gen-*` to specific versions to avoid inconsistencies between dev machines.

Running `make setup` will fetch [retool] (if it's not installed) and install vendored versions of dep and other tools in
the _tools/_ directory (git ignored). This should be sufficient to get you started.

## Submitting a PR

Here are some general guidelines for making PRs for this repo.

1. [Fork this repo](https://github.com/pseudomuto/protokit/fork)
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
* Imports are grouped into external, stdlib, internal groups in each file (see any go file in this repo for an example)
* Test are defined in `<package>_test` packages to ensure only the public interface is tested
* If you export something, make sure you add appropriate godoc comments

## Tagging a Release

* Set the `Version` in _version.go_
* Update CHANGELOG.md with the relevant changes
* `make release` - will commit everything, create a tag and push

[retool]: https://github.com/twitchtv/retool
