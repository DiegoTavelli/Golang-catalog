# Roadmap — Go Product Catalog API

Mini proyecto para aprender Go orientado a entrevistas backend (MercadoLibre).
Comparaciones constantes con NestJS/TypeScript para acelerar el aprendizaje.

---

## FASE 1 — Setup e instalación ✅
- [x] Instalar Go en Windows
- [x] Verificar instalación (`go version`)
- [x] Crear el módulo del proyecto (`go mod init`)
- [x] Instalar gin (el framework web)
- [x] Levantar un servidor HTTP

## FASE 2 — Estructura del proyecto ✅
- [x] Organizar carpetas (equivalente a modules de NestJS)
- [x] Crear el struct `Product` y `CreateProductInput` (modelo + DTO)
- [x] Crear el "repositorio" en memoria (slice de productos)
- [x] Separar rutas, handlers y modelos en archivos distintos

## FASE 3 — CRUD completo ✅
- [x] GET /products — listar todos los productos
- [x] GET /products?category=X — filtrar por categoría
- [x] GET /products/:id — obtener uno por ID
- [x] POST /products — crear producto con validación
- [x] PUT /products/:id — actualizar producto
- [x] DELETE /products/:id — eliminar producto

## FASE 4 — Búsqueda y filtros ✅
- [x] GET /products?q=nombre — búsqueda por nombre (case-insensitive)
- [x] GET /products?category=X&q=Y — combinación de filtros (pipeline de filtros)
- [x] Paginación básica (?page=1&limit=10) con guard clauses

## FASE 5 — Estructura profesional y capas
El proyecto actual tiene toda la lógica en el handler: filtra, pagina y responde
en la misma función. En producción eso no escala. Esta fase frena y reorganiza.

**Por qué este momento:** entre el CRUD (Fase 3-4) y el error handling (Fase 6)
es el lugar natural para separar capas — antes de agregar más features encima de
una estructura que ya empieza a doler.

- [ ] Extraer `service/product_service.go` — lógica de negocio (filter, paginate)
- [ ] Extraer `repository/product_repository.go` — acceso a datos (slice en memoria por ahora)
- [ ] Handler queda solo como capa HTTP: leer params → llamar service → responder JSON
- [ ] Entender `internal/` vs `pkg/` — visibilidad enforceada por el compilador
- [ ] Estructura de carpetas estándar Go: `cmd/`, `internal/`, `pkg/`
- [ ] Documentar el flujo handler → service → repository en arquitectura.md

## FASE 6 — Buenas prácticas
- [ ] Manejo de errores consistente (equivalente a HttpException en NestJS)
- [ ] Middleware de logging (equivalente a interceptors)
- [ ] Respuestas con estructura estandarizada { data, error, status }

## FASE 7 — Testing
- [ ] Test unitario de un service (lógica pura, sin HTTP)
- [ ] Test de integración de un endpoint con httptest
- [ ] Por qué separar capas hace el testing trivial

## FASE 8 — Base de datos real
- [ ] Conectar PostgreSQL con GORM
- [ ] Migrar el repository (solo esa capa cambia — handler y service no se tocan)
- [ ] CRUD completo con persistencia

---

## Conceptos estudiados

| Concepto Go | Equivalente NestJS/TS | Estado |
|---|---|---|
| `struct` | `class` / interface / DTO | ✅ |
| `gin.Context` | `req, res` de Express | ✅ |
| `go.mod` / `go.sum` | `package.json` / `yarn.lock` | ✅ |
| `go get` | `npm install` | ✅ |
| `:=` vs `var` vs `const` | `let` / `var` / `const` | ✅ |
| Early return / Guard clauses | misma práctica en JS | ✅ |
| Punteros (`*` y `&`) | referencias en JS | ✅ |
| Packages y visibilidad | módulos + `export` en TS | ✅ |
| Manejo de errores como valor | try/catch en JS | ✅ |
| `range` | `for...of` / `.forEach()` | ✅ |
| `append()` | `.push()` (pero devuelve nuevo slice) | ✅ |
| Actualizar struct por índice (`slice[i].Field`) | spread/assign en JS | ✅ |
| Eliminar de slice con `append([:i], [i+1:]...)` | `.filter()` en JS | ✅ |
| `ShouldBindJSON` vs `BindJSON` | pipe + validación en NestJS | ✅ |
| `if err := fn(); err != nil` | try/catch compacto | ✅ |
| `func` con receiver | método de clase | ⬜ |
| `interface{}` / `any` | `any` en TypeScript | ⬜ |
| Goroutines | async/await (concepto distinto) | ⬜ |
