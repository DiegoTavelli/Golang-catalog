package handler

// El handler es la capa HTTP — la única que conoce gin.Context.
// Sus únicas responsabilidades:
//   1. Leer parámetros del request (query params, path params, body)
//   2. Llamar al service con esos parámetros ya parseados
//   3. Responder con el status code correcto y el JSON
//
// Lo que NO hace: lógica de negocio, acceso a datos, validaciones de dominio.
// Si el handler crece más allá de "leer → llamar → responder", algo está mal.
//
// Equivalente a:
//   @Controller('products') export class ProductsController { ... } ← NestJS
//   En NestJS los decorators hacen el parseo automáticamente (@Query(), @Param(), @Body())
//   En Go lo hacemos explícitamente — más verbose pero más claro qué está pasando.

import (
	"net/http" // constantes HTTP: StatusOK, StatusNotFound, StatusBadRequest, etc.
	"strconv"  // conversiones de tipo: string → int (Atoi), int → string (Itoa), etc.
	// en JS harías parseInt() o toString() — en Go hay una librería dedicada para esto

	"github.com/DiegoTavelli/Golang-catalog/internal/model"
	"github.com/DiegoTavelli/Golang-catalog/internal/service"
	"github.com/gin-gonic/gin"
)

// GetProducts maneja GET /products con soporte de filtros y paginación.
//
// c *gin.Context es el objeto central de gin — contiene todo: request, response, params, headers.
// El * significa "puntero a gin.Context" — recibimos la referencia, no una copia del objeto.
// Go pasa todo por valor (copia) por defecto. Con * le decimos "pasame la referencia".
// Equivalente a: (req: Request, res: Response) en Express, o los decoradores de NestJS.
func GetProducts(c *gin.Context) {

	// c.Query() lee query params — devuelve "" si el param no existe (nunca nil/undefined)
	// Equivalente a: req.query.category en Express
	category := c.Query("category")
	search := c.Query("q")

	// valores por defecto para paginación
	page, limit := 1, 10 // declaración múltiple en una línea — equivalente a: let page = 1, limit = 10

	// strconv.Atoi convierte "42" → 42 (Atoi = ASCII to Integer)
	// devuelve DOS valores: (número, error) — Go siempre devuelve errores como valores, no excepciones
	// el patrón "if err := fn(); err == nil" ejecuta la función y entra al if solo si no hubo error
	// Equivalente a: const p = parseInt(pageStr); if (!isNaN(p) && p > 0) page = p
	if p, err := strconv.Atoi(c.Query("page")); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	// llamamos al service con los parámetros ya tipados (string, string, int, int)
	// el service no sabe que esto vino de un query string — recibe tipos Go normales
	// Equivalente a: return this.productsService.findAll({ category, search, page, limit })
	result := service.GetProducts(category, search, page, limit)

	// c.JSON serializa el struct a JSON y escribe la respuesta HTTP
	// gin.H{} es un alias de map[string]any — forma rápida de armar un objeto JSON
	// Equivalente a: res.status(200).json({ data: ..., page: ..., ... })
	c.JSON(http.StatusOK, gin.H{
		"data":  result.Data,
		"page":  result.Page,
		"limit": result.Limit,
		"total": result.Total,
	})
}

// GetProductByID maneja GET /products/:id
func GetProductByID(c *gin.Context) {

	// c.Param() lee parámetros de ruta — GET /products/42 → c.Param("id") devuelve "42" (string)
	// Equivalente a: req.params.id en Express
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// el ID no era un número — respondemos 400 y salimos con return
		// "early return" o "guard clause" — patrón igual en JS/TS
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	// el service devuelve (producto, encontrado) — patrón (valor, bool) de Go
	// equivalente a: const product = await this.service.findById(id) → product ?? throw new NotFoundException()
	p, found := service.GetProductByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// CreateProduct maneja POST /products
func CreateProduct(c *gin.Context) {

	// declaramos la variable con valor cero — en Go todas las variables tienen valor cero al nacer
	// string → ""  |  int → 0  |  float64 → 0.0  |  struct → todos los campos en cero
	// Equivalente a: const input: CreateProductInput = {} (pero tipado, sin undefined)
	var input model.CreateProductInput

	// ShouldBindJSON hace dos cosas en un solo llamado:
	//   1. Parsea el body JSON al struct input
	//   2. Valida los campos según los struct tags (binding:"required", binding:"gt=0")
	// El & antes de input es el operador "dirección de" — pasamos un puntero para que gin pueda modificar input
	// Sin &, gin recibiría una copia de input y los cambios se perderían (Go pasa por valor por defecto)
	// Equivalente a: @Body() input: CreateProductDto  ← NestJS con class-validator hace esto automáticamente
	if err := c.ShouldBindJSON(&input); err != nil {
		// err.Error() devuelve el mensaje de error como string — equivalente a err.message en JS
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := service.CreateProduct(input)
	c.JSON(http.StatusCreated, gin.H{"data": p}) // 201 Created — status correcto para un POST que crea
}

// UpdateProduct maneja PUT /products/:id
func UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	var input model.CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// el service devuelve (productoActualizado, encontrado)
	// si found es false → el ID no existía → 404
	p, found := service.UpdateProduct(id, input)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// DeleteProduct maneja DELETE /products/:id
func DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID debe ser un número"})
		return
	}

	// el service devuelve un bool simple — true si existía y fue eliminado
	// el handler es el que decide qué status HTTP corresponde a cada caso
	// esta decisión HTTP no le pertenece al service — esa es la razón de separar capas
	if !service.DeleteProduct(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado", "id": id})
}
