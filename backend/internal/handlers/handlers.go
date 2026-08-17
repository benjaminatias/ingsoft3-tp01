// Package handlers contiene los handlers HTTP de Gin.
// Los handlers utilizan GORM directamente: no existe una capa Repository.
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler agrupa las dependencias de los handlers. Solamente necesita la base de datos.
type Handler struct {
	DB *gorm.DB
}

// Nuevo crea un Handler con la conexión indicada.
func Nuevo(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

// responderError devuelve un error con una estructura simple y consistente.
// Nunca se envían stack traces al frontend.
func responderError(c *gin.Context, estado int, mensaje string) {
	c.JSON(estado, gin.H{"error": mensaje})
}

// leerID obtiene el parámetro :id de la ruta.
// Si no es válido responde 400 y devuelve ok = false.
func leerID(c *gin.Context) (uint, bool) {
	valor, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || valor == 0 {
		responderError(c, http.StatusBadRequest, "El identificador indicado no es válido.")
		return 0, false
	}
	return uint(valor), true
}

// RegistrarRutas registra todas las rutas de la API en el router de Gin.
func RegistrarRutas(router *gin.Engine, h *Handler) {
	router.GET("/health", h.Health)

	api := router.Group("/api")

	// Rutas públicas: registrarse e iniciar sesión.
	api.POST("/auth/registro", h.Registro)
	api.POST("/auth/login", h.Login)

	// El resto de la API exige un token JWT válido.
	privado := api.Group("")
	privado.Use(RequiereAutenticacion())
	{
		privado.GET("/auth/perfil", h.Perfil)

		privado.GET("/peliculas", h.ListarPeliculas)
		privado.GET("/peliculas/:id", h.ObtenerPelicula)
		privado.POST("/peliculas", h.CrearPelicula)
		privado.PUT("/peliculas/:id", h.ActualizarPelicula)
		privado.DELETE("/peliculas/:id", h.EliminarPelicula)
		privado.PATCH("/peliculas/:id/estado", h.CambiarEstado)
		privado.PATCH("/peliculas/:id/puntuacion", h.PuntuarPelicula)

		privado.GET("/generos", h.ListarGeneros)
		privado.GET("/generos/:id", h.ObtenerGenero)
		privado.POST("/generos", h.CrearGenero)
		privado.PUT("/generos/:id", h.ActualizarGenero)
		privado.DELETE("/generos/:id", h.EliminarGenero)

		privado.GET("/estadisticas", h.ObtenerEstadisticas)
	}
}

// NuevoRouter crea el router completo de la aplicación.
func NuevoRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	RegistrarRutas(router, Nuevo(db))
	return router
}
