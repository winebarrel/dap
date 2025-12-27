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

.PHONY: docker/build
docker/build:
	docker build -t dap .

.PHONY: docker/run
docker/run: docker/build
	docker run dap -h
