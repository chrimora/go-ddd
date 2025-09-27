# Go DDD

Simple DDD template for Go.

- Infra
  - Docker compose DB
  - Taskfile commands
  - Declarative migrations - Atlas

- Code
  - Environment config - env & dotenv
  - Contextual logging - slog
  - Open API gen & request validation - Huma
  - Dependency injection - fx
  - Repository implementation - sqlc

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
- `task worker`

Checkout http://localhost:8080/docs

## SQL

Using SQLC and Atlas to manage sql + migrations.

When creating a new table, add your sql schema into infrastructure/sql, then update:
- sqlc.yaml
- atlas.hcl

Then run:
- `task sqlgen`
- `task migrate`

## TODOs

- Worker
- Split service into commands/queries

