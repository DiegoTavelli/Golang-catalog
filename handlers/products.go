package handlers // equivalente a un @Controller en NestJS
// agrupa todas las funciones que manejan requests de productos

import (
	"net/http" // librería estándar de Go para constantes HTTP (StatusOK, StatusNotFound, etc.)
	// equivalente a: HttpStatus de @nestjs/common

	"strconv" // librería estándar para conversiones de tipos — string a int, int a string, etc.
	// en TS harías parseInt() o toString() directamente, en Go hay una librería dedicada

	"strings" // librería estándar para operaciones con strings — Contains, ToLower, TrimSpace, etc.
	// equivalente a: métodos nativos de string en JS (.includes(), .toLowerCase(), .trim())

	"github.com/DiegoTavelli/Golang-catalog/models"
	"github.com/gin-gonic/gin"
)

// products simula una base de datos en memoria
//
// var declara una variable a nivel de package (fuera de funciones) — equivalente a una variable global en TS
// []models.Product es un "slice" — el array dinámico de Go
//
// Diferencia entre array y slice en Go:
//
//	[3]Product → array fijo de exactamente 3 elementos (raro de usar)
//	[]Product  → slice, tamaño dinámico, como los arrays de JS
//
// Equivalente a:
//
//	const products: Product[] = [...]
var products = []models.Product{
	// {} inicializa un struct — equivalente a un objeto literal en JS
	// podés nombrar los campos (Name: "valor") o no, pero nombrarlos es más claro
	{ID: 1, Name: "Notebook Lenovo", Price: 1200.00, Category: "electronics", Stock: 10},
	{ID: 2, Name: "Mouse Logitech", Price: 35.00, Category: "electronics", Stock: 50},
	{ID: 3, Name: "Silla Gamer", Price: 450.00, Category: "furniture", Stock: 5},
}

// nextID simula un autoincrement de base de datos
// en una app real esto lo maneja Postgres/Mongo automáticamente
var nextID = 4

// GetProducts maneja GET /products
//
// En Go las funciones handler reciben un *gin.Context — el objeto central de gin
// Contiene todo: el request, el response, los params, el body, los headers, etc.
// Equivalente a:
//
//	(req: Request, res: Response) en Express
//	(@Req() req, @Res() res) en NestJS
//	O directamente @Query(), @Param(), @Body() en NestJS decorators
//
// El * antes de gin.Context significa "puntero a gin.Context"
// Un puntero es una referencia al objeto en memoria — no una copia
// Go pasa los parámetros por valor (copia) por defecto
// Con *gin.Context le decimos "pasame la referencia, no copies todo el objeto"
// Equivalente a: pasar objetos por referencia en JS (todos los objetos en JS son referencias)
func GetProducts(c *gin.Context) {

	// leemos todos los query params de una vez
	// c.Query() devuelve "" si el param no existe — nunca nil/undefined
	// Equivalente a: req.query.category en Express
	category := c.Query("category")
	search := c.Query("q")        // GET /products?q=notebook → busca por nombre
	pageStr := c.Query("page")    // GET /products?page=2
	limitStr := c.Query("limit")  // GET /products?limit=10

	// --- PASO 1: filtrar por category y búsqueda por nombre ---

	// empezamos con todos los productos y vamos filtrando
	// esta técnica se llama "pipeline de filtros" — cada condición reduce el slice
	// Equivalente a: products.filter(p => ...).filter(p => ...)
	var result []models.Product

	for _, p := range products {

		// filtro por categoría — si category está vacío, este filtro no aplica
		if category != "" && p.Category != category {
			continue // "saltá este elemento y seguí con el siguiente" — equivalente a continue en JS
		}

		// filtro por búsqueda de texto en el nombre
		// strings.Contains(s, substr) → equivalente a s.includes(substr) en JS
		// strings.ToLower() → equivalente a s.toLowerCase() en JS
		// los aplicamos a ambos lados para que la búsqueda sea case-insensitive
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			continue
		}

		// si pasó todos los filtros, lo agregamos al resultado
		result = append(result, p)
	}

	// --- PASO 2: paginación ---

	// valores por defecto si no se pasan los params
	// strconv.Atoi devuelve error si el string no es número — en ese caso usamos el default
	page := 1
	limit := 10

	// el segundo valor de retorno (err) lo descartamos con _ si no nos importa el error
	// en este caso: si page=abc, Atoi falla y usamos el default 1 — comportamiento correcto
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	// calculamos el índice de inicio y fin del "slice" de resultados
	// page=1, limit=10 → start=0, end=10
	// page=2, limit=10 → start=10, end=20
	// Equivalente a lo que haría .skip() y .take() en TypeORM
	start := (page - 1) * limit
	end := start + limit

	// guard clause — si start supera el total de resultados, devolvemos lista vacía
	// sin esto, products[start:end] podría hacer panic (equivalente a un crash en Go)
	if start >= len(result) {
		c.JSON(http.StatusOK, gin.H{
			"data":  []models.Product{}, // slice vacío explícito — mejor que null en la respuesta
			"page":  page,
			"limit": limit,
			"total": len(result),
		})
		return
	}

	// ajustamos end si supera el largo del slice
	// sin esto: result[0:15] cuando result tiene 8 elementos → panic
	if end > len(result) {
		end = len(result)
	}

	// result[start:end] es un "sub-slice" — equivalente a arr.slice(start, end) en JS
	// no copia los datos, solo crea una nueva "vista" del slice original
	paginated := result[start:end]

	c.JSON(http.StatusOK, gin.H{
		"data":  paginated,
		"page":  page,
		"limit": limit,
		"total": len(result), // total de resultados sin paginar — útil para el frontend
	})
}

