build:
	go build -o bin/pomegranate cmd/pomegranate/*.go

test:
	go test -v ./...
