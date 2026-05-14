package main // package main es especial en Go — indica el punto de entrada del programa
// Go busca este package para saber desde dónde arrancar
// Equivalente a: el archivo main.ts de NestJS con el bootstrap()

// En la estructura estándar, main.go vive en cmd/api/ y no en la raíz.
// Esto permite que el proyecto tenga múltiples puntos de entrada:
//   cmd/api/     → el servidor HTTP (este)
//   cmd/worker/  → un worker de background (en el futuro)
//   cmd/migrate/ → un script de migración de DB (en el futuro)
// Cada uno tiene su propio main.go pero comparten todo el código de internal/

import (
	// importamos solo el package handler — main.go no conoce service ni repository
	// esa es la idea: main solo conecta las piezas, no implementa lógica
	"github.com/DiegoTavelli/Golang-catalog/internal/handler"
	"github.com/gin-gonic/gin" // framework web — equivalente a express o @nestjs/core
)

func main() {

	// gin.Default() crea el router con dos middlewares ya incluidos:
	//   Logger:   imprime cada request (método, ruta, status, tiempo de respuesta)
	//   Recovery: si un handler hace "panic" (crash), lo captura y devuelve 500 — el server no cae
	// Equivalente a:
	//   const app = express(); app.use(morgan('dev')); app.use(errorHandler)
	//   o en NestJS: NestFactory.create(AppModule) — que incluye todo esto automáticamente
	router := gin.Default()

	// Group() agrupa rutas bajo un prefijo común — todas quedan con /products automáticamente
	// Equivalente a:
	//   const router = express.Router(); app.use('/products', router)  ← Express
	//   @Controller('products')                                         ← NestJS
	products := router.Group("/products")
	{
		// Las llaves {} son solo estilo visual — no crean un scope nuevo en Go
		// Es una convención para que quede claro que estas rutas pertenecen al grupo
		// Equivalente a los métodos del controller en NestJS: @Get(), @Post(), etc.
		products.GET("", handler.GetProducts)        // GET /products
		products.GET("/:id", handler.GetProductByID) // GET /products/:id
		products.POST("", handler.CreateProduct)     // POST /products
		products.PUT("/:id", handler.UpdateProduct)  // PUT /products/:id
		products.DELETE("/:id", handler.DeleteProduct) // DELETE /products/:id
	}

	// Run() levanta el servidor en el puerto y bloquea el proceso (loop infinito esperando requests)
	// El ":" antes del número significa "escuchar en todas las interfaces de red en este puerto"
	// Equivalente a: app.listen(8080) en Express
	router.Run(":8080")
}
