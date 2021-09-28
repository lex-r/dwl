.PHONY:all
all: build

.PHONY:build
build:
	go build -v ./...

.PHONY:test
test:
	go test -v ./...

.PHONY:deps
deps:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY:mockgen
mockgen:
	mockgen -source=./downloader/downloader.go -destination=./downloader/downloader_mock.go -package=downloader
