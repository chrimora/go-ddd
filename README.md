# Go Template

An opinionated template for API services in Go.

Doesn't include everything but holds the essentials.

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

And most importantly, a solid pattern to follow.  
It's pretty simple, a data model, routes, services, and schema & queries (repository).  
Following a domain first approach.

Keep service methods transactional!

Auth is excluded because it is often dictated by business requirements.

Data is updated using the read, modify, and update pattern, with optimistic locking (updated_at versioning).  
This is to help keep the code simple.

Extension (for growing projects):
- Wrap the sqlc generated pkg in a repository layer to separate db concerns
- If facing frequent race conditions, split the update methods (UpdateUserX, UpdateUserY, ...)
- As routes etc. grow, split the methods into separate files (i.e. in user/routes/) each of these files can have their own:
  - Struct for dependencies
  - Input & Output
  - Handler

TODO; testing framework
TODO; async worker

## Requirements

- Go
- Taskfile
- sqlc
- Atlas
- Docker compose
- yq

## Quick start

- `task up`
- `go run main.go`

Checkout http://localhost:8080/docs