// GetProductByID maneja GET /products/:id
func GetProductByID(c *gin.Context) {

	// c.Param() lee parámetros de ruta
	// GET /products/42 → c.Param("id") devuelve "42" (string, no número)
	// Equivalente a: req.params.id en Express
	idStr := c.Param("id")

	// strconv.Atoi convierte "42" → 42 (Atoi = ASCII to Integer)
	// Go siempre devuelve errores como valores — no hay excepciones ni try/catch
	// la función devuelve DOS valores: el resultado y el error
	// si err == nil → todo bien | si err != nil → algo falló
	//
	// Equivalente a:
	//   const id = parseInt(idStr)
	//   if (isNaN(id)) { res.status(400).json({ error: '...' }) }
	//
	// El manejo explícito de errores es una de las diferencias más grandes con JS/TS
	// En TS usás try/catch para excepciones — en Go el error es un valor de retorno normal
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	// buscamos el producto en el slice
	// Go no tiene .find() como JS — usás un for con range
	// Equivalente a: products.find(p => p.ID === id)
	for _, p := range products {
		if p.ID == id {
			c.JSON(http.StatusOK, gin.H{"data": p})
			return // encontramos el producto, respondemos y salimos
		}
	}

	// si el for termina sin hacer return → no encontramos el producto
	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

// CreateProduct maneja POST /products
func CreateProduct(c *gin.Context) {

	// declaramos una variable del tipo CreateProductInput con valor cero
	// en Go todas las variables tienen un "valor cero" al declararse:
	//   string → ""  |  int → 0  |  float64 → 0.0  |  bool → false  |  struct → todos los campos en cero
	// Equivalente a: const input: CreateProductInput = {} (pero tipado y sin undefined)
	var input models.CreateProductInput

	// ShouldBindJSON intenta:
	//   1. Parsear el body JSON al struct input
	//   2. Validar los campos según las struct tags (binding:"required", binding:"gt=0", etc.)
	// Si falla cualquiera de los dos → devuelve error
	//
	// El & antes de input es el operador "dirección de" — pasamos un puntero al struct
	// Gin necesita el puntero para poder modificar input directamente (llenarlo con los datos del body)
	// Sin &, Go pasaría una copia de input y los cambios se perderían
	//
	// Equivalente a:
	//   @Body() input: CreateProductDto  ← NestJS hace esto automáticamente con decorators
	//   const input = plainToClass(CreateProductDto, req.body)  ← versión manual con class-transformer
	if err := c.ShouldBindJSON(&input); err != nil {
		// err.Error() devuelve el mensaje de error como string
		// Equivalente a: err.message en JS
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// construimos el nuevo producto con los datos validados
	// nextID++ incrementa el contador — equivalente a nextID += 1
	newProduct := models.Product{
		ID:       nextID,
		Name:     input.Name,
		Price:    input.Price,
		Category: input.Category,
		Stock:    input.Stock,
	}
	nextID++

	// agregamos el producto al slice global
	products = append(products, newProduct)

	// 201 Created es el status correcto para un POST que crea un recurso
	// Equivalente a: res.status(201).json({ data: newProduct })
	c.JSON(http.StatusCreated, gin.H{"data": newProduct})
}

// UpdateProduct maneja PUT /products/:id
//
// En Go no existe el "spread operator" ni Object.assign() de JS
// Para actualizar un struct, accedés a cada campo directamente por índice
// Esto es más verbose pero más explícito — sabés exactamente qué estás cambiando
func UpdateProduct(c *gin.Context) {

	// mismo patrón de siempre — leer y convertir el ID del param de ruta
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	// parseamos y validamos el body — igual que en CreateProduct
	var input models.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// buscamos el producto por índice, no por valor
	// necesitamos el ÍNDICE (i) porque vamos a modificar el slice original
	// si usáramos "for _, p := range products" obtendríamos una COPIA de cada producto
	// modificar la copia no afecta el slice original — bug clásico en Go
	//
	// Equivalente a:
	//   const index = products.findIndex(p => p.ID === id)
	//   if (index === -1) return res.status(404)...
	//   products[index] = { ...products[index], ...input }
	for i, p := range products {
		if p.ID == id {

			// products[i] accede al elemento original del slice (no a una copia)
			// actualizamos campo por campo — en Go no hay spread ni merge automático
			products[i].Name = input.Name
			products[i].Price = input.Price
			products[i].Category = input.Category
			products[i].Stock = input.Stock

			// respondemos con el producto actualizado
			c.JSON(http.StatusOK, gin.H{"data": products[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
}

// DeleteProduct maneja DELETE /products/:id
//
// En Go no hay .filter() como en JS
// Para eliminar un elemento de un slice, se construye un slice nuevo sin ese elemento
// Es un patrón idiomático de Go — se ve raro al principio pero es la forma correcta
func DeleteProduct(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	// buscamos el índice del producto a eliminar
	// -1 indica "no encontrado" — mismo concepto que findIndex() en JS
	indexToDelete := -1
	for i, p := range products {
		if p.ID == id {
			indexToDelete = i
			break // encontramos el índice, no hace falta seguir iterando
		}
	}

	// si no encontramos el producto, respondemos 404
	if indexToDelete == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	// eliminamos el elemento del slice — este es el patrón idiomático de Go:
	//
	// products[:indexToDelete]  → todos los elementos ANTES del índice
	// products[indexToDelete+1:] → todos los elementos DESPUÉS del índice
	// append() los une saltando el elemento que queremos borrar
	//
	// Equivalente a: products.filter(p => p.ID !== id)
	//
	// Ejemplo visual con índice 1:
	//   slice original: [A, B, C, D]
	//   [:1]          → [A]
	//   [2:]          → [C, D]
	//   resultado     → [A, C, D]
	//
	// El ... (spread) en Go desempaqueta el slice para que append pueda recibirlo
	// Equivalente al spread operator ... de JS pero solo para slices en append
	products = append(products[:indexToDelete], products[indexToDelete+1:]...)

	// 200 con mensaje — algunos usan 204 No Content (sin body)
	// usamos 200 para ser explícitos y confirmar qué se eliminó
	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado", "id": id})
}
