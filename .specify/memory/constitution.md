<!--
SYNC IMPACT REPORT
==================
Version Change: N/A → 1.0.0 (Initial ratification)
Modified Principles: N/A (All principles newly established)
Added Sections:
  - Core Principles (5 principles)
  - Technical Constraints
  - Development Workflow
  - Governance
Removed Sections: N/A

Templates Requiring Updates:
  ✅ .specify/templates/plan-template.md - Verified Constitution Check section exists
  ✅ .specify/templates/spec-template.md - Verified structure aligns with principles
  ✅ .specify/templates/tasks-template.md - Verified task categorization supports principles

Follow-up TODOs: None
-->

# GEMBOK Constitution

## Core Principles

### Clean Architecture (NON-NEGOTIABLE)

All code MUST follow strict Clean Architecture principles with dependency inversion. The architecture is organized in four concentric layers:

1. **Domain Layer** (`internal/domain/`): Entities and repository interfaces only. No external dependencies. Core business rules live here.
2. **Use Case Layer** (`internal/usecase/`): Application business logic. Depends only on domain interfaces. No framework imports.
3. **Interface Layer** (`internal/interface/`): HTTP handlers, DTOs, middleware. Depends on use cases. Framework-specific code isolated here.
4. **Infrastructure Layer** (`internal/infrastructure/`): External implementations (DB, MikroTik, GenieACS, Tripay). Implements repository interfaces from domain.

**Dependency Rule**: Dependencies MUST only point inward. Inner layers MUST NOT know anything about outer layers. Frameworks, databases, and external systems are pluggable details in outer layers.

**Rationale**: Ensures testability, maintainability, and flexibility to swap external implementations without affecting business logic.

### Golang Standards

All code MUST adhere to standard Golang project layout and conventions:

- **Directory Structure**: Follow `go-std-layout` guidelines: `cmd/`, `internal/`, `pkg/`, `configs/`, `scripts/`
- **Package Organization**: Packages MUST be named clearly and have a single responsibility
- **Interface Design**: Accept interfaces, return structs. Keep interfaces small and focused
- **Error Handling**: Explicit error checking. Never ignore errors. Use `errors.Is` and `errors.As` for error comparisons
- **Context**: Use `context.Context` for request-scoped values, cancellation, and deadlines throughout
- **Configuration**: Use `Viper` for configuration management. Support both YAML config files and environment variables
- **Logging**: Use `Zap` for structured logging. Log levels: Debug, Info, Warn, Error
- **Code Style**: Follow `gofmt`, `golint`, and `go vet` guidelines. No underscores in package names
- **Documentation**: Exported functions and types MUST have GoDoc comments

**Rationale**: Ensures code is idiomatic, maintainable, and familiar to Golang developers.

### Test-First

TDD (Test-Driven Development) is MANDATORY for all critical business logic:

1. **Write tests first** for all use cases and repository implementations
2. **Red-Green-Refactor** cycle strictly enforced:
   - Red: Write failing test
   - Green: Write minimal code to pass
   - Refactor: Improve implementation while keeping tests green
3. **Test Coverage**: Minimum 80% coverage for use cases and domain logic
4. **Unit Tests**: Test business logic in isolation using mocks for external dependencies
5. **Integration Tests**: Test database interactions, MikroTik API calls, and GenieACS integrations
6. **Test Organization**: Use `tests/` directory with `unit/`, `integration/`, and `contract/` subdirectories

**Rationale**: Catches bugs early, documents expected behavior, and enables confident refactoring.

### API First

All features MUST be designed with clear, versioned API contracts:

1. **RESTful Design**: Follow REST principles. Use proper HTTP methods and status codes
2. **API Versioning**: Include version in URL path (e.g., `/api/v1/`) from the start
3. **Request/Response DTOs**: Define clear DTOs in `internal/interface/http/dto/`. Separate from domain entities
4. **Validation**: Validate all input at the interface layer before reaching use cases
5. **Error Responses**: Consistent error format with status, code, message, and optional details
6. **OpenAPI/Swagger**: Document all endpoints with Swagger annotations for auto-generated documentation
7. **Authentication**: JWT-based auth for all protected endpoints. Include middleware for authentication and authorization

**Rationale**: Ensures API contracts are stable, documented, and can be consumed by frontend clients or third-party integrations.

### Security & Observability

Security and observability are NON-NEGOTIABLE requirements:

