# Golang Catalog API

REST API built in Go as a learning project.

Coming from 4+ years of NestJS/TypeScript backend experience in production fintech environments,
this project is my structured path into Go — covering the same backend patterns I already know,
but learning the language, idioms and tooling from scratch.

## What it does

Product catalog REST API with full CRUD, text search, category filtering, pagination and
a professional layered architecture (handler → service → repository).

## Stack

- **Go** — language
- **Gin** — web framework (equivalent to Express in Node.js)
- In-memory storage for now — PostgreSQL with GORM coming in next phase

## Endpoints

| Method | Route | Description |
|---|---|---|
| GET | /products | List all products |
| GET | /products?category=X | Filter by category |
| GET | /products?q=name | Search by name (case-insensitive) |
| GET | /products?category=X&q=name | Combined filters |
| GET | /products?page=1&limit=5 | Pagination |
| GET | /products/:id | Get product by ID |
| POST | /products | Create product |
| PUT | /products/:id | Update product |
| DELETE | /products/:id | Delete product |

## Run locally

```bash
git clone https://github.com/DiegoTavelli/Golang-catalog
cd Golang-catalog
go mod tidy
go run ./cmd/api/
```

Server starts at `http://localhost:8080`

> Note: `go run ./cmd/api/` instead of `go run main.go` — the entry point lives in
> `cmd/api/` following the standard Go project layout.

## Project structure

```
go-catalog/
├── cmd/
│   └── api/
│       └── main.go              # entry point, router setup
├── internal/
│   ├── handler/
│   │   └── product_handler.go  # HTTP layer — reads params, calls service, returns JSON
│   ├── service/
│   │   └── product_service.go  # business logic — filtering, pagination, rules
│   ├── repository/
│   │   └── product_repository.go # data layer — in-memory slice (GORM next)
│   └── model/
│       └── product.go           # structs and DTOs
├── pkg/
│   └── pagination/              # reusable utilities (coming soon)
├── docs/
│   ├── arquitectura.md          # architecture notes and Go vs NestJS comparisons
│   ├── roadmap.md               # learning phases and progress
│   └── insomnia-collection.json # API collection ready to import
├── go.mod                       # dependency manifest (equivalent to package.json)
└── go.sum                       # dependency lockfile (equivalent to package-lock.json)
```

## Learning roadmap

- [x] Project setup, Gin framework, first HTTP server
- [x] Full CRUD (GET, POST, PUT, DELETE)
- [x] Input validation with struct tags
- [x] Text search and category filtering
- [x] Pagination
- [x] Layered architecture — handler / service / repository separation
- [ ] Error handling middleware and standardized response structure
- [ ] Unit testing with Go's native testing package
- [ ] PostgreSQL with GORM
