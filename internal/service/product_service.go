package service

// El service contiene la lógica de negocio — las reglas de la aplicación.
// No sabe nada de HTTP (no conoce gin.Context, no arma respuestas JSON).
// No sabe nada de cómo se guardan los datos (no conoce el slice ni la DB).
// Solo sabe: "dame productos filtrados y paginados", "creá uno con estas reglas".
//
// Equivalente a:
//   @Injectable() export class ProductsService { ... } ← NestJS
//
// Esta separación tiene un beneficio enorme para testing:
// podés testear toda la lógica de negocio sin levantar un servidor HTTP.
// En NestJS harías: const service = new ProductsService(mockRepo)
// En Go: simplemente llamás service.GetProducts(...) en el test — sin mocks de gin.

import (
	"strings" // strings.Contains, strings.ToLower — equivalente a str.includes(), str.toLowerCase() en JS

	"github.com/DiegoTavelli/Golang-catalog/internal/model"
	"github.com/DiegoTavelli/Golang-catalog/internal/repository"
)

// ProductList es el struct de respuesta del listado paginado.
//
// En NestJS devolverías un objeto plano: { data: Product[], total: number, page: number, limit: number }
// En Go definimos un struct para eso — tipado explícito, sin objetos anónimos.
// El handler recibe este struct y lo convierte a JSON con gin.H{}.
//
// Nota: este struct vive en el package service porque es una construcción interna
// del service — no es un modelo de DB ni un DTO de entrada del cliente.
type ProductList struct {
	Data  []model.Product
	Total int
	Page  int
	Limit int
}

// GetProducts aplica filtros y paginación sobre todos los productos.
//
// Recibe los parámetros ya parseados (el handler convirtió los strings a int antes de llamar acá).
// Esta función no sabe que esos parámetros vinieron de un query string HTTP.
// Equivalente a: findAll(filter: FilterDto): Promise<ProductList> en un service NestJS
func GetProducts(category, search string, page, limit int) ProductList {
	// pedimos todos los datos al repository — esta capa no accede al slice directamente
	all := repository.FindAll()

	// pipeline de filtros — técnica: empezamos con todo y vamos descartando
	// cada condición reduce el slice con "continue" (saltear este elemento)
	// Equivalente a: all.filter(p => ...).filter(p => ...) en JS
	// Go prefiere un solo for con múltiples condiciones sobre múltiples .filter() encadenados
	var filtered []model.Product // var sin inicializar = nil en Go, pero append funciona igual sobre nil

	for _, p := range all {
		// si el filtro de categoría está activo y no coincide → saltear
		if category != "" && p.Category != category {
			continue // equivalente a continue en JS dentro de un for
		}
		// búsqueda case-insensitive: convertimos ambos lados a minúscula antes de comparar
		// strings.Contains(s, substr) → equivalente a s.includes(substr) en JS
		// strings.ToLower(s)         → equivalente a s.toLowerCase() en JS
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			continue
		}
		// si pasó todos los filtros → lo incluimos
		filtered = append(filtered, p)
	}

	// paginación — la aplicamos sobre el resultado filtrado, no sobre todos los productos
	total := len(filtered) // len() → equivalente a arr.length en JS

	// calculamos el índice de inicio en el slice
	// page=1, limit=10 → start=0  (empieza desde el primer elemento)
	// page=2, limit=10 → start=10 (salta los primeros 10)
	// Equivalente a lo que haría .skip() y .take() en TypeORM
	start := (page - 1) * limit

	// guard clause: si start supera el total, no hay nada para devolver
	// sin esto, filtered[start:end] haría panic (equivalente a un crash — no hay try/catch en Go)
	if start >= total {
		// []model.Product{} es un slice vacío explícito — mejor que nil en la respuesta JSON
		// nil se serializa como null en JSON, []model.Product{} se serializa como []
		return ProductList{Data: []model.Product{}, Total: total, Page: page, Limit: limit}
	}

	end := start + limit
	if end > total { // ajustamos end si supera el largo del slice
		end = total
	}

	// filtered[start:end] es un "sub-slice" — no copia los datos, crea una vista del slice original
	// Equivalente a: filtered.slice(start, end) en JS
	return ProductList{Data: filtered[start:end], Total: total, Page: page, Limit: limit}
}

// GetProductByID delega directamente al repository.
// El service es el intermediario — en el futuro podría agregar lógica acá
// (por ejemplo: registrar que alguien consultó este producto, o aplicar permisos).
func GetProductByID(id int) (model.Product, bool) {
	return repository.FindByID(id)
	// Go permite devolver múltiples valores — acá "pasamos" los dos valores del repository
	// directamente sin reasignarlos. Equivalente a: return this.repo.findById(id) en TS
}

// CreateProduct delega al repository y devuelve el producto ya creado con su ID asignado.
func CreateProduct(input model.CreateProductInput) model.Product {
	return repository.Create(input)
}

// UpdateProduct devuelve (producto actualizado, encontrado).
// El handler usa el bool para decidir si responde 200 o 404 — esa decisión HTTP no le pertenece al service.
func UpdateProduct(id int, input model.CreateProductInput) (model.Product, bool) {
	return repository.Update(id, input)
}

// DeleteProduct devuelve true si el producto existía y fue eliminado, false si no existía.
// El handler convierte ese bool en 200 o 404 — la capa HTTP decide el status code, no el service.
func DeleteProduct(id int) bool {
	return repository.Delete(id)
}
