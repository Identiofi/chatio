# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go get -d -v golang.org/x/net/html \
    && CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM busybox
COPY --from=builder /app/app /app
CMD /app