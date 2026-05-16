package middleware

// Un middleware en Go/gin es una función que se ejecuta ANTES o DESPUÉS de cada handler.
// Equivalente a: un @Injectable() middleware o un interceptor en NestJS,
// o app.use(fn) en Express.
//
// gin ya incluye un logger básico en gin.Default() — este lo reemplaza con
// un formato más limpio y controlado, útil para entender cómo funciona el mecanismo.

import (
	"fmt"
	"time" // librería estándar de Go para manejo de tiempo y duración

	"github.com/gin-gonic/gin"
)

// Logger devuelve una gin.HandlerFunc — que es simplemente func(*gin.Context).
// Al devolver la función en vez de ser la función, podemos parametrizarla en el futuro
// (por ejemplo: Logger(prefix string) para distintos servicios).
//
// Equivalente a:
//   export function LoggerMiddleware(): NestMiddleware { use(req, res, next) { ... } }
func Logger() gin.HandlerFunc {

	// esta función anónima es la que gin va a ejecutar en cada request
	return func(c *gin.Context) {

		// capturamos el tiempo ANTES de procesar el request
		start := time.Now()

		// c.Next() le dice a gin "seguí con el siguiente handler en la cadena"
		// todo lo que escribamos ANTES de Next() se ejecuta antes del handler
		// todo lo que escribamos DESPUÉS de Next() se ejecuta después del handler
		// Equivalente a: await next() en un middleware de NestJS o Express
		c.Next()

		// llegamos acá DESPUÉS de que el handler respondió
		// time.Since(start) calcula cuánto tiempo pasó desde "start"
		// Equivalente a: Date.now() - startTime en JS
		duration := time.Since(start)

		// c.Writer.Status() devuelve el status code que escribió el handler
		status := c.Writer.Status()

		// fmt.Printf imprime a stdout con formato — equivalente a console.log() en JS
		// %s = string  |  %d = integer  |  %v = cualquier valor (Go lo formatea solo)
		fmt.Printf("[API] %s %s → %d (%v)\n",
			c.Request.Method, // GET, POST, PUT, DELETE
			c.Request.URL.Path, // /products, /products/1
			status,             // 200, 201, 404, etc.
			duration,           // 142µs, 1.2ms, etc.
		)
	}
}