**Security**:
1. **Authentication**: JWT tokens with expiration. Validate tokens on every request to protected endpoints
2. **Authorization**: Role-based access control (RBAC). Admin, operator, and read-only roles
3. **Input Validation**: Validate and sanitize all user inputs. Prevent SQL injection and XSS
4. **Password Security**: Use bcrypt for password hashing with cost factor >= 10
5. **Sensitive Data**: Never log passwords, API keys, or JWT secrets. Use environment variables for secrets
6. **HTTPS**: Enforce HTTPS in production. Configure CORS properly

**Observability**:
1. **Structured Logging**: Log all significant events with request IDs, user IDs, and relevant context
2. **Error Tracking**: Capture stack traces for errors with appropriate logging level
3. **Metrics**: Track key metrics: request duration, error rates, active connections, database performance
4. **Health Checks**: Implement `/health` endpoint for system health monitoring
5. **Audit Logging**: Log all sensitive operations: user changes, router modifications, payment actions

**Rationale**: Ensures system is secure, debuggable, and meets operational requirements for production ISP management.

## Technical Constraints

### Technology Stack

- **Language**: Go 1.21 or higher
- **Web Framework**: Gin (lightweight, high-performance)
- **ORM**: GORM (MySQL driver)
- **Database**: MySQL 5.7 or higher
- **Authentication**: JWT (github.com/golang-jwt/jwt)
- **Configuration**: Viper
- **Logging**: Zap
- **Password Hashing**: bcrypt

### Performance Requirements

- **Response Time**: API responses MUST complete within 500ms p95 for standard operations
- **Concurrent Connections**: Support minimum 1000 concurrent connections
- **Database Connections**: Configure connection pooling (max_idle_conns: 10, max_open_conns: 100)
- **MikroTik API**: Optimize batch operations to minimize API calls

### Deployment Constraints

- **Platform**: Linux server (Ubuntu 20.04+ recommended)
- **Containerization**: Docker support required. Dockerfile provided at project root
- **Environment Variables**: Sensitive configuration via `.env` file (never commit)
- **Database Migrations**: Version-controlled migrations in `database/migrations/`

## Development Workflow

### Code Review Process

1. **Pull Required**: All changes MUST be reviewed via pull requests
2. **Automated Checks**: PRs must pass linting (`golint`, `gofmt`) and tests (`go test ./...`)
3. **Principle Compliance**: Reviewers must verify Clean Architecture and Golang Standards compliance
4. **Review Criteria**: Code quality, test coverage, documentation, adherence to principles

### Feature Development Flow

1. **Add entity** to `internal/domain/entities/`
2. **Define repository interface** in `internal/domain/repositories/`
3. **Implement repository** in `internal/infrastructure/repositories/`
4. **Create use case** in `internal/usecase/` (with tests written first)
5. **Create handler** in `internal/interface/http/handlers/`
6. **Add route** in `internal/interface/http/router.go`
7. **Write integration tests** for the complete flow
8. **Update Swagger documentation**

### Quality Gates

- **Linting**: `golangci-lint run` must pass without errors
- **Formatting**: `gofmt -s -w .` must show no changes
- **Testing**: `go test ./... -cover` must pass with coverage >= 80% for use cases
- **Build**: `go build ./...` must succeed
- **Documentation**: All exported functions must have GoDoc comments

## Governance

### Amendment Procedure

1. **Proposal**: Any team member can propose constitution amendments
2. **Review**: Proposal must be reviewed by at least 2 senior developers
3. **Approval**: Requires unanimous approval from core maintainers
4. **Migration Plan**: Major amendments must include a migration plan for existing code
5. **Documentation**: All amendments must be documented with version bump

### Versioning Policy

- **MAJOR**: Backward-incompatible changes to principles or architecture
- **MINOR**: New principles added or significant guidance expansions
- **PATCH**: Clarifications, wording improvements, non-semantic refinements

### Compliance Review

- **Pre-Commit**: Developers must ensure code complies with constitution before committing
- **Code Review**: PR reviewers must check compliance during review
- **Quarterly Audit**: Quarterly review of codebase to verify ongoing compliance
- **Violations**: Justified violations must be documented with rationale in `COMPLEXITY_TRACKING.md`

### Complexity Management

If constitution principles cannot be followed without compromising critical functionality:

1. Document the violation with clear rationale
2. Explain why simpler alternatives were rejected
3. Propose future refactoring to align with principles
4. Require approval from 2 senior developers

**Version**: 1.0.0 | **Ratified**: 2026-02-19 | **Last Amended**: 2026-02-19
