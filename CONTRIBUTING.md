# Contributing

Thanks for considering a contribution to `pleco-api`.

## Before You Start

- Open an issue first for large changes, architectural changes, or new modules.
- Keep pull requests focused. Small, reviewable PRs are preferred over broad mixed changes.
- Make sure secrets, real credentials, and local environment files are never committed.

## Local Setup

1. Copy the local environment file:

```bash
cp .env.example .env
```

2. Run migrations and seed data:

```bash
go run ./cmd/migrate
go run ./cmd/seed
```

3. Start the app:

```bash
go run ./cmd/api
```

Or use Docker:

```bash
cp .env.docker.example .env.docker
make docker-up
```

## Development Guidelines

- Follow the existing modular structure under `internal/modules/`.
- Keep API responses consistent with the project response envelope.
- Prefer clear, explicit dependency wiring over hidden globals.
- Add or update tests when changing behavior.
- Update README, OpenAPI, and Postman assets when public API behavior changes.

## Module Structure

New business domains should live under `internal/modules/<domain>/` and follow the existing handler-service-repository shape. Keep module code focused on one domain so auth, users, roles, permissions, tokens, social login, audit, and future modules can evolve independently.

Expected files for a typical module:

```text
internal/modules/<domain>/
├── dto.go          # request and response DTOs
├── handler.go      # HTTP handlers and response mapping
├── model.go        # GORM/domain models
├── module.go       # dependency wiring for the module
├── repository.go   # database access behind an interface
├── routes.go       # route registration
└── service.go      # business logic
```

Small modules may omit files they do not need, but keep the same ownership boundaries:

- Handlers parse HTTP input and return envelope responses.
- Services own business rules, transactions, cache invalidation, and cross-module orchestration.
- Repositories own persistence queries and expose interfaces for testing.
- DTOs keep public API shapes separate from database models.
- Routes register endpoints and middleware close to the module they belong to.

When adding a module, wire it from `internal/appsetup/`, add migrations when schema changes, add tests for meaningful behavior, and update README, OpenAPI, and Postman assets if the public API changes.

## Quality Checks

Run the baseline checks before opening a PR:

```bash
go test ./...
```

Optional repository integration tests can be run with:

```bash
TEST_DATABASE_URL="postgresql://..." go test ./tests -run 'TestPermissionRepository|TestAuditRepository' -v
```

## Pull Request Checklist

- Code compiles and tests pass
- No secrets are committed
- API/docs are updated when needed
- Migration and seed changes are included when schema changes
- PR description explains the problem, change, and impact
