package handlers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"gestor-peliculas/internal/models"
	"gestor-peliculas/internal/validation"
)

// estadisticas resume la colección. La puntuación promedio es null cuando
// todavía no existe ninguna película puntuada.
type estadisticas struct {
	TotalPeliculas     int64    `json:"totalPeliculas"`
	Vistas             int64    `json:"vistas"`
	Pendientes         int64    `json:"pendientes"`
	PuntuacionPromedio *float64 `json:"puntuacionPromedio"`
}

// ObtenerEstadisticas devuelve información simple sobre la colección.
func (h *Handler) ObtenerEstadisticas(c *gin.Context) {
	var resultado estadisticas

	if err := h.DB.Model(&models.Pelicula{}).Count(&resultado.TotalPeliculas).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener las estadísticas.")
		return
	}
	if err := h.DB.Model(&models.Pelicula{}).Where("estado = ?", validation.EstadoVista).Count(&resultado.Vistas).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener las estadísticas.")
		return
	}
	if err := h.DB.Model(&models.Pelicula{}).Where("estado = ?", validation.EstadoPendiente).Count(&resultado.Pendientes).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener las estadísticas.")
		return
	}

	// El promedio considera solamente películas vistas que tengan puntuación.
	// Si no hay ninguna, AVG devuelve NULL y el promedio queda en null.
	var fila struct {
		Promedio *float64
	}
	if err := h.DB.Model(&models.Pelicula{}).
		Where("estado = ? AND puntuacion IS NOT NULL", validation.EstadoVista).
		Select("AVG(puntuacion) AS promedio").
		Scan(&fila).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener las estadísticas.")
		return
	}

	if fila.Promedio != nil {
		redondeado := math.Round(*fila.Promedio*10) / 10
		resultado.PuntuacionPromedio = &redondeado
	}

	c.JSON(http.StatusOK, resultado)
}
