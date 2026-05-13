package main // todo archivo Go pertenece a un "package".
             // "main" es especial — indica que este es el punto de entrada del programa.
             // Sin package main, Go no sabe desde dónde arrancar.
             // Equivalente a: el archivo main.ts de NestJS que tiene el bootstrap()

// import agrupa todas las dependencias externas e internas que usa este archivo
// Go no permite imports sin usar — si importás algo y no lo usás, el compilador da error
// Equivalente a: los import al tope de cada archivo TS, pero más estricto
import (
	"github.com/gin-gonic/gin"                        // framework web — equivalente a express
	"github.com/DiegoTavelli/Golang-catalog/handlers"     // nuestro package local de handlers
)

// func main() es el punto de entrada del programa
// Go busca esta función en el package main para saber por dónde arrancar
// Equivalente a: async function bootstrap() { const app = await NestFactory.create(AppModule); app.listen(3000) }
func main() {

	// gin.Default() crea una instancia del router con dos middlewares ya configurados:
	//   - Logger: imprime en consola cada request (método, ruta, status, tiempo)
	//   - Recovery: si un handler hace "panic" (crash), lo captura y devuelve 500 en vez de tirar el servidor
	// Equivalente a:
	//   const app = express()
	//   app.use(morgan('dev'))          ← logger
	//   app.use(errorHandlerMiddleware) ← recovery
	router := gin.Default()

	// Group() agrupa rutas bajo un prefijo común
	// Todas las rutas definidas adentro van a tener /products como prefijo automáticamente
	// Equivalente a:
	//   const productsRouter = express.Router()
	//   app.use('/products', productsRouter)
	// O en NestJS: @Controller('products')
	products := router.Group("/products")
	{
		// las llaves {} son solo estilo visual para agrupar — no crean un scope nuevo en Go
		// es una convención para que quede claro que estas rutas pertenecen al grupo

		products.GET("", handlers.GetProducts)        // GET /products
		products.GET("/:id", handlers.GetProductByID) // GET /products/1, /products/2, etc.
		products.POST("", handlers.CreateProduct)     // POST /products
	}

	// Run() levanta el servidor HTTP en el puerto indicado y bloquea el proceso
	// El ":" antes del número es la sintaxis de Go para "escuchar en todas las interfaces en este puerto"
	// Equivalente a: app.listen(8080)
	router.Run(":8080")
}
