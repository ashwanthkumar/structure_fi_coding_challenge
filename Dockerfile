FROM golang:1.17-alpine3.13

RUN apk add --no-cache git make

WORKDIR /app
COPY . /app/

# Delete the go mod cache location to reduce the final image size
RUN make && rm -rf ${GOPATH}/pkg/mod
EXPOSE 8080

CMD ["./structure_fi_coding_challenge"]