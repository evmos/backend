## Server Binary

### Requirements

In order to run the server you must first set the following environment variables:

- `NUMIA_API_KEY` - Required
- `NUMIA_RPC_ENDPOINT` - Required
- `RPC_SERVER_PORT` - optional

### Build

To build run:

```
  go build ./cmd/server/
```

### Install

To install the binary run:

```
  go install ./cmd/server/
```

### Run

To run and test locally run:

```
  go run ./cmd/server/
```

Note: Do not forget to set env variables first
