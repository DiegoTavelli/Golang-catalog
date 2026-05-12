# Golang Catalog API

REST API built in Go as a learning project.

Coming from 4+ years of NestJS/TypeScript backend experience in production fintech environments,
this project is my structured path into Go — covering the same backend patterns I already know,
but learning the language, idioms and tooling from scratch.

## What it does

Product catalog REST API with full CRUD, filtering by category and input validation.

## Stack

- **Go** — language
- **Gin** — web framework (equivalent to Express in Node.js)
- In-memory storage for now — database integration coming in next phases

## Endpoints

| Method | Route | Description |
|---|---|---|
| GET | /products | List all products |
| GET | /products?category=X | Filter by category |
| GET | /products/:id | Get product by ID |
| POST | /products | Create product |
| PUT | /products/:id | Update product (coming soon) |
| DELETE | /products/:id | Delete product (coming soon) |

## Run locally

```bash
git clone https://github.com/diegotavelli/golang-catalog
cd golang-catalog
go mod tidy
go run main.go
```

Server starts at `http://localhost:8080`

## Project structure

```
go-catalog/
├── main.go            # entry point, router setup
├── handlers/
│   └── products.go    # request handlers (equivalent to controllers in NestJS)
├── models/
│   └── product.go     # structs and DTOs
├── go.mod             # dependency manifest (equivalent to package.json)
└── go.sum             # dependency lockfile (equivalent to package-lock.json)
```

## Learning roadmap

- [x] Project setup, gin framework, first HTTP server
- [x] GET endpoints with query filtering
- [x] POST with input validation
- [ ] PUT and DELETE
- [ ] Error handling middleware
- [ ] Pagination
- [ ] Unit testing with Go's native testing package
- [ ] Connect to a real database (PostgreSQL)
