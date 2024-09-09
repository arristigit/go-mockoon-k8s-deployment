FROM golang:1.17-bullseye as build

WORKDIR /app
ADD . /app

# RUN go get -d -v ./... # Not required step
RUN go mod init temp-module
RUN go build -o /go/bin/app go_code.go

# FROM gcr.io/distroless/base-debian11
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y bash

COPY --from=build /go/bin/app /

EXPOSE 8090
CMD ["/app"]