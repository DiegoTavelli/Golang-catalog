package pagination

// Este package vive en pkg/ porque no tiene nada de "productos" — es lógica pura
// que funciona igual para cualquier slice de cualquier tipo.
// Otro proyecto Go podría importar esto sin cambiar una línea.
//
// Equivalente a: una utility function genérica en TS
//   function paginate<T>(items: T[], page: number, limit: number): Page<T>

// Result es el struct de respuesta de una paginación.
//
// El [T any] después del nombre es la sintaxis de GENERICS en Go (disponible desde Go 1.18).
// T es un "type parameter" — un tipo que se define en el momento de usar el struct, no al definirlo.
// "any" es la restricción: T puede ser cualquier tipo (equivalente a any en TS).
//
// Equivalente en TypeScript:
//
//	interface Page<T> { data: T[]; total: number; page: number; limit: number }
type Result[T any] struct {
	Data  []T
	Total int
	Page  int
	Limit int
}

// Paginate aplica paginación sobre cualquier slice y devuelve un Result[T].
//
// La firma [T any] hace que esta función acepte []Product, []User, []Order — lo que sea.
// Go infiere el tipo T automáticamente según lo que le pases — no necesitás escribirlo.
//
// Uso:
//
//	result := pagination.Paginate(products, 1, 10)  // T = models.Product, inferido solo
//	result := pagination.Paginate(users, 2, 5)      // T = models.User, inferido solo
//
// Equivalente en TypeScript:
//
//	function paginate<T>(items: T[], page: number, limit: number): Page<T>
//	paginate(products, 1, 10)  // TypeScript infiere T = Product
func Paginate[T any](items []T, page, limit int) Result[T] {
	total := len(items)

	// calculamos el índice de inicio en el slice
	// page=1, limit=10 → start=0  |  page=2, limit=10 → start=10
	// Equivalente a: .skip() y .take() en TypeORM
	start := (page - 1) * limit

	// guard clause: si start supera el total no hay nada que devolver
	// []T{} es un slice vacío del tipo T — se serializa como [] en JSON, no como null
	if start >= total {
		return Result[T]{Data: []T{}, Total: total, Page: page, Limit: limit}
	}

	end := start + limit
	if end > total {
		end = total
	}

	// items[start:end] es un sub-slice — no copia datos, crea una vista del slice original
	// Equivalente a: items.slice(start, end) en JS
	return Result[T]{Data: items[start:end], Total: total, Page: page, Limit: limit}
}
