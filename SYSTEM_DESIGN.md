# System Design Document: Hexagonal Architecture Overview

## 1. Architectural Pattern
This project follows the Hexagonal Architecture (Ports & Adapters), promoting separation of concerns and testability. The core domain logic is isolated from external systems (database, HTTP, etc.) via interfaces (ports) and adapters.

---

## 2. Main Components

### 2.1 Domain Layer (Core)
- **Entities/Models**: Defines business objects (e.g., `User` in internal/repository/models.go).
- **Services**: Implements business logic (e.g., `UserService` in internal/user/service.go).
- **Interfaces (Ports)**: Abstracts operations for adapters (e.g., `Querier` for data access).

### 2.2 Adapters Layer
- **Primary Adapters**: Handle inbound requests (e.g., HTTP handlers in internal/user/http.go).
- **Secondary Adapters**: Connect to external systems (e.g., database implementation in pkg/database/database.go).

### 2.3 Infrastructure Layer
- **Database**: Connection and health checks (pkg/database/database.go).
- **Repository**: SQL queries and data mapping (internal/repository/querier.go, db.go).

---

## 3. Key Interfaces & Flow

### User Registration Example
1. **HTTP Request** → `UserHandler.RegisterUser` (adapter)
2. **Service Call** → `UserService.RegisterUser` (domain)
3. **Repository Call** → `Querier.CreateUser` (port)
4. **Database Access** → SQL via `Queries` (adapter)

### Interface Definitions
- `UserService`: Business logic for users
- `UserHandler`: HTTP endpoints for users
- `Querier`: Data access abstraction
- `Database`: Database connection abstraction

---

## 4. File Structure Mapping
- **cmd/main.go**: Application entrypoint, wiring dependencies
- **internal/user/**: User domain logic & HTTP handlers
- **internal/repository/**: Data access interfaces, models, SQLC-generated code
- **pkg/database/**: Database connection logic
- **pkg/schema/**: SQL migration scripts

---

## 5. Data & Dependency Flow
- Adapters depend on domain interfaces, not implementations
- Domain logic is independent of infrastructure
- Dependency injection is used for wiring (see constructors like `NewUserService`, `NewUserHandler`)

---

## 6. Extensibility
- Add new features by defining domain interfaces and implementing adapters
- Swap infrastructure (e.g., database) by changing adapter implementations

---

## 7. Example: User Entity
```go
// internal/repository/models.go
 type User struct {
   ID    int64   `db:"id" json:"id"`
   Name  string  `db:"name" json:"name"`
   Email string  `db:"email" json:"email"`
   Bio   *string `db:"bio" json:"bio"`
 }
```

---

## 8. Summary Table
| Layer           | Example File/Type         | Purpose                  |
|-----------------|--------------------------|--------------------------|
| Domain/Core     | service.go, models.go     | Business logic, entities |
| Adapter         | http.go, database.go      | I/O, external systems    |
| Infrastructure  | db.go, querier.go         | Data access, persistence |

---

## 9. References
- Hexagonal Architecture: https://alistair.cockburn.us/hexagonal-architecture/
- SQLC: https://sqlc.dev/

---

**This document is structured for AI agents to quickly understand the system's architecture, interfaces, and extensibility points.**
