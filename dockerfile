FROM golang:1.19.11-bullseye

WORKDIR /go/src/
COPY . .
RUN go build ./cmd/server/
EXPOSE 8081

CMD ["./server"]
