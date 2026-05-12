# Log del proyecto — Go Product Catalog API

Registro paso a paso de todo lo que se fue haciendo y por qué.

---

## PASO 1 — Instalar Go ✅

Instalado desde go.dev/dl — instalador Windows .msi.
Go viene con todo incluido: compilador, gestor de módulos, test runner, formatter.
No necesitás nada extra (equivalente a Node + npm + jest + prettier todo en uno).

Verificación: `go version` → `go version go1.22.x windows/amd64`

---

## PASO 2 — Crear el módulo ✅

```bash
go mod init github.com/diegotavelli/go-catalog
```

Esto creó el archivo `go.mod` — equivalente a `package.json`.
El nombre del módulo usa formato de URL para evitar colisiones globales.
No hace falta que exista el repo en GitHub, es solo un identificador único.

Go también sugirió correr `go mod tidy` — ese comando limpia y sincroniza dependencias.
Es el equivalente de `npm install` cuando ya tenés un package.json.

---

## PASO 3 — Instalar gin ✅

```bash
go get github.com/gin-gonic/gin
```

Esto descargó gin y todas sus dependencias transitivas.
Actualizó `go.mod` con las dependencias y creó `go.sum` con los hashes de verificación.

Las dependencias se guardan en un cache global en `C:\Users\Diego\go\pkg\mod\`
— no hay node_modules por proyecto, todo es centralizado.

**Error que apareció:** se escribió `gin-tonic` en vez de `gin-gonic` → repositorio no encontrado.
Corrección: el paquete correcto es `github.com/gin-gonic/gin`.

---

## PASO 4 — Estructura del proyecto ✅

Se crearon tres archivos:

```
go-catalog/
├── main.go               ← punto de entrada, configura el router y las rutas
├── models/
│   └── product.go        ← structs: Product (modelo) y CreateProductInput (DTO)
└── handlers/
    └── products.go       ← funciones que manejan cada endpoint
```

Equivalente en NestJS:
```
src/
├── main.ts               ← bootstrap()
├── products/
│   ├── product.entity.ts ← @Entity() class
│   ├── create-product.dto.ts
│   └── products.controller.ts
```

---

## PASO 5 — Servidor funcionando ✅

```bash
go run main.go
```

`go run` compila y ejecuta en un solo paso — equivalente a `ts-node` o `npx ts-node`.
Para producción usarías `go build` que genera un binario ejecutable standalone.

Endpoints probados con Insomnia — todos respondiendo correctamente:
- GET  /products               → lista los 3 productos
- GET  /products?category=electronics → filtra por categoría
- GET  /products/1             → devuelve el producto con ID 1
- POST /products               → crea un producto nuevo con validación

---

## Conceptos nuevos que aparecieron en el camino

**Punteros (&):** cuando pasás `&input` a ShouldBindJSON, le das la dirección de memoria
del struct para que gin lo pueda modificar directamente. Sin &, gin recibiría una copia
y los datos del body se perderían. En JS todos los objetos son referencias automáticamente.

**Manejo de errores:** Go no tiene try/catch. Las funciones devuelven el error como segundo
valor de retorno. Siempre verificás `if err != nil` — es más explícito que las excepciones.

**Valor cero:** en Go toda variable declarada tiene un valor inicial automático.
string → "", int → 0, bool → false, slice → nil. No existe undefined.

**Slice vs Array:** []Product es un slice (dinámico, como arrays de JS).
[3]Product sería un array fijo de exactamente 3 elementos (raro de usar).
