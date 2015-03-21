export GOPATH=$(shell pwd)

test:
	go test -v .

format:
	gofmt -l -w $(wildcard *.go)

.PHONY: test
