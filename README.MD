# Go JWT Authentication & Clean Architecture API

A production-style user authentication microservice built in Go featuring bcrypt password hashing, stateless HS256 JWT tokens, context-injected auth middleware, and interface-driven repository decoupling.

---

### 🛠️ Architecture & Go Concepts Applied

*   **Clean Architecture via Interfaces (`UserStore`):** Decouples HTTP handlers (`AuthServer`) from storage implementations (`PostgresStore` for production, `FakeStore` for fast unit testing).
*   **Adaptive Password Hashing (`bcrypt`):** Salted and hashed passwords before database storage, hiding hashes from JSON serialization via `json:"-"`.
*   **JWT Middleware & Context (`context.WithValue`):** Validates `Authorization: Bearer <token>` headers and injects authenticated user claims down the request lifecycle.
*   **HTTP Handler Testing (`net/http/httptest`):** Simulates full registration and login HTTP cycles in memory without external database dependencies.

---

### 🔌 API Endpoints & Testing

```bash
# Start the authentication API
go run .

# Run interface-mocked HTTP unit tests
go test -v
```
