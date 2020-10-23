FROM golang:alpine AS builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o epitaf

FROM scratch
COPY --from=builder /app/epitaf /app/epitaf
WORKDIR /app
ENTRYPOINT ["./epitaf"]
