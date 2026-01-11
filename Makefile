.PHONY: all
all: vet test build

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
vet:
	go test -v ./...

.PHONY: build
build:
	go build ./cmd/dap

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docker/build
docker/build:
	docker build -t dap .

.PHONY: docker/run
docker/run: docker/build
	docker run dap -h
