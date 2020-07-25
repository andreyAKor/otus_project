PROJECT=image_previewer
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/$(PROJECT)' ./cmd/$(PROJECT)/main.go || exit

run-dev: build
	'$(GOBIN)/$(PROJECT)' --config='$(GOBASE)/configs/$(PROJECT).yml'

run:
	docker-compose up

test:
	go test -race -count 100 ./...

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

install:
	go mod download

generate:
	go generate ./...

.PHONY: build
