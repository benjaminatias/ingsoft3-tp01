package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"gestor-peliculas/internal/models"
	"gestor-peliculas/internal/validation"
)

type peticionGenero struct {
	Nombre *string `json:"nombre"`
}

// ListarGeneros devuelve todos los géneros ordenados por nombre.
func (h *Handler) ListarGeneros(c *gin.Context) {
	generos := []models.Genero{}
	if err := h.DB.Order("nombre ASC").Find(&generos).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudieron obtener los géneros.")
		return
	}
	c.JSON(http.StatusOK, generos)
}

// ObtenerGenero devuelve un género por su identificador.
func (h *Handler) ObtenerGenero(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	var genero models.Genero
	err := h.DB.First(&genero, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(c, http.StatusNotFound, "El género no existe.")
		return
	}
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo obtener el género.")
		return
	}

	c.JSON(http.StatusOK, genero)
}

// CrearGenero registra un género nuevo, sin permitir duplicados.
func (h *Handler) CrearGenero(c *gin.Context) {
	nombre, ok := leerNombreGenero(c)
	if !ok {
		return
	}

	duplicado, err := h.nombreGeneroRepetido(nombre, 0)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear el género.")
		return
	}
	if duplicado {
		responderError(c, http.StatusConflict, "Ya existe un género con ese nombre.")
		return
	}

	genero := models.Genero{Nombre: nombre}
	if err := h.DB.Create(&genero).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo crear el género.")
		return
	}

	c.JSON(http.StatusCreated, genero)
}

// ActualizarGenero cambia el nombre de un género existente.
func (h *Handler) ActualizarGenero(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	nombre, ok := leerNombreGenero(c)
	if !ok {
		return
	}

	var genero models.Genero
	err := h.DB.First(&genero, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(c, http.StatusNotFound, "El género no existe.")
		return
	}
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo actualizar el género.")
		return
	}

	duplicado, err := h.nombreGeneroRepetido(nombre, genero.ID)
	if err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo actualizar el género.")
		return
	}
	if duplicado {
		responderError(c, http.StatusConflict, "Ya existe un género con ese nombre.")
		return
	}

	genero.Nombre = nombre
	if err := h.DB.Save(&genero).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo actualizar el género.")
		return
	}

	c.JSON(http.StatusOK, genero)
}

// EliminarGenero borra un género.
// Si existen películas que lo utilizan responde 409 y no elimina nada.
func (h *Handler) EliminarGenero(c *gin.Context) {
	id, ok := leerID(c)
	if !ok {
		return
	}

	var cantidad int64
	if err := h.DB.Model(&models.Pelicula{}).Where("genero_id = ?", id).Count(&cantidad).Error; err != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo eliminar el género.")
		return
	}
	if cantidad > 0 {
		responderError(c, http.StatusConflict, "No se puede eliminar el género porque existen películas que lo utilizan.")
		return
	}

	resultado := h.DB.Delete(&models.Genero{}, id)
	if resultado.Error != nil {
		responderError(c, http.StatusInternalServerError, "No se pudo eliminar el género.")
		return
	}
	if resultado.RowsAffected == 0 {
		responderError(c, http.StatusNotFound, "El género no existe.")
		return
	}

	c.Status(http.StatusNoContent)
}

// leerNombreGenero lee y valida el nombre recibido en el cuerpo de la petición.
func leerNombreGenero(c *gin.Context) (string, bool) {
	var peticion peticionGenero
	if err := c.ShouldBindJSON(&peticion); err != nil {
		responderError(c, http.StatusBadRequest, "El cuerpo de la petición no es válido.")
		return "", false
	}
	if peticion.Nombre == nil {
		responderError(c, http.StatusBadRequest, "El nombre del género es obligatorio.")
		return "", false
	}

	nombre := validation.LimpiarTexto(*peticion.Nombre)
	if err := validation.ValidarNombreGenero(nombre); err != nil {
		responderError(c, http.StatusBadRequest, err.Error())
		return "", false
	}

	return nombre, true
}

// nombreGeneroRepetido comprueba duplicados ignorando mayúsculas y minúsculas.
// excluirID permite ignorar el propio género al editarlo.
func (h *Handler) nombreGeneroRepetido(nombre string, excluirID uint) (bool, error) {
	var cantidad int64
	consulta := h.DB.Model(&models.Genero{}).Where("LOWER(nombre) = LOWER(?)", nombre)
	if excluirID != 0 {
		consulta = consulta.Where("id <> ?", excluirID)
	}
	if err := consulta.Count(&cantidad).Error; err != nil {
		return false, err
	}
	return cantidad > 0, nil
}
