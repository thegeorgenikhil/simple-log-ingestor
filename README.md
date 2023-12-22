# Simple Log Ingestor 

## Some notes

Q: *What happens if there is no queue bound to the exchange? ie. Publisher publishes to an exchange but no consumer is alive thus no queue is bound to the exchange.*

That is the way RabbitMQ is designed - publishers publish to exchanges, not queues.

If there is no queue bound (with a matching routing key if the exchange requires one), the message is simply discarded.

You can enable `publisher returns` and set the `mandatory` flag when publishing and the broker will return the message (but it arrives on a different thread, not the publishing thread).

---

## Setup

### Install dependencies

```bash
go mod tidy
```

### Run

Use the make command to start the `server` and `rabbitmq-instance(docker)`

```bash
make start
```

### Ingest a Log using the HTTP Server

`Endpoint - http://localhost:9119/log`

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
    "level":"error",
    "message":"Error while generating PDF",
    "from":"pdf-service-1"
    }' \
  http://localhost:9119/log
```

Response - `Log received!`

