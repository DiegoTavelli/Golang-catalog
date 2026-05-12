# go.mod — explicación

Equivalente exacto a `package.json` en Node/TS.
Go lo genera automáticamente con `go mod init`.

---

## Línea por línea

```
module github.com/diegotavelli/go-catalog
```
El nombre único de tu módulo. Go usa URLs como nombres para evitar colisiones globales.
No hace falta que el repo exista en GitHub — es solo un identificador.
Equivalente a: el campo `"name"` en package.json.

---

```
go 1.22.x
```
La versión mínima de Go que necesita este proyecto.
Equivalente a: el campo `"engines": { "node": ">=18" }` en package.json.

---

```
require (
    github.com/gin-gonic/gin v1.12.0
    ...
)
```
Las dependencias del proyecto con su versión exacta.
Equivalente a: el bloque `"dependencies"` en package.json.

Las marcadas con `// indirect` son dependencias de dependencias
(gin necesita estas librerías internamente).
Equivalente a: lo que npm mete en node_modules pero no en tu package.json directamente.

---

## Comandos equivalentes

| Go | npm/yarn |
|---|---|
| `go mod init nombre` | `npm init` |
| `go get paquete` | `npm install paquete` |
| `go mod tidy` | `npm install` (limpia y sincroniza todo) |
| `go mod download` | `npm ci` (descarga sin modificar) |
