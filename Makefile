.PHONY: all
all: vet build

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build ./cmd/dap

.PHONY: lint
lint:
	golangci-lint run
