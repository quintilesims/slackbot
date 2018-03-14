SHELL:=/bin/bash
VERSION?=$(shell git describe --tags --always)
CURRENT_DOCKER_IMAGE=quintilesims/slackbot:$(VERSION)
LATEST_DOCKER_IMAGE=quintilesims/slackbot:latest

deps: 
	go get github.com/golang/mock/mockgen/model
	go install github.com/golang/mock/mockgen
	go get golang.org/x/tools/cmd/goimports
	go install golang.org/x/tools/cmd/goimports

mocks:
	mockgen -package mock github.com/quintilesims/slack SlackClient > mock/mock_slack_client.go
	cd mock && goimports -w . && cd -

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --ldflags "-X main.Version=$(VERSION)" -o slackbot . 
	docker build -t $(CURRENT_DOCKER_IMAGE) .

test:
	go test ./... -v

run:
	ngrok http 9090

release: build
	docker push $(CURRENT_DOCKER_IMAGE)
	docker tag  $(CURRENT_DOCKER_IMAGE) $(LATEST_DOCKER_IMAGE)
	docker push $(LATEST_DOCKER_IMAGE)

.PHONY: deps mocks build test run release
