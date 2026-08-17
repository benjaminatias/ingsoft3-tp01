package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Health comprueba que la API responda y que PostgreSQL sea accesible.
// Se utiliza en Docker Compose para determinar si el backend está saludable.
func (h *Handler) Health(c *gin.Context) {
	if h.DB == nil {
		responderNoSaludable(c)
		return
	}

	sqlDB, err := h.DB.DB()
	if err != nil {
		responderNoSaludable(c)
		return
	}

	ctx, cancelar := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancelar()

	if err := sqlDB.PingContext(ctx); err != nil {
		responderNoSaludable(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func responderNoSaludable(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"status": "unhealthy",
		"error":  "No se pudo conectar con la base de datos.",
	})
}
