COMMIT = $(shell git rev-parse --short HEAD)
FLAGS="-s -w -X 'main.COMMIT=${COMMIT}'"

build:
	go build -ldflags ${FLAGS}

test:
	@go test ./...

clean:
	@go clean ./...
