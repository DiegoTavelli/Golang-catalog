# Arquitectura de Go — contexto para este proyecto

---

## Cómo piensa Go (vs JavaScript)

JavaScript es dinámico y flexible — podés pasar cualquier cosa a cualquier función.
Go es estático y explícito — todo tiene un tipo, todo se declara, nada es implícito.

Esto cambia cómo se estructura el código:

| JS/TS | Go |
|---|---|
| Clases con herencia | Structs + interfaces (sin herencia) |
| try/catch | error como valor de retorno |
| Cualquier función es async | Goroutines (concurrencia real, no event loop) |
| npm + node_modules por proyecto | go mod + cache global |
| Tipado opcional (TS) | Tipado obligatorio siempre |
| Decorators (@Body, @Get) | Explícito — vos escribís cada paso |

---

## Packages — la unidad de organización

En Go no hay clases ni módulos como en TS. La unidad es el **package**.

```
go-catalog/
├── main.go          → package main   (punto de entrada)
├── handlers/
│   └── products.go  → package handlers
└── models/
    └── product.go   → package models
```

Regla: todos los archivos en la misma carpeta deben tener el mismo package name.
Para usar código de otro package: `import "github.com/DiegoTavelli/Golang-catalog/handlers"`

Equivalente en NestJS: cada carpeta es como un módulo, pero sin decorators — Go lo resuelve con imports directos.

---

## Visibilidad — mayúscula vs minúscula

En Go no hay `public` / `private`. La regla es simple:

```go
type Product struct { ... }    // Mayúscula → exportado (público)
type productHelper struct { }  // Minúscula → privado al package
```

```go
func GetProducts(c *gin.Context) { }  // Mayúscula → usable desde otros packages
func buildQuery(filter string) { }    // Minúscula → solo dentro de handlers/
```

Equivalente a: `export` en TS. Si no ponés `export`, es privado al módulo.

---

## Structs — reemplazan las clases

Go no tiene clases ni herencia. Los structs son la forma de agrupar datos:

```go
type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

Para agregar comportamiento, usás funciones con **receiver** — equivalente a métodos de clase:

```go
// esto es como un método: product.IsAvailable()
func (p Product) IsAvailable() bool {
    return p.Stock > 0
}
```

Sin herencia — Go prefiere composición: un struct puede incluir otro struct adentro.

---

## Punteros — cuándo y por qué

Go pasa todo por valor (copia) por defecto. Los punteros permiten pasar referencias:

```go
func fill(input *CreateProductInput) {  // recibe referencia
    input.Name = "algo"                 // modifica el original
}

fill(&myInput)  // & = dame la dirección de memoria
```

Regla práctica para este proyecto:
- `*gin.Context` → siempre puntero (gin lo requiere)
- `&input` en ShouldBindJSON → para que gin pueda llenarlo
- Structs propios → por valor está bien mientras sean pequeños

---

## El flujo de un request en este proyecto

```
HTTP Request
    ↓
gin router (main.go)
    ↓ matchea la ruta
handler function (handlers/products.go)
    ↓ lee params/body
    ↓ valida con ShouldBindJSON
    ↓ aplica lógica
    ↓ construye respuesta
c.JSON(status, data)
    ↓
HTTP Response
```

Sin middlewares de por medio por ahora. gin.Default() ya incluye logger y recovery automáticamente.

---

## Aprendizajes — CRUD completo (Fase 3)

### Actualizar un struct: índice vs copia

En Go el `range` devuelve una copia del elemento, no una referencia:

```go
// ❌ esto NO modifica el original
for _, p := range products {
    p.Name = "nuevo"  // modifica la copia, el slice no cambia
}

// ✅ esto SÍ modifica el original
for i, p := range products {
    if p.ID == id {
        products[i].Name = "nuevo"  // accede directamente al slice
    }
}
```

Equivalente en JS: la diferencia entre mutar el objeto referenciado vs reasignar la variable.

---

### Eliminar de un slice — patrón idiomático

Go no tiene `.splice()` ni `.filter()` built-in. El patrón estándar:

```go
// eliminar el elemento en índice i
products = append(products[:i], products[i+1:]...)
```

- `products[:i]` → todo antes del índice
- `products[i+1:]` → todo después del índice
- `...` → desempaqueta el segundo slice para que append lo reciba elemento por elemento
- `append` concatena las dos partes → resultado sin el elemento eliminado
- `products =` → reasignamos porque append devuelve un slice nuevo

Equivalente a: `products = products.filter(p => p.ID !== id)`

---

### ShouldBindJSON — validación y parseo en uno

```go
var input models.CreateProductInput
if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

Hace dos cosas: parsea el JSON del body AL struct, y valida los campos según los struct tags (`binding:"required"`, `binding:"gt=0"`). Si falla cualquiera de los dos devuelve error.

La forma `if err := fn(); err != nil` declara `err` con scope solo dentro del if — patrón muy común en Go para mantener el código limpio.

---

### El orden importa: modificar antes de responder

```go
products = append(...)           // 1. modificás el estado
c.JSON(200, gin.H{...})          // 2. informás al cliente
```

Siempre en ese orden. El `c.JSON` no "aplica" nada — solo serializa y manda la respuesta HTTP. La modificación real ya ocurrió antes. Con DB real es igual: `db.Delete()` primero, `c.JSON` después.

---

## Lo que viene — con DB real (GORM)

Cuando conectemos PostgreSQL, el struct Product va a tener tags adicionales:

```go
type Product struct {
    ID       uint    `gorm:"primaryKey" json:"id"`
    Name     string  `gorm:"not null"   json:"name"`
    Price    float64                    `json:"price"`
}
```

GORM mapea el struct a la tabla automáticamente.
El handler va a cambiar: en vez de buscar en el slice `products`, va a llamar a `db.Find()` o `db.First()`.
La estructura del proyecto no cambia — solo el handler deja de usar memoria y empieza a usar la DB.
