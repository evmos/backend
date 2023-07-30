FROM golang:1.20.5-bullseye

WORKDIR /go/src/
COPY . .
RUN go build ./cmd/server/
EXPOSE 8081

CMD ["./server"]
