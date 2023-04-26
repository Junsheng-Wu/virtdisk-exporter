GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

build:
	CGO_ENABLED=0 go build -o ./bin/virtdisk-exporter .
.PHONY: docker-build
docker-build: build
docker-build:
	@docker build -t virtdisk-exporter:v0.1.0 . --network host
