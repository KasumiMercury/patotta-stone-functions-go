FROM golang:1.21.6-alpine

RUN apk add --no-cache git bash

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0

CMD ["air", "-c", ".air.toml"]
