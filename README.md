# Simple Log Ingestor 

This is a simple log ingestor that uses RabbitMQ as a message broker. The log ingestor is a HTTP server that accepts logs and publishes them to a RabbitMQ exchange. The logs can be of different levels - `info`, `error` and `debug`. The logs are published to the exchange with the routing key as the log level. The exchange is then subscribed to by two consumers - one `log_consumer` and one `alert_consumer`. 

The `log_consumer` consumes all types of log levels (routing key - `info`, `error` and `debug`) and writes them to the log file named `all_logs.log`. 

The `alert_consumer` consumes only logs with `error` as the log level (routing key - `error`) and writes them to the log file named `mail_logs.log`. The idea with the `mail_logs.log` is that we are considering the log_file as record of mails sent. Like a proxy for a mail server.

The whole idea with this repo is to learn RabbitMQ, how queues, exchanges and bindings work by building a log ingestor and alerting system.

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
# Start the rabbitmq instance
make start_mq

# Start the server in a new terminal
make start_server

# Start the consumer in a new terminal
make start_log_consumer
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

