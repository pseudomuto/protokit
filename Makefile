.PHONY: setup test

setup:
	$(info Synching dev tools and dependencies...)
	@if test -z $(which retool); then go get github.com/twitchtv/retool; fi
	@retool sync
	@retool do dep ensure

fixtures/fileset.pb: fixtures/*.proto
	$(info Generating fixtures...)
	@cd fixtures && go generate

test: fixtures/fileset.pb
	@go test -race -cover ./

test-ci: fixtures/fileset.pb
	@retool do goverage -race -coverprofile=coverage.txt -covermode=atomic ./
