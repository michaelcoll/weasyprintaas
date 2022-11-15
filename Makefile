.PHONY: build
build:
	go build -v .

run:
	go run . serve

dep-upgrade:
	go get -u
	go mod tidy
