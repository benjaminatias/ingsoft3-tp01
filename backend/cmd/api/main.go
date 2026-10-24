// Comando principal de la API de Gestor de Películas.
package main

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	"gestor-peliculas/internal/database"
	"gestor-peliculas/internal/handlers"
)

func main() {
	// 0. Sentry: el DSN llega por el ENTORNO (SENTRY_DSN), no se escribe en el código. Si la variable
	// no está (en mi máquina, en los tests, en QA), Sentry queda apagado y la aplicación anda igual.
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      os.Getenv("SENTRY_ENVIRONMENT"),
		TracesSampleRate: 1.0, // medir todos los pedidos (alcanza para una aplicación de cursada)
	}); err != nil {
		log.Printf("No se pudo iniciar Sentry (la aplicación sigue sin él): %v", err)
	}
	defer sentry.Flush(2 * time.Second)

	// 1. Conectarse a PostgreSQL con las variables de entorno.
	db, err := database.Conectar()
	if err != nil {
		log.Fatalf("Error de conexión con la base de datos: %v", err)
	}

	// 2. Ejecutar migraciones.
	if err := database.Migrar(db); err != nil {
		log.Fatalf("Error al ejecutar las migraciones: %v", err)
	}

	// 3. Crear los géneros iniciales.
	if err := database.CrearGenerosIniciales(db); err != nil {
		log.Fatalf("Error al crear los géneros iniciales: %v", err)
	}

	// 4. Inicializar Gin y registrar las rutas.
	router := handlers.NuevoRouter(db)

	// 5. Ejecutar el servidor.
	puerto := os.Getenv("SERVER_PORT")
	if puerto == "" {
		puerto = "8080"
	}

	log.Printf("Gestor de Películas escuchando en http://localhost:%s", puerto)
	if err := router.Run(":" + puerto); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func init() {
	// En producción Gin no necesita el modo debug.
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
}
