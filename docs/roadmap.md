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

## FASE 4 — Búsqueda y filtros
- [ ] GET /products?q=nombre — búsqueda por nombre
- [ ] Paginación básica (?page=1&limit=10)

## FASE 5 — Buenas prácticas
- [ ] Manejo de errores consistente (equivalente a HttpException en NestJS)
- [ ] Middleware de logging (equivalente a interceptors)
- [ ] Respuestas con estructura estandarizada { data, error, status }

## FASE 6 — Testing
- [ ] Test unitario de un handler con `testing` package nativo de Go
- [ ] Test de integración de un endpoint

## FASE 7 — Base de datos real
- [ ] Conectar PostgreSQL con GORM
- [ ] Migrar el slice en memoria a DB real
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
