# Go DDD

Simple CQRS DDD template for Go.  
Designed for larger and scaling code bases.

A rest service + async event consumer.

- Infra
  - Docker compose DB
  - Declarative migrations - Atlas

- Code
  - Environment config - env & dotenv
  - Contextual logging - slog
  - Dependency injection - fx
  - Open API gen & request validation - Huma
  - Generated SQL code - sqlc
  - Transactional event outbox
  - Event pubsub - watermill

## Requirements

- Go
- sqlc
- yq
- Atlas
- docker compose
- Taskfile

## Quick start

- `task infra`
- `task migrate`
- `task server`
- `task consumer`

Checkout http://localhost:8080/docs

## TODO

- Testcontainers
- SQS pubsub - FIFO + MessageGroupId (AggregateId)

## Guides

### Auth

Auth is intentionally left out as this will depend on business requirements.
A dummy auth middleware has been added, ideally update this to validate a jwt or session cookie/token.

### SQL

Using SQLC and Atlas to manage sql + migrations.

When adding new schema and queries update:
- sqlc.yaml
- atlas.hcl

Then run:
- `task sqlgen`
- `task migrate`

Atlas migrations are declarative. There is no setup for data migrations.

