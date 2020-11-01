FROM golang AS builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN make

FROM scratch
# Get SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/build/epitaf /app/epitaf

WORKDIR /app
ENTRYPOINT ["./epitaf"]
CMD ["start"]
