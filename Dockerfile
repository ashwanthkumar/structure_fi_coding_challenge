# syntax=docker/dockerfile:1

FROM golang:1.17-alpine3.13
EXPOSE 8080

WORKDIR /build
COPY . /build/

RUN apk add --no-cache git make gcc musl-dev bash && \
      make && \
      mkdir -p /app && \
      cp structure_fi_coding_challenge /app/ && \
      # Delete the go mod cache and custom packages to reduce the final image size
      rm -rf /build && \
      rm -rf ${GOPATH}/pkg/mod && \
      apk del git make gcc musl-dev bash

CMD ["/app/structure_fi_coding_challenge"]