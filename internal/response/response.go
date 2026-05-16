package response

// Este package define la estructura estandarizada de todas las respuestas de la API.
//
// El problema sin esto: cada endpoint responde distinto:
//   GET /products      → { "data": [...], "total": 10 }
//   GET /products/:id  → { "data": {...} }
//   DELETE             → { "message": "...", "id": 1 }
//   Error              → { "error": "..." }
//
// Con un formato estándar el frontend siempre sabe qué esperar.
// Equivalente a: un interceptor de respuesta en NestJS que envuelve todo en { data, error, status }
// O a: el patrón ApiResponse<T> que se ve en muchas APIs REST.

// Success arma la respuesta estándar para casos exitosos.
// Equivalente a: res.status(200).json({ data, meta }) en Express
func Success(data any) map[string]any {
	return map[string]any{
		"data": data,
	}
}

// Paginated arma la respuesta estándar para listados paginados.
// Incluye meta con info de paginación — el frontend la necesita para renderizar el paginador.
func Paginated(data any, total, page, limit int) map[string]any {
	return map[string]any{
		"data": data,
		"meta": map[string]any{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	}
}

// Error arma la respuesta estándar para errores.
// Equivalente a: throw new HttpException({ error: msg }, status) en NestJS
func Error(message string) map[string]any {
	return map[string]any{
		"error": message,
	}
}
