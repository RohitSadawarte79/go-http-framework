# 🚀 Go HTTP Framework

A custom HTTP framework built **from scratch** in Go using only the standard library (`net/http`). No external dependencies.

## Why From Scratch?

I was curious about what actually happens when a request hits a server. I didn't want to just call `router.GET()` and move on — I wanted to understand how routing, middleware, and request handling work at the lowest level. So I built it all myself using only Go's standard library.

## Features

### Custom Router
- **Path parameter extraction** — `/users/:id` extracts `id` from the URL
- **Static route priority** — `/users/profile` matches before `/users/:id`
- **REST semantics** — proper `404 Not Found` vs `405 Method Not Allowed`
- **Conflict detection** — panics at registration time if routes conflict

### Middleware Stack
- **Chain** — compose multiple middlewares with correct execution order
- **Recovery** — catches panics, returns 500, keeps the server alive
- **Logger** — logs request method, path, and completion
- **RequestID** — generates unique ID per request, adds to context and response headers
- **CORS** — configurable origin whitelist with preflight (OPTIONS) handling

### Testing
- 17 tests covering routing, middleware, and full-stack integration
- Tests for edge cases: trailing slashes, empty params, method conflicts

## Architecture

```
Request → CORS → Recovery → Logger → RequestID → Router → Handler → Response
```

Each middleware wraps the next handler using the signature:
```go
type Middleware func(http.Handler) http.Handler
```

## Quick Start

```bash
git clone https://github.com/RohitSadawarte79/go-http-server.git
cd go-http-server
go run .
```

Server starts on `http://localhost:8080`.

### Example Endpoints

```bash
# Get all users
curl http://localhost:8080/user

# Get user by ID (path parameter)
curl http://localhost:8080/user/42

# Create a user
curl -X POST http://localhost:8080/user \
  -H "Content-Type: application/json" \
  -d '{"first_name": "Rohit", "last_name": "S", "age": 21}'

# Test CORS preflight
curl -X OPTIONS http://localhost:8080/user \
  -H "Origin: http://localhost:3000" -I
```

### Run Tests

```bash
go test -v ./...
```

## Tech Stack

- **Language:** Go 1.23+
- **Dependencies:** None (standard library only)
- **Testing:** `net/http/httptest`

## What I Learned

Building this taught me that frameworks are just organized patterns on top of simple primitives. The entire middleware system is based on one idea: a function that takes a handler and returns a handler. Once you understand that, everything else is composition.

## Roadmap

- [x] Phase 1 — `net/http` internals
- [x] Phase 2 — Router + Middleware
- [ ] Phase 3 — Clean Architecture (handlers → services → repositories)
- [ ] Phase 4 — Database Integration (PostgreSQL)
- [ ] Phase 5 — Error Handling & Structured Logging
- [ ] Phase 6 — Authentication (JWT)
- [ ] Phase 7 — Concurrency & Performance
- [ ] Phase 8 — Docker & Deployment

## Author

**Rohit Sadawarte** — Building production-grade backend systems in Go.
