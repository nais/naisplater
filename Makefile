.PHONY: test alpine bump

test:
	go test -v -count=1 ./...

alpine:
	go build -a -installsuffix cgo -o bin/naisplater cmd/naisplater/main.go

bump:
	/bin/bash bump.sh
