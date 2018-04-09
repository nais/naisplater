SHELL   := bash
VERSION := $(shell cat ./version)
NAME    := navikt/naisplater
IMAGE   := ${NAME}:${VERSION}

.PHONY: test bump build push

test:
	/bin/bash test/run

bump:
	/bin/bash bump.sh

build:
	docker image build -t ${NAME} -t ${IMAGE} .

push:
	docker image push ${IMAGE}

