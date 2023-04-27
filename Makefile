.PHONY: backend

test:
	@go test -v ./...

stest:
	@go test ./...

build:
	@go build

run:
	@go build && ./dashboard-backend

lint:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0 --config .golangci.yml --color always ./...