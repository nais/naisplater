.PHONY: all test alpine

all:
	go build -o bin/naisplater cmd/naisplater/main.go

test:
	go test -v -count=1 ./...

alpine:
	go build -a -installsuffix cgo -o bin/naisplater cmd/naisplater/main.go
