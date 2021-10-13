build:
	go build -o bin/pomegranate cmd/pomegranate/*.go

run:
	export $(grep -v '^#' .env | xargs)
	go run cmd/pomegranate/*.go

test:
	go test -v ./...
