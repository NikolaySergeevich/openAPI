PKG_LIST := $(shell go list ./... | grep -v /vendor/)
PATH := $(PATH):$(GOPATH)/bin


.PHONY: build
build:
	go build -o bin/objects-srv cmd/objects-srv/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: lint
lint:
	golangci-lint run --timeout 5m -v ./...

.PHONY: genid
genid:
	go run cmd/genid/main.go

.PHONY: generate
generate:
	go generate ./...

.PHONY: install
install:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: test
test:
	@$(GO_TEST)
