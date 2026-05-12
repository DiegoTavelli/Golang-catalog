package handlers // equivalente a un @Controller en NestJS
                 // agrupa todas las funciones que manejan requests de productos

import (
	"net/http" // librería estándar de Go para constantes HTTP (StatusOK, StatusNotFound, etc.)
	           // equivalente a: HttpStatus de @nestjs/common

	"strconv" // librería estándar para conversiones de tipos — string a int, int a string, etc.
	          // en TS harías parseInt() o toString() directamente, en Go hay una librería dedicada

	"github.com/gin-gonic/gin"
	"github.com/diegotavelli/go-catalog/models"
)

// products simula una base de datos en memoria
//
// var declara una variable a nivel de package (fuera de funciones) — equivalente a una variable global en TS
// []models.Product es un "slice" — el array dinámico de Go
//
// Diferencia entre array y slice en Go:
//   [3]Product → array fijo de exactamente 3 elementos (raro de usar)
//   []Product  → slice, tamaño dinámico, como los arrays de JS
//
// Equivalente a:
//   const products: Product[] = [...]
var products = []models.Product{
	// {} inicializa un struct — equivalente a un objeto literal en JS
	// podés nombrar los campos (Name: "valor") o no, pero nombrarlos es más claro
	{ID: 1, Name: "Notebook Lenovo",  Price: 1200.00, Category: "electronics", Stock: 10},
	{ID: 2, Name: "Mouse Logitech",   Price: 35.00,   Category: "electronics", Stock: 50},
	{ID: 3, Name: "Silla Gamer",      Price: 450.00,  Category: "furniture",   Stock: 5},
}

// nextID simula un autoincrement de base de datos
// en una app real esto lo maneja Postgres/Mongo automáticamente
var nextID = 4

// GetProducts maneja GET /products
//
// En Go las funciones handler reciben un *gin.Context — el objeto central de gin
// Contiene todo: el request, el response, los params, el body, los headers, etc.
// Equivalente a:
//   (req: Request, res: Response) en Express
//   (@Req() req, @Res() res) en NestJS
//   O directamente @Query(), @Param(), @Body() en NestJS decorators
//
// El * antes de gin.Context significa "puntero a gin.Context"
// Un puntero es una referencia al objeto en memoria — no una copia
// Go pasa los parámetros por valor (copia) por defecto
// Con *gin.Context le decimos "pasame la referencia, no copies todo el objeto"
// Equivalente a: pasar objetos por referencia en JS (todos los objetos en JS son referencias)
func GetProducts(c *gin.Context) {

	// c.Query() lee query parameters de la URL
	// GET /products?category=electronics → c.Query("category") devuelve "electronics"
	// GET /products → c.Query("category") devuelve "" (string vacío, no nil/undefined)
	// Equivalente a: req.query.category en Express
	category := c.Query("category")

	// en Go no hay if/else ternario ni valores "falsy"
	// "" (string vacío) no es false — tenés que comparar explícitamente
	if category == "" {

		// c.JSON serializa el segundo argumento a JSON y lo manda como response
		// primer argumento: HTTP status code
		// segundo argumento: cualquier dato — gin.H es un shortcut para map[string]any{}
		//
		// gin.H{"data": products} genera: { "data": [...] }
		// Equivalente a: res.status(200).json({ data: products }) en Express
		c.JSON(http.StatusOK, gin.H{
			"data": products,
		})
		return // return corta la función — en Go no hay "else" implícito después de responder
		       // si no ponés return, Go sigue ejecutando el código de abajo (bug clásico)
	}

	// declaramos un slice vacío para acumular resultados filtrados
	// var filtered []models.Product → equivalente a: const filtered: Product[] = []
	// nil en Go es el valor cero de un slice — distinto de un slice vacío, pero ambos funcionan con append
	var filtered []models.Product

	// range itera sobre el slice — equivalente a for...of en JS
	// devuelve dos valores: el índice y el elemento
	// _ descarta el índice (Go da error si declarás una variable y no la usás)
	// Equivalente a:
	//   for (const p of products) { ... }
	//   products.forEach(p => { ... })
	for _, p := range products {
		if p.Category == category {
			// append() agrega un elemento al slice y devuelve el nuevo slice
			// IMPORTANTE: en Go append puede crear un slice nuevo internamente
			// por eso siempre reasignás: filtered = append(filtered, p)
			// Equivalente a: filtered.push(p) — pero push muta el array, append devuelve uno nuevo
			filtered = append(filtered, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": filtered})
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
