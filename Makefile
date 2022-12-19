.PHONY: build
build:
	go build -v .

build-docker:
	docker build . -t weasyprintaas:latest --build-arg VERSION=latest-local

run:
	go run . serve

dep-upgrade:
	go get -u
	go mod tidy
