# Simple Log Ingestor 

## Setup

### Install dependencies

```bash
go mod tidy
```

### Run

Use the make command to start the `server` and `rabbitmq-instance(docker)`

```bash
make run
```

### Ingest a Log using the HTTP Server

`Endpoint - http://localhost:9119/log`

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
    "level":"debug",
    "message":"Error while generating PDF",
    "from":"pdf-service-1"
    }' \
  http://localhost:9119/log
```

Response - `Log received!`
