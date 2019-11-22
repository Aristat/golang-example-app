GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")

GOOS ?= $(shell go env GOOS || echo linux)
GOARCH ?= $(shell go env GOARCH || echo amd64)
CGO_ENABLED ?= 0

DOCKER_IMAGE ?= unknown
TAG ?= unknown
CACHE_TAG ?= unknown_cache

DATABASE_URL ?= golang_example_development

REMOVE_CONTAINERS ?= OFF

define build_resources
 	find "$(GO_DIR)/resources" -maxdepth 1 -mindepth 1 -exec cp -R -f {} $(GO_DIR)/artifacts/${1} \;
endef

install: init ## install cli tools
	export GO111MODULE=off ;\
    go get -v github.com/rubenv/sql-migrate/... ;\
    go get -u github.com/google/wire/cmd/wire ;\
    go get -u github.com/vektah/dataloaden ;

init: ## init packages
	mkdir -p artifacts ;\
    rm -rf artifacts/*

start: ## start daemon on development mode
	./artifacts/bin daemon -c ./artifacts/configs/development.yaml -d

vendor: ## generate vendor
	rm -rf $(GO_DIR)/vendor ;\
	GO111MODULE=on \
	go mod vendor

gqlgen-generate: ## generate graphql server
	go run github.com/99designs/gqlgen

prototool-generate: ## generate proto file
	prototool generate resources/proto

build: init ## build binary file
	$(call build_resources) ;\
	GO111MODULE=on GOOS=${GOOS} CGO_ENABLED=${CGO_ENABLED} GOARCH=${GOARCH} \
	go build -mod vendor -ldflags "-X $(GO_PKG)/cmd/version.appVersion=$(TAG)-$$(date -u +%Y%m%d%H%M)" -o "$(GO_DIR)/artifacts/bin" main.go

docker-image: ## build docker image
	REMOVE_CONTAINERS=${REMOVE_CONTAINERS} DOCKER_IMAGE=${DOCKER_IMAGE} ./scripts/remove_docker_containers.sh
	docker rmi ${DOCKER_IMAGE}:${TAG} || true ;\
	docker build --cache-from ${DOCKER_IMAGE}:${CACHE_TAG} -f "${GO_DIR}/docker/app/Dockerfile" -t ${DOCKER_IMAGE}:${TAG} ${GO_DIR}

test: ## test application with race
	GO111MODULE=on \
	go test -mod vendor  -race -v ./...

coverage: ## test coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out

createdb: ## create database
	createdb $(DATABASE_URL)

dropdb: ## drop database
	dropdb $(DATABASE_URL)

.PHONY: install init vendor gqlgen-generate prototool-generate

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
