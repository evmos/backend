FROM golang:1.21.1-alpine3.18

WORKDIR /go/src/
COPY . .
RUN go build ./cmd/server/
EXPOSE 8081

CMD ["./server"]
