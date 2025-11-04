FROM golang:1.25.3 AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

EXPOSE 50051

ENV DB_HOST="host.docker.internal:3306"

CMD ["/build/server"]
