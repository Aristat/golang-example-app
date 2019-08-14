GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")

GOOS ?= $(shell go env GOOS || echo linux)
GOARCH ?= $(shell go env GOARCH || echo amd64)
CGO_ENABLED ?= 0

DOCKER_IMAGE ?= unknown
TAG ?= unknown
CACHE_TAG ?= unknown_cache

DATABASE_URL ?= oauth2_development

define build_resources
 	find "$(GO_DIR)/resources" -maxdepth 1 -mindepth 1 -exec cp -R -f {} $(GO_DIR)/artifacts/${1} \;
endef

init:
	mkdir -p generated artifacts ;\
    rm -rf artifacts/*

vendor:
	rm -rf $(GO_DIR)/vendor ;\
	GO111MODULE=on \
	go mod vendor

build: init
	$(call build_resources) ;\
	GO111MODULE=on GOOS=${GOOS} CGO_ENABLED=${CGO_ENABLED} GOARCH=${GOARCH} \
	go build -mod vendor -ldflags "-X $(GO_PKG)/cmd/version.appVersion=$(TAG)-$$(date -u +%Y%m%d%H%M)" -o "$(GO_DIR)/artifacts/bin" main.go

docker-image: ## build docker image
	docker rmi ${DOCKER_IMAGE}:${TAG} || true ;\
	docker build --cache-from ${DOCKER_IMAGE}:${CACHE_TAG} -f "${GO_DIR}/docker/app/Dockerfile" -t ${DOCKER_IMAGE}:${TAG} ${GO_DIR}

test: ## test application with race
	GO111MODULE=on \
	go test -mod vendor  -race -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out

createdb:
	createdb $(DATABASE_URL)

dropdb:
	dropdb $(DATABASE_URL)
