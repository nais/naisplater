SHELL   := bash
VERSION := $(shell cat ./version)
NAME    := naisplater
IMAGE   := ghcr.io/nais/${NAME}:${VERSION}
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: test bump build push

test:
	rm -rf ./test/out/*
	/bin/bash naisplater dev ./test/templates ./test/vars ./test/out supersecret69 > /dev/null
	diff ./test/out ./test/expected && echo "OK" || (echo "FAILED" && exit 1)

docker-test:
	docker run -v ${ROOT_DIR}/test/templates:/templates -v ${ROOT_DIR}/test/vars:/vars -v ${ROOT_DIR}/test/out:/out --rm ${IMAGE} naisplater dev /templates /vars /out supersecret69
	diff ./test/out ./test/expected && echo "OK" || (echo "FAILED" && exit 1)

bump:
	/bin/bash bump.sh

build:
	docker image build -t ${IMAGE} .

push:
	docker image push ${IMAGE}

