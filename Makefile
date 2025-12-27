.PHONY: all
all: vet build

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build ./cmd/dap

.PHONY: lint
lint:
	golangci-lint run
