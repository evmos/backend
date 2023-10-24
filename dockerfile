FROM golang:1.21.3 as builder

WORKDIR /go/src/

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o ./server --ldflags '-w -s -extldflags "-static"' ./cmd/server/

FROM alpine:3.18 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch as runner

WORKDIR /app

COPY --from=builder /go/src/server .
COPY --from=builder /go/src/api/config/config.toml ./api/config/config.toml
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8081
CMD ["./server"]
