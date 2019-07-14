DATABASE_URL=oauth2_development
GO_DIR ?= $(shell pwd)
GO_PKG ?= $(shell go list -e -f "{{ .ImportPath }}")

TAG?=unknown

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html coverage.out

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
	GO111MODULE=on \
	go build -mod vendor -ldflags "-X $(GO_PKG)/cmd/version.appVersion=$(TAG)-$$(date -u +%Y%m%d%H%M)" -o "$(GO_DIR)/artifacts/bin" main.go

test: ## test application with race
	GO111MODULE=on \
	go test -mod vendor  -race -v ./...

createdb:
	createdb $(DATABASE_URL)

dropdb:
	dropdb $(DATABASE_URL)
