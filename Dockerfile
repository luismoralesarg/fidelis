# Build stage
FROM golang:alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o kvs

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/kvs .

CMD ["./kvs"]