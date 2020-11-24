FROM golang AS builder
RUN go get -u github.com/swaggo/swag/cmd/swag

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN make swag
RUN make test