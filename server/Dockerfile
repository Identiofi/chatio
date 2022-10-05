# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o server .


FROM busybox 
## Copy the binary from the builder image
COPY --from=builder /app/server /app
## Copy the static files
COPY --from=builder /app/static /static

EXPOSE 8080

CMD /app