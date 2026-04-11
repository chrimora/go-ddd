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

- `cp .env.local .env`
- `task install`
- `task infra`
- `task migrate`
- `task run`

Checkout http://localhost:8080/docs

## TODO

- Testcontainers
- SQS pubsub - FIFO + MessageGroupId (AggregateId)

## Guides

### Adding an aggregate

The `user` package is an example bounded context containing a single aggregate. A bounded context can hold many aggregates — e.g. an `order` context might contain `Order`, `OrderLine`, and `Invoice` aggregates, each with their own repository, commands, queries, and events, all wired together in a shared `module.go`.

To add a new aggregate (within a new or existing bounded context), create the following:

```
internal/<context>/
  domain/
    aggregate.go       # struct embedding AggregateRoot, constructor, mutators
    events.go          # event types and constructors
    repository.go      # repository interface + SQL implementation
  application/
    commands/          # one file per command (create, update, delete)
    queries/           # one file per query
    eventhandlers/     # one file per event handler
  infrastructure/sql/
    schema.sql         # table definition
    repository.sql     # sqlc write queries (create, update, delete)
    queries.sql        # sqlc read queries
  interfaces/rest/
    routes.go          # Huma route registration
  module.go            # fx CoreModule, APIModule, ConsumerModule
  test/
    module.go          # UnitTestModule + IntegrationTestModule
    factory.go         # Mock factory using RehydrateX + repo.Create
```

Wire the new context into `cmd/api/main.go` and `cmd/consumer/main.go`, then run `task migrate`.

### Auth

Auth is intentionally left out as this will depend on business requirements.
A dummy auth middleware has been added, ideally update this to validate a jwt or session cookie/token.

### Testing

The tests are split into unit and integration tests defined by the build tags:
- `//go:build unit`
- `//go:build integration`

- `task test-unit`
- `task test-integration`

Unit tests make use of mock depependencies (in memory db), integration tests run
serially and require `task infra` to be up.

Factories create and clean up test data via `t.Cleanup`:
```go
user := s.uf.Mock(t, ctx)                                        // default fields
user := s.uf.Mock(t, ctx, map[string]any{"Name": "Hello World"}) // overrides
```

### SQL

Using SQLC and Atlas to manage sql + migrations.

When adding new schema and queries update:
- sqlc.yaml
- atlas.hcl

Then run:
- `task sqlgen`
- `task migrate`

yq is used for anchor expansion (sqlc does not suport this).
Atlas migrations are declarative. There is no setup for data migrations.

