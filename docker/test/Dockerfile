FROM golang AS builder
RUN apt update && apt install postgresql-client -y

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
ENTRYPOINT ["/app/docker/test/entrypoint.sh"]