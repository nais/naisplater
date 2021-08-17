.PHONY: all install test alpine

all:
	go build -o bin/naisplater cmd/naisplater/main.go

install: all
	sudo install bin/naisplater /usr/local/bin

test:
	go test -v -count=1 ./...

alpine:
	go build -a -installsuffix cgo -o bin/naisplater cmd/naisplater/main.go
