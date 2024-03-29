.PHONY: build
.DEFAULT_GOAL := build

TIMESTAMP=$(shell date +%s)
VERSION=0.1.0

# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63#goarch-values
GOARCH?=amd64
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63#goos-values
GOOS?=linux

build: swag
	env GOARCH=$(GOARCH) GOOS=$(GOOS) go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.Build=$(TIMESTAMP)"

swag:
	swag init

run: swag
	gow -v run *.go

docker: swag
	docker build -t s3:$(VERSION)-alpine3.18 .
