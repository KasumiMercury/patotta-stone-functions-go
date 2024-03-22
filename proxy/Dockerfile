FROM golang:1.21.6-alpine

RUN apk add --no-cache git bash

WORKDIR /app/proxy

CMD ["go", "run", "proxy.go"]
