FROM golang AS builder
RUN go get -u github.com/swaggo/swag/cmd/swag

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . /app
RUN make

# FROM alpine
# RUN apk update
# RUN apk upgrade
# RUN apk add ca-certificates && update-ca-certificates
# # Change TimeZone
# RUN apk add --update tzdata
# ENV TZ=Europe/Paris
# # Clean APK cache
# RUN rm -rf /var/cache/apk/*
# # Get SSL certificates
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# # Retrieve binary
# COPY --from=builder /app/build/epitaf /app/epitaf

WORKDIR /app
ENTRYPOINT ["/app/epitaf"]
CMD ["start"]
