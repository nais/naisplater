FROM golang:alpine AS go-getter

RUN apk upgrade --update && \
    apk add bash git

RUN go get github.com/tsg/gotpl

FROM alpine:3.7

RUN apk upgrade --update && \
    apk add bash curl

ENV KUBECTL_VERSION="1.9.2"
ENV YQ_VERSION="1.14.0"

RUN curl -Ls https://github.com/mikefarah/yq/releases/download/${YQ_VERSION}/yq_linux_amd64 > /usr/bin/yq && \
    chmod +x /usr/bin/yq

COPY --from=go-getter /go/bin/gotpl /usr/bin/gotpl
COPY naisplater /usr/bin/

CMD bash

WORKDIR /root
