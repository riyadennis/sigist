# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Event management system consisting of two microservices communicating via Kafka:
- **graphql-service**: GraphQL API for user feedback management (Go 1.24+, gqlgen, PostgreSQL)
- **rest-service**: REST API for email management (Go 1.19, SQLite)
- **benthos**: Message broker that consumes Kafka events and forwards to rest-service

Data flows from graphql-service → Kafka → benthos → rest-service.

## Service Architecture

### graphql-service (Port 4000)
- **Stack**: gqlgen GraphQL server with federation support, PostgreSQL, Kafka producer
- **Entry point**: `graphql-service/main.go`
- **GraphQL schema**: `graphql-service/graph/schema.graphqls`
- **Resolvers**: `graphql-service/graph/schema.resolvers.go`
- **Database**: PostgreSQL with pgx driver, migrations in `graphql-service/migrations/`
- **Kafka**: Produces UserFeedback events to `data-pipe` topic via `internal/kafka.go`
- **Config**: Environment variables loaded via `internal/config.go`

Key components:
- `service/service.go`: Service initialization and lifecycle
- `service/postgres.go`: Database operations
- `graph/resolver.go`: GraphQL resolver root
- `internal/kafka.go`: Kafka producer implementation

### rest-service (Port 8080)
- **Stack**: Chi router, SQLite, structured logging with zap
- **Entry point**: `rest-service/main.go`
- **Routes**: `/email` (POST), `/emails` (GET)
- **Database**: SQLite with golang-migrate, migrations in `rest-service/migrations/`
- **Handler**: `service/http.go` contains email endpoints
- **Config**: Environment variables via `internal/config.go`

Note: The rest-service go.mod references the old module name `github.com/riyadennis/sigist/rest-service` instead of `github.com/riyadennis/event-management/rest-service`.

### benthos
- **Config**: `environment/benthos/kafka-consumer.yaml`
- Consumes from Kafka topics `data-pipe` on brokers `kafka:19092` and `localhost:9092`
- Posts messages to `http://rest-service:8080/email`
- Consumer group: `data-pipe`

## Development Commands

### Running the Full Stack
```bash
make docker-run
# Runs docker-compose up with all services: postgres, kafka, kafka-ui, graphql-service, rest-service, benthos
```

### Cleanup
```bash
make docker-clean
# Stops all containers and removes images
```

### Debugging
```bash
make check-logs
# View Kafka broker logs

make ssh-rest-service
# SSH into rest-service container
```

### Access Points
- GraphQL service: http://localhost:4000
- REST service: http://localhost:8080
- Kafka UI: http://localhost:9019
- Kafka broker: localhost:9092
- PostgreSQL: localhost:5434 (user: username, password: password, db: feedback)

## Running Tests

### Unit Tests
```bash
cd graphql-service
go test ./...

cd rest-service
go test ./...
```

### Integration Tests (BDD with Godog)
```bash
cd graphql-service/integration
go test -v
# Runs Cucumber/Godog scenarios from features/ directory
```

### Contract Tests (Pact)
```bash
cd graphql-service/pacts
go test -v

cd rest-service/pacts
go test -v
```

### Single Test
```bash
go test -run TestName ./path/to/package
```

## GraphQL Code Generation

When modifying GraphQL schema:
```bash
cd graphql-service
go run github.com/99designs/gqlgen generate
```

This regenerates:
- `graph/generated/generated.go`
- `graph/generated/federation.go`
- `graph/model/models_gen.go`

Configuration in `graphql-service/gqlgen.yaml`.

## Database Migrations

### graphql-service (PostgreSQL)
Migrations in `graphql-service/migrations/` using naming: `000001_name.up.sql` and `000001_name.down.sql`

### rest-service (SQLite)
Migrations in `rest-service/migrations/` with same naming convention

Migrations run automatically on service startup via golang-migrate.

## Docker Compose Network

All services communicate via the `datapipe` bridge network:
- Kafka advertised as `kafka:19092` internally, `localhost:9092` externally
- Services reference each other by container name (e.g., `rest-service:8080`)

## Module Naming Inconsistency

The rest-service go.mod still references `github.com/riyadennis/sigist/rest-service` instead of `github.com/riyadennis/event-management/rest-service`. This affects imports in rest-service code.

## Testing Infrastructure

- Integration tests use gnomock for containerized Kafka
- Unit tests use go-sqlmock for database mocking
- Contract tests use pact-foundation/pact-go
- Logging in tests uses zap with test configuration