# Compile stage
FROM golang:1.20.1-alpine3.17 AS build-env

WORKDIR /src
COPY go.mod /src/
COPY go.sum /src/
COPY src/*.go /src/

ENV GO111MODULE=on
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/go-redis/redis/v8
RUN go get -d -v gopkg.in/olahol/melody.v1

RUN go build -o /server *.go

# Serve stage
FROM alpine:3.13
WORKDIR /
COPY static/ /static/
COPY --from=build-env /server /
EXPOSE 5000
CMD ["/server"]