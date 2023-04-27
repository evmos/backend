FROM golang:1.19.0-bullseye

WORKDIR /go/src/
COPY . .
RUN go build
EXPOSE 8081
CMD ["./dashboard-backend"]
