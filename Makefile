PROJECT=image_previewer
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	go build -ldflags="-s -w" -o '$(GOBIN)/$(PROJECT)' ./cmd/$(PROJECT)/main.go || exit

run-dev: build
	'$(GOBIN)/$(PROJECT)' --config='$(GOBASE)/configs/$(PROJECT).yml'

run:
	docker-compose up --build 

test:
	go test -race -count 100 ./...

test-integration:
	set -e ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml run integration_tests go test || test_status_code=$$? ;\
	docker-compose -f docker-compose.yml -f docker-compose.test.yml down ;\
	exit $$test_status_code ;

install-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: install-deps
	golangci-lint run ./...

install:
	go mod download

generate:
	go generate ./...

.PHONY: build
