package repository

// El repository es la capa más cercana a los datos.
// Su única responsabilidad: saber cómo leer y escribir datos — nada más.
// No sabe qué reglas de negocio existen. No sabe qué significa "filtrar por categoría".
// Solo sabe: "dame todos", "dame el de ID=3", "guardá este", "borrá ese".
//
// Equivalente a:
//   @Injectable() export class ProductRepository extends Repository<Product> {} ← NestJS + TypeORM
//   O el patrón Repository manual: class ProductRepository { findAll() { return db.query(...) } }
//
// Hoy usa un slice en memoria — en Fase 8 este archivo cambia a GORM/PostgreSQL.
// El handler y el service NO se van a tocar cuando eso pase.
// Eso es lo que hace valiosa la separación en capas.

import "github.com/DiegoTavelli/Golang-catalog/internal/model"

// products simula una base de datos en memoria
//
// var a nivel de package (fuera de funciones) = variable global del package
// Solo este package puede acceder a ella directamente — es privada (minúscula)
// Equivalente a: una variable de módulo en TS, o un array in-memory en un servicio NestJS de prueba
var products = []model.Product{
	// {} inicializa un struct — equivalente a un objeto literal en JS
	{ID: 1, Name: "Notebook Lenovo", Price: 1200.00, Category: "electronics", Stock: 10},
	{ID: 2, Name: "Mouse Logitech", Price: 35.00, Category: "electronics", Stock: 50},
	{ID: 3, Name: "Silla Gamer", Price: 450.00, Category: "furniture", Stock: 5},
}

// nextID simula un autoincrement de base de datos
// en una app real esto lo maneja Postgres automáticamente con SERIAL o BIGSERIAL
var nextID = 4

// FindAll devuelve todos los productos sin filtrar.
// El service es el que filtra — el repository solo entrega los datos crudos.
//
// Nota: devuelve el slice directamente (no una copia).
// En Go los slices son referencias a un array subyacente — pasar un slice no copia los datos.
// Equivalente a: return this.products en un servicio NestJS
func FindAll() []model.Product {
	return products
}

// FindByID busca un producto por ID y devuelve DOS valores: el producto y un bool.
//
// Este es un patrón muy común en Go: devolver (valor, ok) en vez de valor | null.
// En JS harías: return products.find(p => p.id === id) ?? null
// En Go no existe null para structs — si no encontrás nada, devolvés el "valor cero" del struct
// y false como señal de "no encontrado".
//
// La firma (model.Product, bool) es el equivalente Go de: Product | undefined en TS
func FindByID(id int) (model.Product, bool) {
	// Go no tiene .find() — usás un for con range
	// range sobre un slice devuelve (índice, copia del elemento)
	// como no necesitamos el índice acá, lo descartamos con _
	// Equivalente a: products.find(p => p.ID === id)
	for _, p := range products {
		if p.ID == id {
			return p, true // encontrado → devolvemos el producto y true
		}
	}
	return model.Product{}, false // no encontrado → valor cero del struct y false
}

// Create agrega un nuevo producto al slice y devuelve el producto creado con su ID asignado.
//
// Recibe un CreateProductInput (sin ID) y construye un Product (con ID).
// La asignación del ID es responsabilidad del repository — simula lo que haría la DB.
// Equivalente a: const saved = await this.repo.save(input) en TypeORM
func Create(input model.CreateProductInput) model.Product {
	p := model.Product{
		ID:       nextID,
		Name:     input.Name,
		Price:    input.Price,
		Category: input.Category,
		Stock:    input.Stock,
	}
	nextID++ // equivalente a nextID += 1 — el autoincrement simulado

	// append devuelve un NUEVO slice — en Go no existe .push() que mute in place
	// hay que reasignar: products = append(products, p)
	// Equivalente a: products.push(p) en JS, pero Go requiere la reasignación explícita
	products = append(products, p)
	return p
}

// Update modifica un producto existente y devuelve (producto actualizado, encontrado).
//
// Patrón clave de Go: buscamos por ÍNDICE (i), no por valor (p).
// Si usáramos "for _, p := range products" obtendríamos una COPIA de cada producto.
// Modificar la copia no afecta el slice original — bug clásico en Go para quien viene de JS.
// En JS todos los objetos son referencias, acá no — hay que ser explícito.
//
// Equivalente a:
//   const index = products.findIndex(p => p.id === id)
//   products[index] = { ...products[index], ...input }
func Update(id int, input model.CreateProductInput) (model.Product, bool) {
	for i, p := range products {
		if p.ID == id {
			// products[i] accede al elemento ORIGINAL del slice (no a la copia p)
			// actualizamos campo por campo — Go no tiene spread ni Object.assign()
			products[i].Name = input.Name
			products[i].Price = input.Price
			products[i].Category = input.Category
			products[i].Stock = input.Stock
			return products[i], true
		}
	}
	return model.Product{}, false
}

// Delete elimina un producto del slice y devuelve true si existía.
//
// En Go no hay .filter() ni .splice() built-in para slices.
// El patrón idiomático: construir un nuevo slice sin el elemento que queremos borrar.
//
//	products[:i]   → todos los elementos ANTES del índice
//	products[i+1:] → todos los elementos DESPUÉS del índice
//	append() une las dos partes saltando el elemento eliminado
//	El ... desempaqueta el segundo slice para que append lo reciba elemento por elemento
//	  (equivalente al spread ... de JS pero solo válido dentro de append)
//
// Equivalente a: products = products.filter(p => p.id !== id)
func Delete(id int) bool {
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return true
		}
	}
	return false
}
