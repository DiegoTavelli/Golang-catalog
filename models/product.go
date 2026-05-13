package models // este archivo pertenece al package "models"
// en Go, el package name determina cómo otros archivos importan esto:
//   import "github.com/DiegoTavelli/Golang-catalog/models"
//   models.Product{} ← así se usa
// Equivalente a: una carpeta entities/ o dto/ en NestJS

// Product representa un producto en el catálogo
//
// En Go no hay clases — los structs son la forma de modelar datos.
// Un struct es como un objeto con forma fija: sabés exactamente qué campos tiene y de qué tipo.
// Equivalente a:
//
//	interface Product {
//	  id: number
//	  name: string
//	  price: number
//	  category: string
//	  stock: number
//	}
//
// O a una Entity de TypeORM:
//
//	@Entity()
//	class Product {
//	  @PrimaryGeneratedColumn() id: number
//	  @Column() name: string
//	  ...
//	}
type Product struct {

	// los backticks al final de cada campo son "struct tags"
	// le dicen a Go cómo serializar/deserializar este campo cuando trabaja con JSON
	// json:"id" → cuando convierta a JSON, este campo se llama "id" (minúscula)
	// sin el tag, Go usaría "ID" (mayúscula) en el JSON, que no es convención REST
	// Equivalente a: @Expose() @Transform() o simplemente camelCase en class-transformer

	ID   int    `json:"id"`
	Name string `json:"name"`

	// float64 es el tipo de punto flotante de 64 bits en Go
	// Go tiene float32 y float64 — siempre usá float64 para dinero o medidas (más precisión)
	// Equivalente a: number en TypeScript (JS no distingue entre enteros y decimales)
	Price float64 `json:"price"`

	Category string `json:"category"`
	Stock    int    `json:"stock"` // int = entero — equivalente a number en TS pero sin decimales
}

// CreateProductInput es el struct que valida el body de un POST /products
//
// Separar el modelo de entrada del modelo de datos es buena práctica:
// - Product tiene ID (que el cliente no manda, lo genera el servidor)
// - CreateProductInput no tiene ID
// Equivalente a: un DTO (Data Transfer Object) en NestJS
//
//	class CreateProductDto {
//	  @IsNotEmpty() name: string
//	  @IsPositive() price: number
//	  ...
//	}
type CreateProductInput struct {

	// binding:"required" activa validación cuando gin parsea el JSON
	// si el campo falta o está vacío, gin devuelve error automáticamente
	// Equivalente a: @IsNotEmpty() de class-validator en NestJS
	Name string `json:"name"         binding:"required"`

	// gt=0 significa "greater than 0" — el precio debe ser positivo
	// Equivalente a: @IsPositive() de class-validator
	Price float64 `json:"price"      binding:"required,gt=0"`

	Category string `json:"category" binding:"required"`
	Stock    int    `json:"stock"    binding:"required"`
}
