package model // todos los archivos en esta carpeta pertenecen al package "model"
// en Go el package name NO tiene que coincidir con la carpeta, pero por convención sí coincide
// antes teníamos package "models" (plural) — el estándar Go usa singular: "model"

// Product representa un producto en el catálogo.
//
// En Go no hay clases — los structs son la forma de modelar datos con forma fija.
// La diferencia clave con JS: en JS un objeto puede tener cualquier campo en runtime.
// En Go el compilador sabe exactamente qué campos existen y de qué tipo son.
//
// Equivalente a:
//
//	interface Product { id: number; name: string; price: number; ... }  ← TS
//	@Entity() class Product { @PrimaryGeneratedColumn() id: number; ... } ← TypeORM
type Product struct {

	// Los backticks son "struct tags" — metadatos que otras librerías leen en runtime
	// json:"id" le dice al serializador de Go: "cuando conviertas esto a JSON, usá la clave 'id'"
	// sin el tag, Go usaría "ID" (mayúscula) en el JSON — no es convención REST
	// Equivalente a: @Expose({ name: 'id' }) en class-transformer, o simplemente camelCase en TS
	ID   int    `json:"id"`
	Name string `json:"name"`

	// float64 = punto flotante de 64 bits — siempre usá float64 para precios o medidas
	// Go tiene float32 (menor precisión) y float64 (más preciso) — en JS no existe esa distinción
	// Equivalente a: number en TypeScript (JS internamente usa float64 para todo)
	Price float64 `json:"price"`

	Category string `json:"category"`
	Stock    int    `json:"stock"` // int = entero sin decimales — equivalente a number en TS pero sin punto flotante
}

// CreateProductInput es el DTO de entrada para POST y PUT.
//
// Separamos modelo de datos (Product) del modelo de entrada (CreateProductInput) porque:
// - Product tiene ID → lo genera el servidor, el cliente no lo manda
// - CreateProductInput no tiene ID → solo los campos que el cliente puede enviar
//
// Esta separación es igual que en NestJS:
//
//	class CreateProductDto {
//	  @IsNotEmpty() name: string
//	  @IsPositive() price: number
//	}
//
// En Go la validación se declara en el mismo struct con struct tags de binding:
// binding:"required" → campo obligatorio (equivale a @IsNotEmpty())
// binding:"gt=0"     → greater than 0 (equivale a @IsPositive())
// gin usa la librería go-playground/validator internamente, igual que class-validator en NestJS
type CreateProductInput struct {
	Name     string  `json:"name"     binding:"required"`
	Price    float64 `json:"price"    binding:"required,gt=0"`
	Category string  `json:"category" binding:"required"`
	Stock    int     `json:"stock"    binding:"required"`
}
