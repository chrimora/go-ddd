# Go DDD

Simple DDD template for Go.

Vertical layered web service + async event consumer.

- Infra
  - Docker compose DB
  - Taskfile commands
  - Declarative migrations - Atlas

- Code
  - Environment config - env & dotenv
  - Contextual logging - slog
  - Dependency injection - fx
  - Open API gen & request validation - Huma
  - Repository implementation - sqlc
  - Transactional event outbox
  - Event pubsub - watermill

## Requirements

- Go
- Taskfile
- sqlc
- Atlas
- docker compose
- yq

## Quick start

- `task infra`
- `task server`
- `task consumer`

Checkout http://localhost:8080/docs

## SQL

Using SQLC and Atlas to manage sql + migrations.

When adding new schema and queries update:
- sqlc.yaml
- atlas.hcl

Then run:
- `task sqlgen`
- `task migrate`

# TODO

- CQRS
- SQS pubsub

