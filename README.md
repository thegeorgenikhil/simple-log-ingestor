# Simple Log Ingestor

This is a simple log ingestor that uses RabbitMQ as a message broker. The log ingestor is a HTTP server that accepts logs and publishes them to a RabbitMQ exchange. The logs can be of different levels - `info`, `error` and `debug`. The logs are published to the exchange with the routing key as the log level. The exchange is then subscribed to by two consumers - one `logger_service` and one `alert_service`.

The `logger_service` consumes all types of log levels (routing key - `info`, `error` and `debug`) and writes them to the log file named `all_logs.log`.

The `alert_service` consumes only logs with `error` as the log level (routing key - `error`) and writes them to the log file named `mail_logs.log`. The idea with the `mail_logs.log` is that we are considering the log_file as record of mails sent. Like a proxy for a mail server.

The whole idea with this repo is to learn RabbitMQ, how queues, exchanges and bindings work by building a log ingestor and alerting system.

## Some notes

Q: _What happens if there is no queue bound to the exchange? ie. Publisher publishes to an exchange but no consumer is alive thus no queue is bound to the exchange._

That is the way RabbitMQ is designed - publishers publish to exchanges, not queues.

If there is no queue bound (with a matching routing key if the exchange requires one), the message is simply discarded.

You can enable `publisher returns` and set the `mandatory` flag when publishing and the broker will return the message (but it arrives on a different thread, not the publishing thread).

---

## Setup

### Install dependencies

```bash
go mod tidy
```

### Get the email credentials for the alert service

Copy the `.env.example` file to `.env` and fill in the details.

Get the credentials from [Brevo](https://app.brevo.com/settings/keys/smtp) (300 free emails per day)

```
EMAIL_SERVER_HOST=
EMAIL_SERVER_PORT=
EMAIL_SERVER_USERNAME=
EMAIL_SERVER_PASSWORD=
EMAIL_FROM_ADDRESS=<can be anything>
EMAIL_TO_ADDRESS=<your own email address or the one you want to send the alert to>
```

### Run

Use the make command to start the `server` and `rabbitmq-instance(docker)`

```bash
# Start the rabbitmq instance
make start_mq

# Start the server in a new terminal
make start_server

# Start the logger service in a new terminal
make start_logger_service

# Start the alert service in a new terminal
make start_alert_service
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
